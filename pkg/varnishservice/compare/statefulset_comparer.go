package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	statefulSetIgnoreFields = cmpopts.IgnoreFields(appsv1.StatefulSet{}, "Spec.Template.Spec.DeprecatedServiceAccount")
	compareQuantity         = cmp.Comparer(func(x, y resource.Quantity) bool { return x.Cmp(y) == 0 })
	deployOpts              = []cmp.Option{cmpopts.IgnoreFields(appsv1.StatefulSet{}, sharedIgnoreMetadata...), cmpopts.IgnoreFields(appsv1.StatefulSet{}, sharedIgnoreStatus...), statefulSetIgnoreFields, compareQuantity}
)

// EqualStatefulSet compares 2 statefulsets for equality
func EqualStatefulSet(found, desired *appsv1.StatefulSet) bool {
	return cmp.Equal(found, desired, deployOpts...)
}

// DiffStatefulSet generates a patch diff between 2 statefulsets
func DiffStatefulSet(found, desired *appsv1.StatefulSet) string {
	return cmp.Diff(found, desired, deployOpts...)
}
