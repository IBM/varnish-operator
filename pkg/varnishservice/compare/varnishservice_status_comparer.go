package compare

import (
	icmv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"

	"github.com/google/go-cmp/cmp"
)

// EqualVarnishServiceStatus compares 2 statuses for equality
func EqualVarnishServiceStatus(old, current *icmv1alpha1.VarnishServiceStatus) bool {
	return cmp.Equal(old, current)
}

// DiffVarnishServiceStatus generates a patch diff between 2 statuses
func DiffVarnishServiceStatus(old, current *icmv1alpha1.VarnishServiceStatus) string {
	return cmp.Diff(old, current)
}
