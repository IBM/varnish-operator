package pods

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &Predicate{}

type Predicate struct {
	selector labels.Selector
}

func NewPredicate(selectorSet map[string]string) (predicate.Predicate, error) {
	selector := labels.SelectorFromSet(selectorSet)
	podPredicate := &Predicate{
		selector: selector,
	}
	return podPredicate, nil
}

func (ep *Predicate) shared(meta metav1.Object) bool {
	if ep.selector.Matches(labels.Set(meta.GetLabels())) {
		return true
	}

	return false
}

func (ep *Predicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *Predicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *Predicate) Update(e event.UpdateEvent) bool {
	return ep.shared(e.MetaNew)
}

func (ep *Predicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta)
}
