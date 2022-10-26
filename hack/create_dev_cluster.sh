#!/bin/bash

set -ex

kube_version="1.25.2" # see https://github.com/kubernetes-sigs/kind/releases for available Kubernetes versions
if [ -n "${KUBERNETES_VERSION}" ]; then
  kube_version="${KUBERNETES_VERSION}"
fi

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

varnish_namespace="varnish-operator"
cluster_name="e2e-tests"
repo="ibmcom"
build_args="--build-arg VERSION=local"
platform="linux/amd64"
podman_in_use=false
ignore_podman=false
manage_cluster=true
use_buildx=false
create_vc=false
create_backends=false
skip_docker_build=flase
dry_run=false

function usage {
  cat << !
USAGE: $0 [-c cluster] [-n namespace] [-p platform] [-r repo] [-b] [-s] [-v] [-x]

Creates a dev cluster and varnish-operator install

-c|--cluster   | cluster
-n|--namespace | namespace
-p|--platform  | platform (not validated so know which build you're calling)
-r|--repo      | CR repository
-b             | create backends
-s             | skip docker build
-v             | create varnish cluster
-x             | ignore podman's presence
!
}

function default_vc_namespace {
  if [[ "$varnish_namespace" == "varnish-operator" ]]; then
    vc=$(kubectl get namespace --no-headers | grep varnish-cluster | wc -l)
    if [ "$vc" -eq 0 ]; then
      kubectl create namespace varnish-cluster
    fi
    varnish_namespace="varnish-cluster"
  fi
}

function create_nginx_backends {
  if [ "$dry_run" = false ]; then
    echo "dry-run: would otherwise be installing nginx"
    return 0
  fi
  default_vc_namespace
  kubectl create deployment nginx-backend --image nginx -n $varnish_namespace --port=80
}

function create_varnishcluster {
  if [ "$dry_run" = true ]; then
    echo "dry-run: would otherwise be installing varnishcluster"
    return 0
  fi

  default_vc_namespace
  cat <<EOF | kubectl create -f -
apiVersion: caching.ibm.com/v1alpha1
kind: VarnishCluster
metadata:
  name: varnishcluster-example
  namespace: $varnish_namespace
spec:
  vcl:
   configMapName: vcl-config
   entrypointFileName: entrypoint.vcl
  backend:
    port: 80
    selector:
      app: nginx-backend
  service:
    port: 80 # Varnish pods will be exposed on that port
EOF
}

while (( "$#" )); do
  opt="$1"; shift;
  case "$opt" in
    "-b"|"--backends") create_backends=true;;
    "-d") dry_run=true;;
    "-s") skip_docker_build=true;;
    "-v") create_vc=true;;
    "-x") ignore_podman=true;;
    "-c"|"--cluster") cluster_name="$1"; manage_cluster=false; shift;;
    "-n"|"--namespace") varnish_namespace="$1"; shift;;
    "-p"|"--platform") platform="$1"; shift;;
    "-r"|"--repo") repo="$1"; shift;;
    *) echo "invalid option: \""$opt"\"" >&2; usage; exit 1;;
  esac
done

if [ "$create_vc" = true ]; then
  create_varnishcluster
  exit 0
fi

if [ "$create_backends" = true ]; then
  create_nginx_backends
  exit 0
fi

container_image="$repo/varnish-operator:local"
images=($container_image $repo/varnish:local $repo/varnish-controller:local $repo/varnish-metrics-exporter:local)
dockerfiles=(Dockerfile Dockerfile.varnishd Dockerfile.controller Dockerfile.exporter)

if [[ $platform =~ ^.*,.*$ ]]; then
  use_buildx=true
  build_args="buildx build $build_args --platform $platform"
else
  build_args="build $build_args --platform $platform"
fi

if [ "$ignore_podman" = false ] && which podman >/dev/null; then
  podman_in_use=true
elif [ "$use_buildx" == true ]; then
  build_args="$build_args --push"
fi

if [ "$dry_run" = true ]; then
  echo "build_args: $build_args, manage_cluster: $manage_cluster; use_buildx: $use_buildx, podman_in_use: $podman_in_use, ignore_podman: $ignore_podman"
  exit 0
fi

if [ "$manage_cluster" = true ]; then
  kind delete cluster --name $cluster_name > /dev/null 2>&1
  kind create cluster --name $cluster_name --image kindest/node:v$kube_version --kubeconfig ./e2e-tests-kubeconfig
  export KUBECONFIG=./e2e-tests-kubeconfig
fi

vo=$(kubectl get namespace --no-headers | grep varnish-operator | wc -l)
if [ "$vo" -eq 0 ]; then
  kubectl create ns $varnish_namespace
fi

# if skipping docker build, ensure all the images are at least in the local docker registry
if [ "$skip_docker_build" = true]; then
  set +e
  for image in "${images[@]}"; do
    res=$(docker image inspect $image)
    if [ "$?" -ne 0 ]; then
      echo "missing $image. cannot skip docker build"
      skip_docker_build=false
    fi
  done
  set -e
fi

if [ "$skip_docker_build" = false ]; then
  for ((i=0; i<${#images[@]}; i++)); do
    echo "docker $build_args -f ${dockerfiles[$i]} -t ${images[$i]} ."
  done

  if [ "$use_buildx" = true ] && [ "$podman_in_use" = false ]; then
    for image in "${images[@]}"; do
      docker pull $image
    done
  fi
fi

for image in "${images[@]}"; do
  kind load docker-image -n $cluster_name image
done

helm install varnish-operator varnish-operator --namespace=$varnish_namespace --wait --set container.imagePullPolicy=Never --set container.image=$container_image
