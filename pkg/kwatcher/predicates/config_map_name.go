package predicates

import (
	"icm-varnish-k8s-operator/pkg/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &configMapNamePredicate{}

type configMapNamePredicate struct {
	configMapName string
	log           *logger.Logger
}

func (ep *configMapNamePredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *configMapNamePredicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *configMapNamePredicate) Update(e event.UpdateEvent) bool {
	return ep.shared(e.MetaNew, e.ObjectNew)
}

func (ep *configMapNamePredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta, e.Object)
}

func (ep *configMapNamePredicate) shared(meta metav1.Object, obj runtime.Object) bool {
	if _, ok := obj.(*v1.ConfigMap); !ok {
		return true
	}

	if meta.GetName() != ep.configMapName {
		return false
	}

	return true
}

// NewConfigMapNamePredicate filters out config maps that doesn't match the name specified
func NewConfigMapNamePredicate(configMapName string, logr *logger.Logger) predicate.Predicate {
	return &configMapNamePredicate{
		configMapName: configMapName,
		log:           logr,
	}
}
