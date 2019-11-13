package predicates

import (
	"icm-varnish-k8s-operator/pkg/logger"
	"testing"

	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestEndpointsUpdatePredicate(t *testing.T) {
	epSelectors := []labels.Selector{
		labels.Set{
			"app": "nginx",
		}.AsSelector(),
		labels.Set{
			"varnish-component": "varnish-service",
			"varnish-owner":     "varnishcluster-example",
			"varnish-uid":       "some-uid",
		}.AsSelector(),
	}

	cases := []struct {
		name              string
		event             event.UpdateEvent
		shouldBeProcessed bool
	}{
		{
			name: "Nothing has changed nginx endpoints",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "Nothing has changed for varnish endpoints",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"varnish-component": "varnish-service",
							"varnish-owner":     "varnishcluster-example",
							"varnish-uid":       "some-uid",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.4",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.5",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"varnish-component": "varnish-service",
							"varnish-owner":     "varnishcluster-example",
							"varnish-uid":       "some-uid",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.4",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.5",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "Only the order has changed",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "An endpoint has been removed",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			shouldBeProcessed: true,
		},
		{
			name: "An address moved to not ready state should not trigger a reconcile",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
							},
							NotReadyAddresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "A removed port should trigger a processing",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
								{
									Name:     "metrics",
									Port:     8080,
									Protocol: "TCP",
								},
							},
						},
					},
				},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Subsets: []v1.EndpointSubset{
						{
							Addresses: []v1.EndpointAddress{
								{
									IP:       "192.168.0.1",
									Hostname: "host1",
									NodeName: func(str string) *string { return &str }("node1"),
								},
								{
									IP:       "192.168.0.2",
									Hostname: "host2",
									NodeName: func(str string) *string { return &str }("node2"),
								},
							},
							NotReadyAddresses: nil,
							Ports: []v1.EndpointPort{
								{
									Name:     "web",
									Port:     80,
									Protocol: "TCP",
								},
							},
						},
					},
				},
			},
			shouldBeProcessed: true,
		},
		{
			name: "Wrong object type in event.ObjectNew should filter out the request",
			event: event.UpdateEvent{
				ObjectNew: &v1.Pod{},
			},
			shouldBeProcessed: true,
		},
		{
			name: "Endpoint not matching the selector should not be processed",
			event: event.UpdateEvent{
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "not-nginx",
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "Event with wrong object type for event.ObjectOld should not be processed",
			event: event.UpdateEvent{
				ObjectOld: &v1.Pod{},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
	}

	for _, c := range cases {
		epPredicate := NewEndpointsSelectors(epSelectors, logger.NewNopLogger())
		if allowToProcess := epPredicate.Update(c.event); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed. Allowed to process %t, Should've been alllowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
	}
}

func TestEndpointsSharedPredicate(t *testing.T) {
	epSelectors := []labels.Selector{
		labels.Set{
			"app": "nginx",
		}.AsSelector(),
		labels.Set{
			"varnish-component": "varnish-service",
			"varnish-owner":     "varnishcluster-example",
			"varnish-uid":       "some-uid",
		}.AsSelector(),
	}

	cases := []struct {
		name              string
		obj               *v1.Endpoints
		shouldBeProcessed bool
	}{
		{
			name: "Matches the backend selector",
			obj: &v1.Endpoints{
				TypeMeta: v12.TypeMeta{},
				ObjectMeta: v12.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
			},
			shouldBeProcessed: true,
		},
		{
			name: "Matches the varnish selector",
			obj: &v1.Endpoints{
				TypeMeta: v12.TypeMeta{},
				ObjectMeta: v12.ObjectMeta{
					Labels: map[string]string{
						"varnish-component": "varnish-service",
						"varnish-owner":     "varnishcluster-example",
						"varnish-uid":       "some-uid",
					},
				},
			},
			shouldBeProcessed: true,
		},
		{
			name: "Not matches any selector",
			obj: &v1.Endpoints{
				TypeMeta: v12.TypeMeta{},
				ObjectMeta: v12.ObjectMeta{
					Labels: map[string]string{
						"other": "app",
					},
				},
			},
			shouldBeProcessed: false,
		},
	}

	for _, c := range cases {
		epPredicate := NewEndpointsSelectors(epSelectors, logger.NewNopLogger())
		if allowToProcess := epPredicate.Create(event.CreateEvent{Object: c.obj, Meta: c.obj}); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed for Create. Allowed to process %t, Should've been alllowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
		if allowToProcess := epPredicate.Delete(event.DeleteEvent{Object: c.obj, Meta: c.obj}); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed for Delete. Allowed to process %t, Should've been alllowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
		if allowToProcess := epPredicate.Generic(event.GenericEvent{Object: c.obj, Meta: c.obj}); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed for Generic. Allowed to process %t, Should've been alllowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
	}

	epPredicate := NewEndpointsSelectors(epSelectors, logger.NewNopLogger())
	allowToProcess := epPredicate.Generic(event.GenericEvent{
		Meta:   nil,
		Object: &v1.Pod{}, //not an enpdpoint
	})

	if !allowToProcess {
		t.Fatalf("Non Enpoints events should not be filtered out")
	}
}
