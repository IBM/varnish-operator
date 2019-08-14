package configmaps

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &ConfigMapPredicate{}

type ConfigMapPredicate struct {
	name string
}

func Predicate(name string) predicate.Predicate {
	return &ConfigMapPredicate{
		name: name,
	}
}

func (cmp *ConfigMapPredicate) shared(meta metav1.Object) bool {
	return meta.GetName() == cmp.name
}

func (cmp *ConfigMapPredicate) Create(e event.CreateEvent) bool {
	return cmp.shared(e.Meta)
}

func (cmp *ConfigMapPredicate) Delete(e event.DeleteEvent) bool {
	return cmp.shared(e.Meta)
}

func (cmp *ConfigMapPredicate) Update(e event.UpdateEvent) bool {
	return cmp.shared(e.MetaNew)
}

func (cmp *ConfigMapPredicate) Generic(e event.GenericEvent) bool {
	return cmp.shared(e.Meta)
}
