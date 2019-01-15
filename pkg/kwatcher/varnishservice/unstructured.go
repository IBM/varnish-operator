package varnishservice

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var UnstructuredVarnishService = &unstructured.Unstructured{}

// Init should ONLY be called by the `init()` function inside `config` package
func Init(group, version, kind string) {
	gvk := schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}
	UnstructuredVarnishService.SetGroupVersionKind(gvk)
}

func GetKind() string {
	return UnstructuredVarnishService.GetKind()
}

func GetAPIVersion() string {
	return UnstructuredVarnishService.GetAPIVersion()
}
