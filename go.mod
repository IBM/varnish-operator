module github.com/ibm/varnish-operator

go 1.15

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/docker/distribution v2.7.1+incompatible
	github.com/go-logr/zapr v0.4.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.6
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.17.0
	go.uber.org/zap v1.16.0
	golang.org/x/tools v0.1.3 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver v0.21.2 // indirect
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.8.3
)
