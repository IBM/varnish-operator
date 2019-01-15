package endpoints

import (
	"github.com/juju/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &EndpointPredicate{}

type EndpointPredicate struct {
	namespace string
	labels    map[string]string
}

func Predicate(selectorString string) (predicate.Predicate, error) {
	ls, err := labels.ConvertSelectorToLabelsMap(selectorString)
	if err != nil {
		return nil, errors.Annotate(err, "could not parse selector string")
	}
	ep := &EndpointPredicate{
		labels: ls,
	}
	return ep, nil
}

func (ep *EndpointPredicate) shared(meta metav1.Object) bool {
	inc := meta.GetLabels()
	for label, v := range ep.labels {
		if incV, found := inc[label]; !found || v != incV {
			return false
		}
	}
	return true
}

func (ep *EndpointPredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *EndpointPredicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *EndpointPredicate) Update(e event.UpdateEvent) bool {
	return ep.shared(e.MetaNew)
}

func (ep *EndpointPredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta)
}
