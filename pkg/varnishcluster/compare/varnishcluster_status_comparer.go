package compare

import (
	vcapi "github.com/cin/varnish-operator/api/v1alpha1"

	"github.com/google/go-cmp/cmp"
)

// EqualVarnishClusterStatus compares 2 statuses for equality
func EqualVarnishClusterStatus(old, current *vcapi.VarnishClusterStatus) bool {
	return cmp.Equal(old, current)
}

// DiffVarnishClusterStatus generates a patch diff between 2 statuses
func DiffVarnishClusterStatus(old, current *vcapi.VarnishClusterStatus) string {
	return cmp.Diff(old, current)
}
