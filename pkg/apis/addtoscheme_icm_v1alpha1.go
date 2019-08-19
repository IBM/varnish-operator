package apis

import (
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kappsv1 "k8s.io/kubernetes/pkg/apis/apps/v1"
	kv1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme, RegisterBuiltInDefaults, v1alpha1.RegisterDefaults)
}

// RegisterBuiltInDefaults adds in necessary default functions from kubernetes library into the scheme for this project
func RegisterBuiltInDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&appsv1.StatefulSet{}, func(obj interface{}) { kappsv1.SetObjectDefaults_StatefulSet(obj.(*appsv1.StatefulSet)) })
	scheme.AddTypeDefaultingFunc(&v1.Service{}, func(obj interface{}) { kv1.SetObjectDefaults_Service(obj.(*v1.Service)) })
	return nil
}
