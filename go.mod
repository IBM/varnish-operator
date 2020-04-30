module icm-varnish-k8s-operator

go 1.13

require (
	github.com/caarlos0/env/v6 v6.2.1
	github.com/docker/distribution v2.7.1+incompatible
	github.com/go-logr/zapr v0.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-cmp v0.3.1
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.7.0
	go.uber.org/zap v1.10.0
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	sigs.k8s.io/controller-runtime v0.5.2
)
