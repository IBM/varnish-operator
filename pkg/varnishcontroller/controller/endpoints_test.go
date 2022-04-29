package controller

import (
	"context"
	"testing"

	"github.com/ibm/varnish-operator/pkg/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/onsi/gomega"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/config"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetBackendsEndpoint(t *testing.T) {
	baseScheme := scheme.Scheme
	utilruntime.Must(clientgoscheme.AddToScheme(baseScheme))
	utilruntime.Must(v1alpha1.AddToScheme(baseScheme))
	backendPortNumber := intstr.FromInt(4314)

	tcs := []struct {
		name              string
		vc                *v1alpha1.VarnishCluster
		podNamespace      string
		podNode           string
		k8sObjects        []client.Object
		k8sLists          []client.ObjectList
		expectedPodNumber int32
		expectedPodInfo   []PodInfo
		expectedErr       error
	}{
		{
			name: "one backend",
			vc: &v1alpha1.VarnishCluster{
				Spec: v1alpha1.VarnishClusterSpec{
					Backend: &v1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "backend"},
						Port:     &backendPortNumber,
					},
				},
			},
			podNamespace: "ns1",
			podNode:      "node1",
			k8sObjects: []client.Object{
				createTestNode("node1", map[string]string{v1.LabelTopologyZone: "zone1"}),
			},
			k8sLists: []client.ObjectList{
				&v1.PodList{
					Items: []v1.Pod{
						createTestPod("pod1", "ns1", "10.24.12.2", "node1",
							map[string]string{"app": "backend"},
							[]v1.ContainerPort{{Name: "backend", ContainerPort: backendPortNumber.IntVal}},
						),
					},
				},
			},
			expectedPodNumber: backendPortNumber.IntVal,
			expectedPodInfo: []PodInfo{
				{IP: "10.24.12.2", NodeLabels: map[string]string{v1.LabelTopologyZone: "zone1"}, PodName: "pod1", Weight: 1},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tcs {
		t.Log(tc.name)
		clientBuilder := fake.NewClientBuilder().WithScheme(baseScheme)
		clientBuilder.WithLists(tc.k8sLists...)
		clientBuilder.WithObjects(tc.k8sObjects...)
		tClient := clientBuilder.Build()

		reconciler := &ReconcileVarnish{
			config: &config.Config{
				Namespace: tc.podNamespace,
				NodeName:  tc.podNode,
			},
			Client: tClient,
			logger: logger.NewNopLogger(),
		}
		podInfo, portNumber, _, _, err := reconciler.getBackendEndpoints(context.Background(), tc.vc)

		a := gomega.NewGomegaWithT(t)
		if tc.expectedErr == nil {
			a.Expect(err).To(gomega.BeNil())
		} else {
			a.Expect(err.Error()).To(gomega.Equal(tc.expectedErr.Error()))
		}
		a.Expect(podInfo).To(gomega.Equal(tc.expectedPodInfo))
		a.Expect(portNumber).To(gomega.Equal(tc.expectedPodNumber))

	}
}

func createTestPod(name, namespace, ip, nodeName string, labels map[string]string, ports []v1.ContainerPort) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.PodSpec{
			NodeName: nodeName,
			Containers: []v1.Container{
				{
					Ports: ports,
				},
			},
		},
		Status: v1.PodStatus{
			PodIP: ip,
		},
	}
}

func createTestNode(name string, labels map[string]string) *v1.Node {
	node := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Status: v1.NodeStatus{
			Addresses: []v1.NodeAddress{
				{
					Type:    v1.NodeInternalIP,
					Address: "192.24.51.2",
				},
			},
		},
	}

	if len(labels) > 0 {
		node.Labels = labels
	}

	return &node
}
