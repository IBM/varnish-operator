package v1alpha1

import (
	"icm-varnish-k8s-operator/operator/controller/pkg/apis/icm"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemaGroupVersion is the group version used to register objects in this folder
var (
	version            = "v1alpha1"
	SchemeGroupVersion = schema.GroupVersion{Group: icm.GroupName, Version: version}
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder stays in k8s.io/kubernetes
	SchemeBuilder runtime.SchemeBuilder
	// localSchemeBuilder stays in k8s.io/kubernetes
	localSchemeBuilder = &SchemeBuilder
	// AddToScheme stays in k8s.io/kubernetes
	AddToScheme = localSchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&VarnishService{},
		&VarnishServiceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func init() {
	// only registry manually written functions here. Registration of generated functions takes place in generated files.
	// This separation makes the code compile even when generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}
