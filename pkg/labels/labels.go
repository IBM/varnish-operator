package labels

import (
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
)

// CombinedComponentLabels create labels for a component and inherits VarnishCluster object labels
func CombinedComponentLabels(instance *vcapi.VarnishCluster, componentName string) (m map[string]string) {
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
func ComponentLabels(instance *vcapi.VarnishCluster, componentName string) map[string]string {
	return map[string]string{
		vcapi.LabelVarnishOwner:     instance.Name,
		vcapi.LabelVarnishComponent: componentName,
		vcapi.LabelVarnishUID:       string(instance.UID),
	}
}

func InheritLabels(instance *vcapi.VarnishCluster) (m map[string]string) {
	m = make(map[string]string, len(instance.Labels))
	for k, v := range instance.Labels {
		m[k] = v
	}
	return
}
