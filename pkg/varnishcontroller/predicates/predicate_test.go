package predicates

import (
	"testing"

	"github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	"k8s.io/apimachinery/pkg/types"

	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestVarnishControllerPredicate_Update(t *testing.T) {
	clusterUID := types.UID("varnishcluster-uid")
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
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
			name: "An address in endpoints moved to not ready state should not trigger a reconcile",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
			name: "A removed port from endpoints should trigger a processing",
			event: event.UpdateEvent{
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name: "epName",
						Namespace: "epNamespace",
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
						Name: "epName",
						Namespace: "epNamespace",
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
			name: "Endpoint not matching the selector should not be processed",
			event: event.UpdateEvent{
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name: "epName",
						Namespace: "epNamespace",
						Labels: map[string]string{
							"app": "not-nginx",
						},
					},
				},
				ObjectOld: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name: "epName",
						Namespace: "epNamespace",
						Labels: map[string]string{
							"app": "not-nginx",
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "An event from correct VarnishCluster with changed backend labels should be processed",
			event: event.UpdateEvent{
				ObjectNew: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       clusterUID,
						Labels: map[string]string{
							"app": "varnish",
						},
					},
					Spec: v1alpha1.VarnishClusterSpec{
						Backend: &v1alpha1.VarnishClusterBackend{
							Selector: map[string]string{
								"app": "nginx",
							},
						},
					},
				},
				ObjectOld: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       clusterUID,
						Labels: map[string]string{
							"app": "varnish",
						},
					},
					Spec: v1alpha1.VarnishClusterSpec{
						Backend: &v1alpha1.VarnishClusterBackend{
							Selector: map[string]string{
								"app": "not-nginx",
							},
						},
					},
				},
			},
			shouldBeProcessed: true,
		},
		{
			name: "An event from VarnishCluster with changed ConfigMap version should be processed",
			event: event.UpdateEvent{
				ObjectNew: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       clusterUID,
						Labels: map[string]string{
							"app": "varnish",
						},
					},
					Spec: v1alpha1.VarnishClusterSpec{
						Backend: &v1alpha1.VarnishClusterBackend{
							Selector: map[string]string{
								"app": "nginx",
							},
						},
					},
					Status: v1alpha1.VarnishClusterStatus{
						VCL: v1alpha1.VCLStatus{
							ConfigMapVersion: "version1",
						},
					},
				},
				ObjectOld: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       clusterUID,
						Labels: map[string]string{
							"app": "varnish",
						},
					},
					Spec: v1alpha1.VarnishClusterSpec{
						Backend: &v1alpha1.VarnishClusterBackend{
							Selector: map[string]string{
								"app": "nginx",
							},
						},
					},
					Status: v1alpha1.VarnishClusterStatus{
						VCL: v1alpha1.VCLStatus{
							ConfigMapVersion: "version2",
						},
					},
				},
			},
			shouldBeProcessed: true,
		},
		{
			name: "An unchanged VarnishCluster event should not be processed",
			event: event.UpdateEvent{
				ObjectNew: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       clusterUID,
						Labels: map[string]string{
							"app": "varnish",
						},
					},
					Spec: v1alpha1.VarnishClusterSpec{
						Backend: &v1alpha1.VarnishClusterBackend{
							Selector: map[string]string{
								"app": "nginx",
							},
						},
					},
					Status: v1alpha1.VarnishClusterStatus{
						VCL: v1alpha1.VCLStatus{
							ConfigMapVersion: "version1",
						},
					},
				},
				ObjectOld: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       clusterUID,
						Labels: map[string]string{
							"app": "varnish",
						},
					},
					Spec: v1alpha1.VarnishClusterSpec{
						Backend: &v1alpha1.VarnishClusterBackend{
							Selector: map[string]string{
								"app": "nginx",
							},
						},
					},
					Status: v1alpha1.VarnishClusterStatus{
						VCL: v1alpha1.VCLStatus{
							ConfigMapVersion: "version1",
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "An event from not corresponding VarnishCluster should not be processed",
			event: event.UpdateEvent{
				ObjectNew: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       "differentUID",
						Labels: map[string]string{
							"app": "not-nginx",
						},
					},
				},
				ObjectOld: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       "differentUID",
						Labels: map[string]string{
							"app": "not-nginx",
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "An event with wrong object type for event.ObjectOld should not be processed",
			event: event.UpdateEvent{
				ObjectOld: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{
						Name: "epName",
						Namespace: "epNamespace",
					},
				},
				ObjectNew: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name: "epName",
						Namespace: "epNamespace",
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
		epPredicate := NewVarnishControllerPredicate(clusterUID, epSelectors, nil)
		if allowToProcess := epPredicate.Update(c.event); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed. Allowed to process %t. Should've been allowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
	}
}

