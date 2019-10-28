package predicates

import (
	"icm-varnish-k8s-operator/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &ignorePodsPredicate{}

type ignorePodsPredicate struct {
}

func (ep *ignorePodsPredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *ignorePodsPredicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *ignorePodsPredicate) Update(e event.UpdateEvent) bool {
	return ep.shared(e.MetaNew, e.ObjectNew)
}

func (ep *ignorePodsPredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *ignorePodsPredicate) shared(meta metav1.Object, obj runtime.Object) bool {
	if _, ok := obj.(*v1alpha1.VarnishService); ok {
		return false
	}

	return true
}

// NewIgnoreVarnishServicesPredicate filters out all VarnishServices.
// Used as a stub for the controller builder that doesn't own a resource but still needs to provide something for the `For` function.
func NewIgnoreVarnishServicesPredicate() predicate.Predicate {
	return &ignorePodsPredicate{}
}
