/*
Copyright 2020 FUJITSU CLOUD TECHNOLOGIES LIMITED. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package userdata

const (
	ScriptTemplate = `#!/bin/bash
# update cloud-init
tdnf install -y yum
tdnf makecache
yum update -y cloud-init 

yum install -y inotify-tools

# enable docker
systemctl enable docker
systemctl start docker

## Install Kubernetes packages
CNI_VERSION="v0.8.2"
mkdir -p /opt/cni/bin
curl -L "https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-amd64-${CNI_VERSION}.tgz" | tar -C /opt/cni/bin -xz

CRICTL_VERSION="v1.16.0"
mkdir -p /opt/bin
curl -L "https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz" | tar -C /opt/bin -xz

# RELEASE="$(curl -sSL https://dl.k8s.io/release/stable.txt)"
RELEASE=v1.17.0

mkdir -p /opt/bin
cd /opt/bin
curl -L --remote-name-all https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/amd64/{kubeadm,kubelet,kubectl}
chmod +x {kubeadm,kubelet,kubectl}

curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/${RELEASE}/build/debs/kubelet.service" | sed "s:/usr/bin:/opt/bin:g" > /etc/systemd/system/kubelet.service
mkdir -p /etc/systemd/system/kubelet.service.d
curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/${RELEASE}/build/debs/10-kubeadm.conf" | sed "s:/usr/bin:/opt/bin:g" > /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

systemctl enable --now kubelet

# path config
export PATH=${PATH}:/opt/bin
echo 'export PATH="${PATH}:/opt/bin"' >> ~/.bash_profile

# config DNS
sed -i -e "s/DNS=127.0.0.1/DNS=8.8.8.8/g" /etc/systemd/resolved.conf
systemctl restart systemd-resolved.service

cat <<"EOF" > /opt/startup.sh
#!/bin/bash
# watch bootstrap file created
inotifywait -e CREATE,MODIFY -m /root | while read line; do
  set $line
  filename=${3}
  cd /root
  if [ ${filename} = "bootstrap.cfg" ]; then
    # path config
    export PATH=${PATH}:/opt/bin
    # replace hostname
    replaced=replaced-${filename}
    sed "s/'{{ .default_hostname}}'/{{ .instance_id}}/" <(cat ${filename} | base64 -d) > ${replaced}
    cat ${replaced} > /etc/cloud/cloud.cfg.d/90-bootstrap.cfg
    # Run cloud-init
    cloud-init init --local
    cloud-init init
    cloud-init modules --mode=config
    cloud-init modules --mode=final

    if [[ ! -f /etc/kubernetes/admin.conf ]]; then
      rm -rf /var/lib/cloud
      cloud-init init --local
      cloud-init init
      cloud-init modules --mode=config
      cloud-init modules --mode=final
    fi
  fi
done
EOF

chmod 755 /opt/startup.sh

cat << "EOF" > /etc/systemd/system/startup.service
[Unit]
Description = watch bootstrap from cluster-api

[Service]
Type = simple
ExecStart = /opt/startup.sh
ExecStop=/usr/bin/kill -QUIT $MAINPID
LimitNOFILE=65536
 
[Install]
WantedBy = multi-user.target
EOF

systemctl enable startup
systemctl start startup
`
)
