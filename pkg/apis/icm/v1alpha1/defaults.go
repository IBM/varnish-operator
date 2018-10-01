// Generate some defaults for VarnishService
//go:generate go run ../../../../vendor/k8s.io/code-generator/cmd/defaulter-gen/main.go -O zz_generated.defaults -i . -h ../../../../hack/boilerplate.go.txt

package v1alpha1

import (
	"icm-varnish-k8s-operator/pkg/config"
	"k8s.io/api/core/v1"
	kv1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

var globalConf = config.GlobalConf

func SetDefaults_VarnishDeployment(in *VarnishDeployment) {

	if in.VarnishMemory == "" {
		in.VarnishMemory = globalConf.DefaultVarnishMemory
	}
	if in.VarnishResources == nil {
		in.VarnishResources = &globalConf.DefaultVarnishResources
	}
	if in.LivenessProbe == nil {
		in.LivenessProbe = globalConf.DefaultLivenessProbe
	}
	if in.ReadinessProbe == nil {
		in.ReadinessProbe = &globalConf.DefaultReadinessProbe
	}
	if in.ImagePullPolicy == nil {
		in.ImagePullPolicy = &globalConf.VarnishImagePullPolicy
	}
	if in.VarnishRestartPolicy == "" {
		in.VarnishRestartPolicy = globalConf.DefaultVarnishRestartPolicy
	}
	if in.BackendsFile == "" {
		in.BackendsFile = globalConf.DefaultBackendsFile
	}
	if in.DefaultFile == "" {
		in.DefaultFile = globalConf.DefaultDefaultFile
	}
	if in.VCLDir == "" {
		in.VCLDir = globalConf.VCLDir
	}
}

func SetDefaults_ServiceSpec(in *v1.ServiceSpec) {
	s := &v1.Service{Spec: *in}
	kv1.SetObjectDefaults_Service(s)
	s.Spec.DeepCopyInto(in)
}
