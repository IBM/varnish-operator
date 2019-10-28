module icm-varnish-k8s-operator

go 1.12

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/docker/distribution v0.0.0-20170726174610-edc3ab29cdff
	github.com/go-logr/zapr v0.1.0
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/google/go-cmp v0.3.0
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/pflag v1.0.3 // indirect
	go.uber.org/zap v1.9.1
	golang.org/x/net v0.0.0-20190926025831-c00fd9afed17 // indirect
	golang.org/x/sys v0.0.0-20190610200419-93c9922d18ae // indirect
	golang.org/x/text v0.3.2 // indirect
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.1
)
