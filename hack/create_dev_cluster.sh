#!/bin/bash

set -ex
kube_version="1.25.2" #see https://hub.docker.com/r/kindest/node/tags for available Kubernetes version
if [ -n "${KUBERNETES_VERSION}" ]; then
  kube_version="${KUBERNETES_VERSION}"
fi

varnish_namespace="varnish-operator"
cluster_name="e2e-tests"
repo="ibmcom"
container_image="$repo/varnish-operator:local"
build_args="--build-arg VERSION=local --platform linux/arm64,linux/amd64 --push"
podman_in_use=false

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

if which podman >/dev/null; then
  img=$(podman images | grep '$repo/varnish-operator' | grep local | cut -f1 -d' ')
  if [ -n "$img" ]; then
    podman_in_use=true
    echo "using podman image: $img"
    # podman doesn't support the `--push` option. hopefully no other options start with `--push`...
    build_args=$(echo ${build_args/--push})
  else
    echo "podman installed but varnish-operator image not found"
  fi
fi

kind delete cluster --name $cluster_name > /dev/null 2>&1
kind create cluster --name $cluster_name --image kindest/node:v$kube_version --kubeconfig ./e2e-tests-kubeconfig

export KUBECONFIG=./e2e-tests-kubeconfig

kubectl create ns $varnish_namespace

docker buildx build $build_args -f Dockerfile -t $container_image .
docker buildx build $build_args -f Dockerfile.varnishd -t $repo/varnish:local .
docker buildx build $build_args -f Dockerfile.controller -t $repo/varnish-controller:local .
docker buildx build $build_args -f Dockerfile.exporter -t $repo/varnish-metrics-exporter:local .

if [ "$podman_in_use" = false ]; then
  docker pull $container_image
  docker pull $repo/varnish:local
  docker pull $repo/varnish-controller:local
  docker pull $repo/varnish-metrics-exporter:local
fi

kind load docker-image -n $cluster_name $container_image
kind load docker-image -n $cluster_name $repo/varnish:local
kind load docker-image -n $cluster_name $repo/varnish-controller:local
kind load docker-image -n $cluster_name $repo/varnish-metrics-exporter:local

helm install vo varnish-operator --namespace=$varnish_namespace --wait --set container.imagePullPolicy=Never --set container.image=$container_image
