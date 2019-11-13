package compare

import (
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"

	"github.com/google/go-cmp/cmp"
)

// EqualVarnishClusterStatus compares 2 statuses for equality
func EqualVarnishClusterStatus(old, current *icmv1alpha1.VarnishClusterStatus) bool {
	return cmp.Equal(old, current)
}

// DiffVarnishClusterStatus generates a patch diff between 2 statuses
func DiffVarnishClusterStatus(old, current *icmv1alpha1.VarnishClusterStatus) string {
	return cmp.Diff(old, current)
}
