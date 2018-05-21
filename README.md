# Varnish As a Kubernetes Operator

## Getting Started With Development

Follow along with [this getting started guide](https://github.com/operator-framework/getting-started).

### Prerequisites

Development requires installation of the [`operator-sdk`](https://github.com/operator-framework/operator-sdk) tool. To install it:

1.  Have a Go environment [set up](http://sourabhbajaj.com/mac-setup/Go/README.html)

    * This primarily involves installing Go and setting a GOPATH

2.  Install [`dep`](https://github.com/golang/dep)

3.  `git clone` the [Github project for operator-sdk](https://github.com/operator-framework/operator-sdk) into the `$GOPATH/src/github.com/operator-framework` directory
4.  `git checkout tags/v0.0.5`
5.  `dep ensure`
6.  `go install github.com/operator-framework/operator-sdk/commands/operator-sdk`

yes I know. A lot of stuff...

### Building

**WARNING**: building requires docker to be running. HOWEVER, if you try to use docker associated with a minikube deployment, trying to build the image will blow up with a "permission denied" error, **and destroy your minikube deployment**. Use `docker-machine` or anything not-minikube to build this.

Part of the project involves code generation. The documentation is...thin on what gets generated and what doesn't, but it appears that the `operator/pkg/apis/icm/v1alpha1/types.go` file drives code generation.

What this means is that, on first run and any time you change the `types.go` file:

1.  `operator-sdk generate-k8s`.
2.  `operator-sdk build varnish-service:vX.X.X`

On subsequent runs:

1.  `operator-sdk build varnish-service:vX.X.X`
