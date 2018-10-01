package apis

import (
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme, v1alpha1.RegisterDefaults)
}