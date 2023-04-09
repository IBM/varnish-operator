package reconcile

import (
	"sync"
	"time"

	vcapi "github.com/cin/varnish-operator/api/v1alpha1"
	"github.com/cin/varnish-operator/pkg/logger"

	"sigs.k8s.io/controller-runtime/pkg/event"
)

func NewReconcileTriggerer(logr *logger.Logger, reconcileChan chan<- event.GenericEvent) *ReconcileTriggerer {
	return &ReconcileTriggerer{
		logger:        logr,
		timetable:     make(map[string]map[string]*time.Timer),
		reconcileChan: reconcileChan,
	}
}

type ReconcileTriggerer struct {
	logger        *logger.Logger
	timetable     map[string]map[string]*time.Timer
	reconcileChan chan<- event.GenericEvent
	sync.Mutex
}

func (q *ReconcileTriggerer) TriggerAfter(key string, triggerAfter time.Duration, instance *vcapi.VarnishCluster) {
	namespacedName := instance.Namespace + "/" + instance.Name
	q.logger.Debugf("Setting timer to trigger in %s", triggerAfter)
	timer := time.AfterFunc(triggerAfter, func() {
		q.Lock()
		delete(q.timetable, namespacedName)
		q.Unlock()
		q.logger.Debugf("Triggering reconcile")
		q.reconcileChan <- event.GenericEvent{
			Object: instance,
		}
	})
	q.Lock()
	defer q.Unlock()
	if timers, exists := q.timetable[namespacedName]; exists {
		if timer, exists := timers[key]; exists {
			timer.Stop()
			delete(q.timetable, namespacedName)
		}
	}

	if _, exists := q.timetable[namespacedName]; !exists {
		q.timetable[namespacedName] = map[string]*time.Timer{}
	}

	q.timetable[namespacedName][key] = timer
}

func (q *ReconcileTriggerer) Stop(key string, instance *vcapi.VarnishCluster) {
	namespacedName := instance.Namespace + "/" + instance.Name
	if timers, exists := q.timetable[namespacedName]; exists {
		if timer, exists := timers[key]; exists {
			timer.Stop()
			q.Lock()
			delete(q.timetable, namespacedName)
			q.Unlock()
		}
	}
}

func (q *ReconcileTriggerer) TimerExists(key string, service *vcapi.VarnishCluster) bool {
	namespacedName := service.Namespace + "/" + service.Name
	q.Lock()
	defer q.Unlock()
	if _, exists := q.timetable[namespacedName]; exists {
		if _, timerExists := q.timetable[namespacedName][key]; timerExists {
			return true
		}
	}

	return false
}
