package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
)

var (
	podDisruptionBudgetOpts = []cmp.Option{cmpopts.IgnoreFields(policyv1beta1.PodDisruptionBudget{}, sharedIgnoreMetadata...), cmpopts.IgnoreFields(policyv1beta1.PodDisruptionBudget{}, sharedIgnoreStatus...)}
)

func EqualPodDisruptionBudget(found, desired *policyv1beta1.PodDisruptionBudget) bool {
	return cmp.Equal(found, desired, podDisruptionBudgetOpts...)
}

func DiffPodDisruptionBudget(found, desired *policyv1beta1.PodDisruptionBudget) string {
	return cmp.Diff(found, desired, podDisruptionBudgetOpts...)
}
