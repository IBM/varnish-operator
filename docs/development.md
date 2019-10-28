# Development

Requirements:

* Kubernetes 1.12 or newer. You can use [minikube](https://kubernetes.io/docs/setup/minikube/) for local development)
* Go 1.12 with enabled go modules
* [Kubebuilder](https://kubebuilder.io/quick-start.html#installation) 2.0.0+
* [kustomize](https://github.com/kubernetes-sigs/kustomize) 3.1.0+
* [helm](https://helm.sh/)
* [docker](https://docs.docker.com/install/)
* [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
* [GolangCI-Lint](https://github.com/golangci/golangci-lint) 1.19.1+
* [gitbook cli](https://github.com/GitbookIO/gitbook-cli), if you want to modify and test the docs locally

### Code structure

The project consists of 2 components working together:

* Varnish operator itself, that manages `VarnishService` deployments
* Kwatcher is a process that's running in the same container as Varnish. It is responsible for watching Kubernetes resources and reacting accordingly. For example, Kwatcher reloads the VCL configuration when backends scale or the VCL configuration has changed in the ConfigMap.
                                                                              
Both components live in one repo and share the same codebase, dependencies, build scripts, etc.

The operator and kwatcher's codebases are located in `/pkg/varnishservice/` and `pkg/kwatcher` folders respectively.
The main packages for both components can be found in the `cmd/` folder.

### Developing the operator
```bash
$ git clone git@github.ibm.com:TheWeatherCompany/icm-varnish-k8s-operator.git
$ cd icm-varnish-k8s-operator
$ go mod download
```

#### Run the operator locally against a Kubernetes cluster
The operator can be run locally without building the image every time the code changes.

First, you need to install the [CRD](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/) for `VarnishService` resource.

```bash
$ make install
<path>/<to>/<controller-gen>/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build <path>/<to>/<repo>/icm-varnish-k8s-operator/config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/varnishservices.icm.ibm.com created
```

You should see the created CRD in your cluster:

```bash
$ kubectl get customresourcedefinitions.apiextensions.k8s.io
 NAME                          CREATED AT
 varnishservices.icm.ibm.com   2019-06-05T09:53:26Z
```

`make install` should be run only for the first time and after changes in the CRD schema because it is only responsible for installing and updating the CRD for the `VarnishService` resource.

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
helm install --name varnish-operator --namespace varnish-operator-system --set container.image=us.icr.io/icm-varnish/varnish-controller:test ./varnish-operator
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

### Developing kwatcher

Varnish pods (with kwatcher inside) can only be tested by running in Kubernetes. That means that we need to build an image every time we make a change in the code related to kwatcher.

After changes are made in the code, the typical workflow will be the following:

```bash
docker build -f Dockerfile.Varnish  -t <image-name> .
docker push <image-name>
```

Then, in your `VarnishService`, specify your image:

```yaml
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishService
...
spec:
  statefulSet:
    container:
      image: <image-name>
...
```
The StatefulSet will reload the pods with new image. If you're reusing the same image name, make sure `spec.statefulSet.container.imagePullPolicy` is set to `Always` and reload the pods manually by deleting them or recreating the `VarnishService`. 

For images uploaded to a private registry, [create an image pull secret](https://pages.github.ibm.com/TheWeatherCompany/icm-docs/managed-kubernetes/container-registry.html#creating-an-image-pull-secret) and set the name of it in the `spec.container.imagePullSecret` field. 

### Tests

To run tests simply run `make test` in repo's root directory. 

Tests depend on binaries provided by Kubebuilder so it has to be [installed](https://kubebuilder.io/quick-start.html#installation) first.