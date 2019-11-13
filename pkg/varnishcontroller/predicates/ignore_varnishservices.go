package predicates

import (
	"icm-varnish-k8s-operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &ignoreVarnishClusterPredicate{}

type ignoreVarnishClusterPredicate struct {
}

func (ep *ignoreVarnishClusterPredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *ignoreVarnishClusterPredicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *ignoreVarnishClusterPredicate) Update(e event.UpdateEvent) bool {
	return ep.shared(e.MetaNew, e.ObjectNew)
}

func (ep *ignoreVarnishClusterPredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *ignoreVarnishClusterPredicate) shared(meta metav1.Object, obj runtime.Object) bool {
	if _, ok := obj.(*v1alpha1.VarnishCluster); ok {
		return false
	}

	return true
}

// NewIgnoreVarnishClusterPredicate filters out all VarnishClusters.
// Used as a stub for the controller builder that doesn't own a resource but still needs to provide something for the `For` function.
func NewIgnoreVarnishClusterPredicate() predicate.Predicate {
	return &ignoreVarnishClusterPredicate{}
}
