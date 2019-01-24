// Generate some defaults for VarnishService
//go:generate go run ../../../../vendor/k8s.io/code-generator/cmd/defaulter-gen/main.go -O zz_generated.defaults -i . -h ../../../../hack/boilerplate.go.txt

package v1alpha1

import (
	"fmt"
	"icm-varnish-k8s-operator/pkg/varnishservice/config"

	"k8s.io/api/core/v1"
	kv1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

var operatorCfg *config.Config

// Init is used to inject config to the package. Autogeneration requires to have functions with only one argument
// of specific type so that makes not possible to have structs where you can inject configs or other dependencies.
// Needs to be called before package usage.
func Init(cfg *config.Config) {
	operatorCfg = cfg
}

func SetDefaults_VarnishService(in *VarnishService) {
	if in.Spec.VCLConfigMap.Name == "" {
		in.Spec.VCLConfigMap.Name = fmt.Sprintf("%s-%s", in.Name, operatorCfg.DefaultVCLConfigMapName)
	}
}

func SetDefaults_VarnishVCLConfigMap(in *VarnishVCLConfigMap) {
	if in.BackendsFile == "" {
		in.BackendsFile = operatorCfg.DefaultBackendsFile
	}
	if in.DefaultFile == "" {
		in.DefaultFile = operatorCfg.DefaultDefaultFile
	}
	in.BackendsTmplFile = in.BackendsFile + ".tmpl"
}

func SetDefaults_VarnishDeployment(in *VarnishDeployment) {
	if in.VarnishMemory == "" {
		in.VarnishMemory = operatorCfg.DefaultVarnishMemory
	}
	if in.VarnishResources == nil {
		in.VarnishResources = &operatorCfg.DefaultVarnishResources
	}
	if in.LivenessProbe == nil {
		in.LivenessProbe = operatorCfg.DefaultLivenessProbe
	}
	if in.ReadinessProbe == nil {
		in.ReadinessProbe = &operatorCfg.DefaultReadinessProbe
	}
	if in.VarnishRestartPolicy == "" {
		in.VarnishRestartPolicy = operatorCfg.DefaultVarnishRestartPolicy
	}
}

func SetDefaults_VarnishDeploymentImage(in *VarnishDeploymentImage) {
	if in.Host == "" {
		in.Host = operatorCfg.VarnishImageHost
	}
	if in.Namespace == "" {
		in.Namespace = operatorCfg.VarnishImageNamespace
	}
	if in.Name == "" {
		in.Name = operatorCfg.VarnishImageName
	}
	if in.Tag == "" {
		in.Tag = operatorCfg.VarnishImageTag
	}
	if in.ImagePullPolicy == nil {
		in.ImagePullPolicy = &operatorCfg.VarnishImagePullPolicy
	}
	if in.ImagePullSecretName == "" {
		in.ImagePullSecretName = operatorCfg.ImagePullSecret
	}
}

func SetDefaults_ServiceSpec(in *v1.ServiceSpec) {
	s := &v1.Service{Spec: *in}
	kv1.SetObjectDefaults_Service(s)
	s.Spec.DeepCopyInto(in)
}
