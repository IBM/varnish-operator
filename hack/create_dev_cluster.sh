#!/bin/bash

set -ex
kube_version="1.17.0" #see https://hub.docker.com/r/kindest/node/tags for available Kubernetes version
if [ -n "${KUBERNETES_VERSION}" ]; then
  kube_version="${KUBERNETES_VERSION}"
fi

varnish_namespace="varnish-operator"
cluster_name="e2e-tests"

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

kind delete cluster --name $cluster_name > /dev/null 2>&1
kind create cluster --name $cluster_name --image kindest/node:v$kube_version --kubeconfig ./e2e-tests-kubeconfig

export KUBECONFIG=./e2e-tests-kubeconfig

kubectl create serviceaccount tiller -n kube-system
kubectl create clusterrolebinding tiller --serviceaccount=kube-system:tiller --clusterrole=cluster-admin

helm init --service-account=tiller --upgrade --wait

kubectl create ns $varnish_namespace

docker build -f Dockerfile  -t ibmcom/varnish-operator:local .
docker build -f Dockerfile.varnishd  -t ibmcom/varnish:local .
docker build -f Dockerfile.controller  -t ibmcom/varnish-controller:local .
docker build -f Dockerfile.exporter  -t ibmcom/varnish-metrics-exporter:local .

kind load --name $cluster_name docker-image ibmcom/varnish-operator:local
kind load --name $cluster_name docker-image ibmcom/varnish:local
kind load --name $cluster_name docker-image ibmcom/varnish-controller:local
kind load --name $cluster_name docker-image ibmcom/varnish-metrics-exporter:local

helm install --name=varnish-operator --namespace=$varnish_namespace varnish-operator --wait --set container.imagePullSecret=docker-reg-secret --set container.imagePullPolicy=Never --set container.image=ibmcom/varnish-operator:local
