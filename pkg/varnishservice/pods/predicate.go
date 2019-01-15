package pods

import (
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &AnnotationsPredicate{}

type AnnotationsPredicate struct {
	annotations map[string]string
}

func NewAnnotationsPredicate(annotations map[string]string) (predicate.Predicate, error) {
	if len(annotations) == 0 {
		return nil, errors.New("can't create predicate with nil or empty annotations")
	}
	podPredicate := &AnnotationsPredicate{
		annotations: annotations,
	}
	return podPredicate, nil
}

func (ep *AnnotationsPredicate) shared(meta metav1.Object) bool {
	for annotation, expected := range meta.GetAnnotations() {
		if actual, found := meta.GetAnnotations()[annotation]; !found || expected != actual {
			return false
		}
	}

	return true
}

func (ep *AnnotationsPredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *AnnotationsPredicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *AnnotationsPredicate) Update(e event.UpdateEvent) bool {
	return ep.shared(e.MetaNew)
}

func (ep *AnnotationsPredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta)
}
