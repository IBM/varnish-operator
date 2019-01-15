package controller

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

const (
	EventRecorderNameVarnishService             = "varnish-service" // appears in the 'From' column of the events list
	EventReasonDeploymentError      EventReason = "DeploymentError"
)

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
	e.Recorder.Event(object, v1.EventTypeWarning, string(reason), message)
}

// Normal creates a 'normal' type event
func (e *EventHandler) Normal(object runtime.Object, reason EventReason, message string) {
	e.Recorder.Event(object, v1.EventTypeNormal, string(reason), message)
}
