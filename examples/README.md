## Requirements

### Environmental values

|                                       | 環境変数                                         |
|---------------------------------------|--------------------------------------------------|
| `NIFCLOUD_ACCESS_KEY`                 | APIアクセスキー                                  |
| `NIFCLOUD_SECRET_KEY`                 | APIシークレットキー                              |
| `NIFCLOUD_REGION`                     | リージョン                                       |
| `CLUSTER_API_SSH_KEY`                 | nifcloudに登録済みの公開鍵に対する秘密鍵のパス   |
| `CLUSTER_API_PRIVATE_KEY_PASS`        | `CLUSTER_API_SSH_KEY`のパスフレーズ              |

## Tools

| tool        | version |
|-------------|---------|
| kubectl     | v1.17.0 |
| kustomize   | v3.5.4  | 
| go          | 1.13.4  |
| kind        | v0.7.0  |

## Quick start

### ニフクラに公開鍵の登録

ニフクラのコントロールパネルからSSHキーを作成します。
取得した秘密鍵のPathとパスフレーズを環境変数に設定します。

```sh
export CLUSTER_API_SSH_KEY=<your-private-key-path>
export CLUSTER_API_PRIVATEKEY_PASS=<your-private-key-pass>

chmod 604 $CLUSTER_API_SSH_KEY 
```

### マニュフェストの作成

```sh
git clone # <<TODO insert repository url>>
./examples/generator.sh
```

### Managementクラスタを作成

ここでは[kind](https://github.com/kubernetes-sigs/kind)を使用します。

```sh
# クラスタの作成
kind create cluster --name=clusterapi

# Kubeconfigの設定
export KUBECONFIG="$(kind get kubeconfig-path --name="clusterapi")"
```
### CRDの登録

```sh
make install
```

### Providerのデプロイ

```sh
kubectl apply -f examples/_out/provider-components.yaml
```

### managerの起動

```sh
make run
```

### Clusterの作成
```sh
kubectl apply -f examples/_out/cluster.yaml
```

### Control Planeの作成
```sh
kubectl apply -f examples/_out/controlplane.yaml
```
サーバーの作成後、クラスタからkubeconfigを取得します。

```sh
kubectl get secret capi-kubeconfig -o jsonpath='{.data.value}' | base64 -d > kubeconfig
```

### アドオンのデプロイ

```sh
KUBECONFIG=./kubeconfig kubectl apply -f examples/_out/addons.yaml
```

### Nodeの作成

```sh
kubectl apply -f examples/_out/machines.yaml
```