func TestVarnishControllerPredicateCreateGeneric(t *testing.T) {
	clusterUID := types.UID("varnishcluster-uid")
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
		event             event.CreateEvent
		shouldBeProcessed bool
	}{
		{
			name: "The VarnishCluster with correct UID should be processed",
			event: event.CreateEvent{
				Object: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       clusterUID,
						Labels: map[string]string{
							"app": "nginx",
						},
					},
				},
			},
			shouldBeProcessed: true,
		},
		{
			name: "The VarnishCluster with wrong UID should not be processed",
			event: event.CreateEvent{
				Object: &v1alpha1.VarnishCluster{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "vcName",
						Namespace: "vcNamespace",
						UID:       "wrongUID",
						Labels: map[string]string{
							"app": "nginx",
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "Endpoints with not matching selector should not be processed",
			event: event.CreateEvent{
				Object: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "epName",
						Namespace: "epNamespace",
						Labels: map[string]string{
							"wrong": "labels",
						},
					},
				},
			},
			shouldBeProcessed: false,
		},
		{
			name: "Endpoints with matching selector should be processed",
			event: event.CreateEvent{
				Object: &v1.Endpoints{
					TypeMeta: v12.TypeMeta{},
					ObjectMeta: v12.ObjectMeta{
						Name:      "epName",
						Namespace: "epNamespace",
						Labels: map[string]string{
							"app": "nginx",
						},
					},
				},
			},
			shouldBeProcessed: true,
		},
	}

	for _, c := range cases {
		vcPredicate := NewVarnishControllerPredicate(clusterUID, epSelectors, logger.NewNopLogger())
		if allowToProcess := vcPredicate.Create(c.event); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed for Create. Allowed to process %t. Should've been allowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
		// the logic for Generic is the same as for Create so reuse the test cases
		if allowToProcess := vcPredicate.Generic(event.GenericEvent{Object: c.event.Object}); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed for Generic. Allowed to process %t. Should've been allowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
	}
}

func TestVarnishControllerPredicate_Delete(t *testing.T) {
	cases := []struct {
		name              string
		event             event.DeleteEvent
		shouldBeProcessed bool
	}{
		{
			name: "Deleted endpoints shouldn't trigger and event",
			event: event.DeleteEvent{
				Object: &v1.Endpoints{},
			},
			shouldBeProcessed: false,
		},
		{
			name: "Deleted VarnishClusters shouldn't trigger and event",
			event: event.DeleteEvent{
				Object: &v1alpha1.VarnishCluster{},
			},
			shouldBeProcessed: false,
		},
		{
			name: "Other watched resources deletion should trigger an event",
			event: event.DeleteEvent{
				Object: &v1.Pod{
					ObjectMeta: v12.ObjectMeta{
						Name:      "testPod",
						Namespace: "testNamespace",
					},
				}, //if a watch for a Pod will be added, it should not be ignored
			},
			shouldBeProcessed: true,
		},
	}

	for _, c := range cases {
		epPredicate := NewVarnishControllerPredicate("someUID", []labels.Selector{}, logger.NewNopLogger())

		if allowToProcess := epPredicate.Delete(c.event); allowToProcess != c.shouldBeProcessed {
			t.Fatalf("Test %q failed for Delete. Allowed to process %t. Should've been allowed to process: %t", c.name, allowToProcess, c.shouldBeProcessed)
		}
	}
}
