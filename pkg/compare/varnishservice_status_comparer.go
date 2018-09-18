package compare

import (
	"github.com/google/go-cmp/cmp"
	icmv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"reflect"
)

func EqualVarnishServiceStatus(old, current *icmv1alpha1.VarnishServiceStatus) bool {
	return reflect.DeepEqual(old, current)
}

func DiffVarnishServiceStatus(old, current *icmv1alpha1.VarnishServiceStatus) string {
	return cmp.Diff(old, current)
}
