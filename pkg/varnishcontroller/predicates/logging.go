package predicates

import (
	"icm-varnish-k8s-operator/pkg/logger"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &loggingPredicate{}

type loggingPredicate struct {
	logger *logger.Logger
}

func (ep *loggingPredicate) Create(e event.CreateEvent) bool {
	ep.logger.Debugf("Create event for resource %T: %s/%s", e.Object, e.Meta.GetNamespace(), e.Meta.GetName())
	return true
}

func (ep *loggingPredicate) Delete(e event.DeleteEvent) bool {
	ep.logger.Debugf("Delete event for resource %T: %s/%s", e.Object, e.Meta.GetNamespace(), e.Meta.GetName())
	return true
}

func (ep *loggingPredicate) Update(e event.UpdateEvent) bool {
	ep.logger.Debugf("Update event for resource %T: %s/%s", e.ObjectNew, e.MetaNew.GetNamespace(), e.MetaNew.GetName())
	return true
}

func (ep *loggingPredicate) Generic(e event.GenericEvent) bool {
	ep.logger.Debugf("Generic event for resource %T: %s/%s", e.Object, e.Meta.GetNamespace(), e.Meta.GetName())
	return true
}

// NewLoggingPredicate creates a predicate that doesn't filter out any requests, but logs them.
// Useful for debugging, when need to figure out what kind of requests the operator receives
func NewLoggingPredicate(logr *logger.Logger) predicate.Predicate {
	return &loggingPredicate{
		logger: logr,
	}
}
