// Generate some defaults for VarnishService
//go:generate go run ../../../../vendor/k8s.io/code-generator/cmd/defaulter-gen/main.go -O zz_generated.defaults -i . -h ../../../../hack/boilerplate.go.txt

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/util/intstr"

	v1 "k8s.io/api/core/v1"
	kv1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

// SetDefaults_VarnishService sets defaults for everything inside VarnishService.
// Normally, this would be handled by the generator code, but in this case, the order that things are defaulted matters,
// since some defaults depend on other values already being defaulted.
// Therefore, add the `defaulter-gen=covers` flag which tells the generator to ignore any recursive `SetDefaults` functions,
// meaning they must be included here instead. With this, the order of defaults can be controlled
// +k8s:defaulter-gen=covers
func SetDefaults_VarnishService(in *VarnishService) {
	SetDefaults_VarnishServiceSpec(&in.Spec)
	SetDefaults_VarnishContainer(&in.Spec.Deployment.Container)
	SetDefaults_VarnishServiceService(&in.Spec.Service)
}

func SetDefaults_VarnishServiceSpec(in *VarnishServiceSpec) {
	if in.LogLevel == "" {
		in.LogLevel = "info"
	}
	if in.LogFormat == "" {
		in.LogFormat = "json"
	}
}

func SetDefaults_VarnishContainer(in *VarnishContainer) {
	if in.ImagePullPolicy == "" {
		in.ImagePullPolicy = v1.PullAlways
	}
	if in.RestartPolicy == "" {
		in.RestartPolicy = v1.RestartPolicyAlways
	}
}

func SetDefaults_VarnishServiceService(in *VarnishServiceService) {
	s := &v1.Service{Spec: in.ServiceSpec}
	kv1.SetObjectDefaults_Service(s)
	s.Spec.DeepCopyInto(&in.ServiceSpec)

	if in.VarnishPort.Name == "" {
		in.VarnishPort.Name = "varnish"
	}
	if in.VarnishPort.TargetPort == (intstr.IntOrString{}) {
		in.VarnishPort.TargetPort = intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: in.VarnishPort.Port,
		}
	}
	if in.VarnishExporterPort.Name == "" {
		in.VarnishExporterPort.Name = "varnishexporter"
	}
	//we don't support running exporter on custom port yet so ignore the value set by user for now
	in.VarnishExporterPort.TargetPort = intstr.FromInt(VarnishPrometheusExporterPort)
}
