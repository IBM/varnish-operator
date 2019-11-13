#!/bin/bash

#TODO add more checks and validations before running any commands

set -ex

varnish_namespace="varnish-operator-system"
pull_secret="docker-reg-secret"

if ! docker info; then
    echo -e "!!! Can't connect to Docker daemon, please make sure it is running"
    exit 1
fi

if ! which kind >/dev/null; then
    go get sigs.k8s.io/kind
fi

kind delete cluster > /dev/null 2>&1
kind create cluster

export KUBECONFIG=$(kind get kubeconfig-path)

kubectl create -f - << EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
EOF

helm init --service-account=tiller --upgrade --wait

kubectl create ns $varnish_namespace

token=$(bx cr token-get -q $(bx cr tokens | grep "for Varnish operator" | awk '{print $1}'))

kubectl create secret -n default \
  docker-registry $pull_secret --docker-server=us.icr.io --docker-username=token --docker-password=$token --docker-email=a@b.com

kubectl create secret -n $varnish_namespace \
  docker-registry $pull_secret --docker-server=us.icr.io --docker-username=token --docker-password=$token --docker-email=a@b.com

helm install --name=icm-varnish varnish-operator --wait

kubectl run nginx --image=nginx --replicas=2
kubectl rollout status deploy nginx --watch=true
kubectl expose deployment nginx --port=80 --target-port=80

kubectl apply -f - << EOF
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishCluster
metadata:
  labels:
    operator: varnish
  name: nginx
  namespace: default
spec:
  vcl:
    configMapName: vcl-file
    entrypointFileName: entrypoint.vcl
  replicas: 1
  varnish:
    resources:
      limits:
        cpu: 100m
        memory: "200Mi"
      requests:
        cpu: 100m
        memory: "200Mi"
    imagePullSecret: ${pull_secret}
  backend:
    selector:
      run: nginx
    port: 80
  service:
    port: 80
EOF

sleep 3; kubectl rollout status deploy $(kubectl get deploy -l operator=varnish -o jsonpath='{.items[*].metadata.name}') --watch

varnish=$(kubectl get po -l varnish-component=varnish -o jsonpath='{.items[*].metadata.name}')

sleep 10; kubectl logs $varnish
kubectl exec -it $varnish cat /etc/varnish/backends.vcl
kubectl port-forward service/$(kubectl get svc -l varnish-component=cache-service -o jsonpath='{.items[*].metadata.name}') 8080:80 &

sleep 10
curl --head http://127.0.0.1:8080/
curl --head http://127.0.0.1:8080/
kill %1
