package predicates

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cin/varnish-operator/pkg/logger"
	"k8s.io/apimachinery/pkg/labels"

	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestLabelMatcherPredicate_Update(t *testing.T) {
	tcs := []struct {
		name               string
		selector           map[string]string
		updateEvent        event.UpdateEvent
		shouldTriggerEvent bool
	}{
		{
			name:     "nothing changed",
			selector: map[string]string{"app": "backend"},
			updateEvent: event.UpdateEvent{
				ObjectOld: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend"}},
					Spec:       v1.PodSpec{NodeName: "node1"},
					Status:     v1.PodStatus{PodIP: "19.43.11.32"},
				},
				ObjectNew: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend"}},
					Spec:       v1.PodSpec{NodeName: "node1"},
					Status:     v1.PodStatus{PodIP: "19.43.11.32"},
				},
			},
			shouldTriggerEvent: false,
		},
		{
			name:     "nothing significant changed",
			selector: map[string]string{"app": "backend"},
			updateEvent: event.UpdateEvent{
				ObjectOld: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend", "one": "two"}},
					Spec:       v1.PodSpec{NodeName: "node1"},
					Status:     v1.PodStatus{PodIP: "19.43.11.32"},
				},
				ObjectNew: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend", "one": "three"}},
					Spec:       v1.PodSpec{NodeName: "node1"},
					Status:     v1.PodStatus{PodIP: "19.43.11.32"},
				},
			},
			shouldTriggerEvent: false,
		},
		{
			name:     "pod doesn't match selector",
			selector: map[string]string{"app": "backend"},
			updateEvent: event.UpdateEvent{
				ObjectOld: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend4"}},
					Spec:       v1.PodSpec{NodeName: "node1"},
					Status:     v1.PodStatus{PodIP: "19.43.11.32"},
				},
				ObjectNew: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend4"}},
					Spec:       v1.PodSpec{NodeName: "node2"},
					Status:     v1.PodStatus{PodIP: "19.43.11.33"},
				},
			},
			shouldTriggerEvent: false,
		},
		{
			name:     "ip changed (pod scheduled)",
			selector: map[string]string{"app": "backend"},
			updateEvent: event.UpdateEvent{
				ObjectOld: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend"}},
					Spec:       v1.PodSpec{NodeName: "node2"},
					Status:     v1.PodStatus{PodIP: ""},
				},
				ObjectNew: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend"}},
					Spec:       v1.PodSpec{NodeName: "node2"},
					Status:     v1.PodStatus{PodIP: "19.43.11.33"},
				},
			},
			shouldTriggerEvent: true,
		},
		{
			name:     "node changed (pod scheduled)",
			selector: map[string]string{"app": "backend"},
			updateEvent: event.UpdateEvent{
				ObjectOld: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend"}},
					Spec:       v1.PodSpec{NodeName: ""},
					Status:     v1.PodStatus{PodIP: "19.43.11.33"},
				},
				ObjectNew: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{Name: "pod1", Labels: map[string]string{"app": "backend"}},
					Spec:       v1.PodSpec{NodeName: "node2"},
					Status:     v1.PodStatus{PodIP: "19.43.11.33"},
				},
			},
			shouldTriggerEvent: true,
		},
	}

	for _, tc := range tcs {
		predicate := NewLabelMatcherPredicate(labels.SelectorFromSet(tc.selector), logger.NewNopLogger())
		if predicate.Update(tc.updateEvent) != tc.shouldTriggerEvent {
			t.Logf(tc.name+": expected %t got %t", tc.shouldTriggerEvent, !tc.shouldTriggerEvent)
			t.Fail()
		}
	}
}
