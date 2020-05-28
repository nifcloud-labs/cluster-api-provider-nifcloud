#!/bin/bash

# Versions
CAPN_VERSION="v0.0.0"
CAPI_VERSION="v0.2.8"
CABPK_VERSION="v0.1.5"
CALICO_VERSION="v3.10"

# Nifcloud Settings
export SSH_KEY_NAME="${SSH_KEY_NAME:-default}"
export NIFCLOUD_REGION="$(echo ${NIFCLOUD_REGION:-jp-east-1} | tr -d '\n' | base64)"   # Tokyo
export NIFCLOUD_BASE64ENCODE_ACCESS_KEY=$(echo ${NIFCLOUD_ACCESS_KEY} | tr -d '\n' | base64)
export NIFCLOUD_BASE64ENCODE_SECRET_KEY=$(echo ${NIFCLOUD_SECRET_KEY} | tr -d '\n' | base64)

# Cluster Settings
export KUBERNETES_VERSION="${KUBERNETES_VERSION:-v1.17.0}"
export CLUSTER_NAME="${CLUSTER_NAME:-capi}"

# Machine Settings
# 55395: photon-3
export CONTROL_PLANE_IMAGE_ID="${CONTROL_PLANE_IMAGE_ID:-'55395'}"
export NODE_IMAGE_ID="${NODE_IMAGE_ID:-'55395'}"
# kubeadm required 2x CPU
export CONTROL_PLANE_INSTANCE_TYPE="${CONTROL_PLANE_INSTANCE_TYPE:-medium}"
export NODE_INSTANCE_TYPE="${NODE_INSTANCE_TYPE:-medium}"

# Output Settings
SOURCE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
OUTPUT_DIR=${OUTPUT_DIR:-${SOURCE_DIR}/_out}
echo "SOURCE: ${SOURCE_DIR}"
echo "OUTPUT: ${OUTPUT_DIR}"

COMPONENTS_CLUSTER_API_GENERATED_FILE=${SOURCE_DIR}/provider-components/provider-components-cluster-api.yaml
COMPONENTS_KUBEADM_GENERATED_FILE=${SOURCE_DIR}/provider-components/provider-components-kubeadm.yaml
COMPONENTS_NIFCLOUD_GENERATED_FILE=${SOURCE_DIR}/provider-components/provider-components-nifcloud.yaml

ADDON_CALICO_GENERATED_FILE=${SOURCE_DIR}/addons/calico.yaml
ADDON_NIFCLOUD_PROVIDER_GENERATED_FILE=${SOURCE_DIR}/addons/cloud-provider-nifcloud.yaml

CLUSTER_GENERATED_FILE=${OUTPUT_DIR}/cluster.yaml
CONTROLPLANE_GENERATED_FILE=${OUTPUT_DIR}/controlplane.yaml
MACHINES_GENERATED_FILE=${OUTPUT_DIR}/machines.yaml
PROVIDER_COMPONENTS_GENERATED_FILE=${OUTPUT_DIR}/provider-components.yaml
ADDON_GENERATED_FILE=${OUTPUT_DIR}/addons.yaml

if [ -d "$OUTPUT_DIR" ]; then
  echo "ERR: Folder ${OUTPUT_DIR} already exists. Delete it manually before running this script."
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"

# Generate cluster manifest
kustomize build "${SOURCE_DIR}/cluster" | envsubst > "${CLUSTER_GENERATED_FILE}"
echo "Generated ${CLUSTER_GENERATED_FILE}"

# Generate controlplane manifest
kustomize build "${SOURCE_DIR}/controlplane" | envsubst > "${CONTROLPLANE_GENERATED_FILE}"
echo "Generated ${CONTROLPLANE_GENERATED_FILE}"

# Generate machine manifest
kustomize build "${SOURCE_DIR}/machine" | envsubst > "${MACHINES_GENERATED_FILE}"
echo "Generated ${MACHINES_GENERATED_FILE}"

# Download & Generate provider-components.yaml
# Cluster API Provider Nifcloud
kustomize build "${SOURCE_DIR}/../config/default" | envsubst > "${COMPONENTS_NIFCLOUD_GENERATED_FILE}"
echo "Generated ${COMPONENTS_NIFCLOUD_GENERATED_FILE}"

## Cluster API
kustomize build "github.com/kubernetes-sigs/cluster-api/config/default/?ref=${CAPI_VERSION}" > "${COMPONENTS_CLUSTER_API_GENERATED_FILE}"
echo "Generated ${COMPONENTS_CLUSTER_API_GENERATED_FILE}"

## Cluster API Bootstrap Provider kubeadm
kustomize build "github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm/config/default/?ref=${CABPK_VERSION}" > "${COMPONENTS_KUBEADM_GENERATED_FILE}"
echo "Generated ${COMPONENTS_KUBEADM_GENERATED_FILE}"

# Download Network Plugin (Calico) manifest
curl -sL https://docs.projectcalico.org/${CALICO_VERSION}/manifests/calico.yaml -o "${ADDON_CALICO_GENERATED_FILE}"
echo "Downloaded ${ADDON_CALICO_GENERATED_FILE}"

# Generate a single provider components file.
kustomize build "${SOURCE_DIR}/provider-components" | envsubst > "${PROVIDER_COMPONENTS_GENERATED_FILE}"
echo "Generated ${PROVIDER_COMPONENTS_GENERATED_FILE}"
echo "⚠️ WARNING: ${PROVIDER_COMPONENTS_GENERATED_FILE} includes Nifcloud credentials"

# Generate a single addon components file.
kustomize build "${SOURCE_DIR}/addons" | envsubst > "${ADDON_GENERATED_FILE}"
echo "Generated ${ADDON_GENERATED_FILE}"
echo "⚠️ WARNING: ${ADDON_GENERATED_FILE} includes Nifcloud credentials"
