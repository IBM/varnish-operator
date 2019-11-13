package reconcile

import (
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"sync"
	"time"

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

func (q *ReconcileTriggerer) TriggerAfter(key string, triggerAfter time.Duration, instance *icmapiv1alpha1.VarnishCluster) {
	namespacedName := instance.Namespace + "/" + instance.Name
	q.logger.Debugf("Setting timer to trigger in %s", triggerAfter)
	timer := time.AfterFunc(triggerAfter, func() {
		q.Lock()
		delete(q.timetable, namespacedName)
		q.Unlock()
		q.logger.Debugf("Triggering reconcile")
		q.reconcileChan <- event.GenericEvent{
			Meta: instance,
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

func (q *ReconcileTriggerer) Stop(key string, instance *icmapiv1alpha1.VarnishCluster) {
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

func (q *ReconcileTriggerer) TimerExists(key string, service *icmapiv1alpha1.VarnishCluster) bool {
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
