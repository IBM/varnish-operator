package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	deployIgnoreFields = cmpopts.IgnoreFields(appsv1.Deployment{}, "Spec.Template.Spec.DeprecatedServiceAccount")
	compareQuantity    = cmp.Comparer(func(x, y resource.Quantity) bool { return x.Cmp(y) == 0 })
	deployOpts         = []cmp.Option{cmpopts.IgnoreFields(appsv1.Deployment{}, sharedIgnoreMetadata...), cmpopts.IgnoreFields(appsv1.Deployment{}, sharedIgnoreStatus...), deployIgnoreFields, compareQuantity}
)

func withDeploymentInheritance(desired, found *appsv1.Deployment) *appsv1.Deployment {
	var desiredCopy appsv1.Deployment
	desired.DeepCopyInto(&desiredCopy)
	if desiredCopy.Annotations == nil {
		desiredCopy.Annotations = make(map[string]string)
	}
	if desiredCopy.Annotations["deployment.kubernetes.io/revision"] == "" {
		desiredCopy.Annotations["deployment.kubernetes.io/revision"] = found.Annotations["deployment.kubernetes.io/revision"]
	}
	return &desiredCopy
}

// EqualDeployment compares 2 deployments for equality
func EqualDeployment(found, desired *appsv1.Deployment) bool {
	desiredCopy := withDeploymentInheritance(desired, found)
	return cmp.Equal(found, desiredCopy, deployOpts...)
}

// DiffDeployment generates a patch diff between 2 deployments
func DiffDeployment(found, desired *appsv1.Deployment) string {
	desiredCopy := withDeploymentInheritance(desired, found)
	return cmp.Diff(found, desiredCopy, deployOpts...)
}
