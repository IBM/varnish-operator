# Varnish operator debugging

## Setup
Setup KinD or KUBECONFIG for a kubernetes cluster.

NOTE: I'm not sure how pushing the images will work in kind...depending on the registry is setup it'll be hard to give instructions. Good news is if you're using already KinD, you probably already know how to push the images. :)

Be sure the following environment variables are set when you are building the code, docker images, or setting up the debug configuration in IntelliJ.
```bash
export NAMESPACE=default
export LOGLEVEL=debug
export LOGFORMAT=console
export VERSION=0.0.1-test
export REPO_PATH=dns.of.cr/namespace
export CONTAINER_IMAGE="$REPO_PATH/ibmcom/varnish-operator:$VERSION-dev"
export LEADERELECTION_ENABLED=false
export WEBHOOKS_ENABLED=false
export KUBECONFIG=xxx

# not necessary but for K* utils to run commands in parallel
export KPE_LABELS="-l varnish-component=varnish"
# allows for things like:
kpe -c varnish 'varnishadm vcl.list'
kpe -s -c varnish 'varnishstat -1 | grep SMA
```

## Build the operator and controller
Make sure your changes are building and update the CRD if necessary:
```bash
make all
```

## Build, tag, and push the docker images
```bash
# build the operator
make docker-build
# build the container images for the varnish pod 
make docker-build-pod
# tag and push the operator image
make docker-tag-push
# tag and push the varnish container images
make docker-tag-push-pod
```

## Apply the CRD
```bash
k apply -f varnish-operator/crds/varnishcluster.yaml
```

## Install the operator helm charts
In order to debug the operator, we still need to deploy the operator's helm chart so everything is in place for the operator. This includes services, service accounts, roles/rolebindings (cluster too), etc. Note that we're setting the `replicas` to 0 because we don't want the operator pod to be created.
```bash
helm upgrade vo varnish-operator -f config/samples/vo.yaml --install --set replicas=0
```

## Debug in IntelliJ
Go to [cmd/varnish-operator/main.go](cmd/varnish-operator/main.go), right click and select `Debug 'go build main.go` or `Debug 'go build github.com/ibm/varnish-operator/cmd/varnish-operator'` depending on where you click. From there, edit the build configuration and add the environment variables as described above. Hit `^D` and :rocket:!

## Install VarnishClusters
```bash
k apply -f config/samples/vc.yaml
```

## HAProxy specific testing
```bash
k apply -f config/samples/proxy-sidecar-config.yaml

# sleep 60
# check to see if the config has been updated across the cluster
kpe -s -c haproxy-sidecar 'cat /usr/local/etc/haproxy/haproxy.cfg'

# force the haproxy to reread the config
kpe -s -c haproxy-sidecar 'PID=$(pgrep -o haproxy) && kill -HUP $PID'
# or, the newly added script
kpe -s -c haproxy-sidecar -- /haproxy-scripts/haproxy-hup.sh
```

## Uninstall the operator helm charts
```bash
k delete vc varnishcluster-example
helm uninstall vo
```
