# Development

Requirements:

* Kubernetes 1.12 or newer. You can use [minikube](https://kubernetes.io/docs/setup/minikube/) for local development.
* Go 1.12+ with enabled go modules
* [Kubebuilder](https://kubebuilder.io/quick-start.html#installation) 2.0.0+
* [kustomize](https://github.com/kubernetes-sigs/kustomize) 3.1.0+
* [helm](https://helm.sh/) v2.14.3+
* [docker](https://docs.docker.com/install/)
* [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
* [GolangCI-Lint](https://github.com/golangci/golangci-lint) 1.19.1+
* [gitbook cli](https://github.com/GitbookIO/gitbook-cli), if you want to modify and test the docs locally
* [kind](https://github.com/kubernetes-sigs/kind) v0.6.0+ for running end to end tests

### Code structure

The project consists of 2 components working together:

* Varnish operator itself, that manages `VarnishCluster` resources
* Varnish Controller is a process that's running in the same container as Varnish. It is responsible for watching Kubernetes resources and reacting accordingly. For example, Varnish Controller reloads the VCL configuration when backends scale or the VCL configuration has changed in the ConfigMap.
                                                                              
Both components live in one repo and share the same codebase, dependencies, build scripts, etc.

The operator and varnish controller's codebases are located in `/pkg/varnishcluster/` and `pkg/varnishcontroller` folders respectively.
The main packages for both components can be found in the `cmd/` folder.

### Developing the operator
```bash
$ git clone git@github.ibm.com:TheWeatherCompany/icm-varnish-k8s-operator.git
$ cd icm-varnish-k8s-operator
$ go mod download
```

#### Run the operator locally against a Kubernetes cluster
The operator can be run locally without building the image every time the code changes.

First, you need to install the [CRD](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/) for `VarnishCluster` resource.

```bash
$ make install
<path>/<to>/<controller-gen>/controller-gen "crd:trivialVersions=true" rbac:roleName=varnish-operator paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build <path>/<to>/<repo>/config/crd > <path>/<to>/<repo>/varnish-operator/templates/customresourcedefinition.yaml
<path>/<to>/<controller-gen>/controller-gen "crd:trivialVersions=true" rbac:roleName=varnish-operator paths="./..." output:crd:none output:rbac:stdout > <path>/<to>/<repo>/varnish-operator/templates/clusterrole.yaml
kustomize build <path>/<to>/<repo>/config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/varnishclusters.icm.ibm.com created
```

You should see the created CRD in your cluster:

```bash
$ kubectl get customresourcedefinitions.apiextensions.k8s.io
 NAME                          CREATED AT
 varnishclusters.icm.ibm.com   2019-06-05T09:53:26Z
```

`make install` should be run only for the first time and after changes in the CRD schema because it is only responsible for installing and updating the CRD for the `VarnishCluster` resource.

After that you're ready to run the operator:

 
```bash
$ make run
make run                                                                                                              
<path>/<to>/<controller-gen>/controller-gen object:headerFile=./hack/boilerplate.go.txt paths="./..."
cd <path>/<to>/<repo>/icm-varnish-k8s-operator/ && go generate ./pkg/... ./cmd/...
cd <path>/<to>/<repo>/icm-varnish-k8s-operator/ && goimports -w ./pkg ./cmd
cd <path>/<to>/<repo>/icm-varnish-k8s-operator/ && go vet ./pkg/... ./cmd/...
NAMESPACE="default" LOGLEVEL=debug LOGFORMAT=console CONTAINER_IMAGE=us.icr.io/icm-varnish/varnish:0.20.0-dev LEADERELECTION_ENABLED=false WEBHOOKS_ENABLED=false go run <path>/<to>/<repo>/icm-varnish-k8s-operator/cmd/manager/main.go...
```

By default the operator will work in the `default` namespace. You can override that behaviour by setting the `NAMESPACE` env var:

```bash
$ NAMESPACE=varnish-operator make run
```

After you've made changes to the operator code, just rerun `make run` to test them.

#### Deploying your operator in a Kubernetes cluster
Some changes can't be tested by running the operator locally:

* Validating and Mutating webhooks.
* Correctness of RBAC permissions. With the `make run` approach, your local `kubectl` configs are used to communicate with the cluster, not the clusterrole as it would be in production.

To test that functionality, you would have to run your operator as a pod in the cluster.

This can be done using the helm template configured to use your custom image:

```bash
docker build -t <image-name> -f Dockerfile .
docker push <image-name>
make manifests #make sure your helm charts are in sync with current CRD and RBAC definitions
helm install --name varnish-operator --namespace varnish-operator-system --set container.image=<image-name> ./varnish-operator
``` 

If your docker image is located in a private container registry, you'll need to [create an image pull secret](https://pages.github.ibm.com/TheWeatherCompany/icm-docs/managed-kubernetes/container-registry.html#creating-an-image-pull-secret) and reference it by adding `--set container.imagePullSecret=<image-pull-secret>` to the `helm install` command.

Check the operator pod logs to make sure all works as expected:

```bash
kubectl logs -n varnish-operator-system varnish-operator-fd96f48f-gn6mc
{"level":"info","ts":1559818986.866487,"caller":"manager/main.go:34","msg":"Version: 0.14.5"}
{"level":"info","ts":1559818986.8665597,"caller":"manager/main.go:35","msg":"Leader election enabled: true"}
{"level":"info","ts":1559818986.866619,"caller":"manager/main.go:36","msg":"Log level: info"}
...
```

### Developing varnish controller

Varnish pods (with varnish-controller inside) can only be tested by running in Kubernetes. That means that we need to build an image every time we make a change in the code related to varnish-controller.

After changes are made in the code, the typical workflow will be the following:

```bash
docker build -f Dockerfile.controller  -t <image-name> .
docker push <image-name>
```

Then, in your `VarnishCluster`, specify your image:

```yaml
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishCluster
...
spec:
  varnish:
...
    controller:
      image: <image-name>
...
```

The StatefulSet will reload the pods with new image. If you're reusing the same image name, make sure `spec.statefulSet.container.imagePullPolicy` is set to `Always` and reload the pods manually by deleting them or recreating the `VarnishCluster`. 

For images uploaded to a private registry, [create an image pull secret](https://pages.github.ibm.com/TheWeatherCompany/icm-docs/managed-kubernetes/container-registry.html#creating-an-image-pull-secret) and set the name of it in the `spec.varnish.imagePullSecret` field.

To change varnishd - varnish daemon or varnish metrics exporter containers refer to appropriate Docker files configurations. Update `VarnishCluster` and specify your images same way as it is described for varnish controller above.

To build new varnishd use:

```bash
docker build -f Dockerfile.varnish -t <image-name> .
docker push <image-name>
```

Then, in your `VarnishCluster`, specify your image for varnishd:

```yaml
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishCluster
...
spec:
  varnish:
    image: <image-name>
...
```

To build new varnish metrics exporter use:

```bash
docker build -f Dockerfile.exporter -t <image-name> .
docker push <image-name>
```

Then, in your `VarnishCluster`, specify your image for varnish metrics exporter:

```yaml
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishCluster
...
spec:
  varnish:
...
    metricsExporter:
      image: <image-name>
...
```

There is `PROMETHEUS_VARNISH_EXPORTER_VERSION` docker's build argument available. This allows you to specify which version of prometheus varnish metrics exporter to use in your image. Just set it to the required value before build the metrics exporter image:

`docker build --build-arg PROMETHEUS_VARNISH_EXPORTER_VERSION=1.5.2  -t <image-name> -f Dockerfile.exporter .`

### Tests

To run tests simply run `make test` in repo's root directory. 

Tests depend on binaries provided by Kubebuilder so it has to be [installed](https://kubebuilder.io/quick-start.html#installation) first.

### End to End tests

To run end-to-end tests you have to have [kind](https://github.com/kubernetes-sigs/kind) installed first. Then, simply run `make e2e-tests` in the root directory.

It will do the following:
 1. Bring up a kubernetes cluster using `kind`. You can also set the Kubernetes version you want to use by setting the `KUBERNETES_VERSION` env var (e.g. `KUBERNETES_VERSION=1.6.3 make e2e-tests`)
 1. Build all docker images and load them into the cluster
 1. Install the operator using the local Helm chart. It will set the image pull policy to `Never` to be able to use the built images in previous step.
 1. Run tests written in Go. The tests also rely on the docker images built and loaded into the cluster. It requires you to remember to use the local docker images, but it also facilitates better portability.
 1. Delete the cluster
