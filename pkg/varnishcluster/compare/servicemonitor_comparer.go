package compare

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// EqualServiceMonitor compares 2 servicemonitors for equality
func EqualServiceMonitor(found, desired *unstructured.Unstructured) bool {
	return reflect.DeepEqual(found.GetLabels(), desired.GetLabels()) && reflect.DeepEqual(found.Object["spec"], desired.Object["spec"])
}

// DiffStatefulSet generates a patch diff between 2 servicemonitors
func DiffServiceMonitor(found, desired *unstructured.Unstructured) string {
	return cmp.Diff(found.Object["spec"], desired.Object["spec"]) + "\n" + cmp.Diff(found.GetLabels(), desired.GetLabels())
}
