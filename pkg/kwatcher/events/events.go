package events

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
)

const (
	EventRecorderName = "varnish"

	EventReasonReloadError         EventReason = "ReloadError"
	EventReasonVCLCompilationError EventReason = "VCLCompilationError"
)

var (
	eventAnnotations map[string]string
	VSObject         *v1.ObjectReference //an object reference used to send events on behalf of VarnishService resource
)

// Init should ONLY be called by the `init()` function inside `config` package
func Init(sourcePod, kind, namespace, name, apiVersion string, uid types.UID) {
	eventAnnotations = map[string]string{
		"sourcePod": sourcePod,
	}
	VSObject = &v1.ObjectReference{
		Kind:       kind,
		Namespace:  namespace,
		Name:       name,
		APIVersion: apiVersion,
		UID:        uid,
	}
}

// EventReason is the reason why the event was create. The value appears in the 'Reason' tab of the events list
type EventReason string

// NewEventHandler creates a new event handler that will use the specified recorder
func NewEventHandler(recorder record.EventRecorder) *EventHandler {
	return &EventHandler{
		Recorder: recorder,
	}
}

// EventHandler handles the operations for events
type EventHandler struct {
	Recorder record.EventRecorder
}

// Warning creates a 'warning' type event
func (e *EventHandler) Warning(object runtime.Object, reason EventReason, message string) {
	e.Recorder.AnnotatedEventf(object, eventAnnotations, v1.EventTypeWarning, string(reason), message)
}

// Normal creates a 'normal' type event
func (e *EventHandler) Normal(object runtime.Object, reason EventReason, message string) {
	e.Recorder.AnnotatedEventf(object, eventAnnotations, v1.EventTypeNormal, string(reason), message)
}
