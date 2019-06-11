package events

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

const (
	EventRecorderName = "varnish"

	EventReasonReloadError         EventReason = "ReloadError"
	EventReasonVCLCompilationError EventReason = "VCLCompilationError"

	annotationSourcePod string = "sourcePod"
)

// EventReason is the reason why the event was create. The value appears in the 'Reason' tab of the events list
type EventReason string

// NewEventHandler creates a new event handler that will use the specified recorder
func NewEventHandler(recorder record.EventRecorder, podName string) *EventHandler {
	return &EventHandler{
		Recorder: recorder,
		podName:  podName,
	}
}

// EventHandler handles the operations for events
type EventHandler struct {
	Recorder record.EventRecorder
	podName  string
}

// Warning creates a 'warning' type event
func (e *EventHandler) Warning(object runtime.Object, reason EventReason, message string) {
	e.Recorder.AnnotatedEventf(object, map[string]string{annotationSourcePod: e.podName}, v1.EventTypeWarning, string(reason), message)
}

// Normal creates a 'normal' type event
func (e *EventHandler) Normal(object runtime.Object, reason EventReason, message string) {
	e.Recorder.AnnotatedEventf(object, map[string]string{annotationSourcePod: e.podName}, v1.EventTypeNormal, string(reason), message)
}
