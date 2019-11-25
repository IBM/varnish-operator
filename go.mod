module icm-varnish-k8s-operator

go 1.12

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff
	github.com/go-logr/zapr v0.1.0
	github.com/gogo/protobuf v1.1.1
	github.com/google/go-cmp v0.3.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pkg/errors v0.8.1
	go.uber.org/zap v1.9.1
	golang.org/x/net v0.0.0-20190926025831-c00fd9afed17 // indirect
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	sigs.k8s.io/controller-runtime v0.3.0
)
