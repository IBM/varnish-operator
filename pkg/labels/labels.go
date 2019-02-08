package labels

import (
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
)

// CombinedComponentLabels create labels for a component and inherits VarnishService object labels
func CombinedComponentLabels(instance *icmapiv1alpha1.VarnishService, componentName string) (m map[string]string) {
	inherited := InheritLabels(instance)
	generated := ComponentLabels(instance, componentName)

	m = make(map[string]string, len(inherited)+len(generated))
	for k, v := range inherited {
		m[k] = v
	}
	for k, v := range generated {
		m[k] = v
	}
	return
}

// CombinedComponentLabels create labels for a component
func ComponentLabels(instance *icmapiv1alpha1.VarnishService, componentName string) map[string]string {
	return map[string]string{
		icmapiv1alpha1.LabelVarnishOwner:     instance.Name,
		icmapiv1alpha1.LabelVarnishComponent: componentName,
		icmapiv1alpha1.LabelVarnishUID:       string(instance.UID),
	}
}

func InheritLabels(instance *icmapiv1alpha1.VarnishService) (m map[string]string) {
	m = make(map[string]string, len(instance.Labels))
	for k, v := range instance.Labels {
		m[k] = v
	}
	return
}
