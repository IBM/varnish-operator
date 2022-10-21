#!/bin/bash

set -ex
kube_version="1.25.2" #see https://hub.docker.com/r/kindest/node/tags for available Kubernetes version
if [ -n "${KUBERNETES_VERSION}" ]; then
  kube_version="${KUBERNETES_VERSION}"
fi

varnish_namespace="varnish-operator"
cluster_name="e2e-tests"
container_image="ibmcom/varnish-operator:local"

if ! which docker; then
    echo -e "Install docker first"
    exit 1
fi

if ! which kind >/dev/null; then
    echo -e "Install kind first"
    exit 1
fi

if ! which helm >/dev/null; then
    echo -e "Install helm first"
    exit 1
fi

#if which podman >/dev/null; then
#  img=$(podman images | grep 'ibmcom/varnish-operator' | grep local | cut -f1 -d' ')
#  if [ -n "$img" ]; then
#    echo "using podman image: $img"
#    container_image=$img
#  else
#    echo "podman installed but varnish-operator image not found. using default container_image"
#  fi
#fi

kind delete cluster --name $cluster_name > /dev/null 2>&1
kind create cluster --name $cluster_name --image kindest/node:v$kube_version --kubeconfig ./e2e-tests-kubeconfig

export KUBECONFIG=./e2e-tests-kubeconfig

kubectl create ns $varnish_namespace

docker buildx build --build-arg VERSION=local --platform linux/arm64,linux/amd64 --push -f Dockerfile -t ibmcom/varnish-operator:local .
docker buildx build --build-arg VERSION=local --platform linux/arm64,linux/amd64 --push -f Dockerfile.varnishd -t ibmcom/varnish:local .
docker buildx build --build-arg VERSION=local --platform linux/arm64,linux/amd64 --push -f Dockerfile.controller -t ibmcom/varnish-controller:local .
docker buildx build --build-arg VERSION=local --platform linux/arm64,linux/amd64 --push -f Dockerfile.exporter -t ibmcom/varnish-metrics-exporter:local .

docker pull ibmcom/varnish-operator:local
docker pull ibmcom/varnish:local
docker pull ibmcom/varnish-controller:local
docker pull ibmcom/varnish-metrics-exporter:local

kind load docker-image -n $cluster_name ibmcom/varnish-operator:local
kind load docker-image -n $cluster_name ibmcom/varnish:local
kind load docker-image -n $cluster_name ibmcom/varnish-controller:local
kind load docker-image -n $cluster_name ibmcom/varnish-metrics-exporter:local

helm install vo varnish-operator --namespace=$varnish_namespace --wait --set container.imagePullPolicy=Never --set container.image=$container_image
