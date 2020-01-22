package tests

import (
	"bytes"
	"context"
	"fmt"
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"net/http"
	"net/url"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	prometheusClient "github.com/prometheus/client_model/go"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getMetric returns the first metric in the specified metric family
func getMetric(metricFamilies map[string]*prometheusClient.MetricFamily, metricName string) (prometheusClient.Metric, bool) {
	mf, found := metricFamilies[metricName]
	if !found {
		return prometheusClient.Metric{}, false
	}

	if mf.Name != nil && *mf.Name == metricName {
		if len(mf.Metric) == 0 {
			return prometheusClient.Metric{}, false
		}

		return *mf.Metric[0], true
	}
	return prometheusClient.Metric{}, false
}

// getMetricByLabel returns the first metric in the specified metric family that has the specified label and label value. The label value parameter is a substring.
func getMetricByLabel(metricFamilies map[string]*prometheusClient.MetricFamily, metricName, labelName, labelValue string) (prometheusClient.Metric, bool) {
	mf, found := metricFamilies[metricName]
	if !found {
		return prometheusClient.Metric{}, false
	}

	if mf.Name != nil && *mf.Name == metricName {
		for _, metric := range mf.Metric {
			for _, label := range metric.Label {
				if *label.Name == labelName && strings.Contains(*label.Value, labelValue) {
					return *metric, true
				}
			}
		}
	}
	return prometheusClient.Metric{}, false
}

func waitForPodsTermination(namespace string, selector map[string]string) {
	Eventually(func() bool {
		podList := &v1.PodList{}
		err := k8sClient.List(context.Background(), podList, client.InNamespace(namespace), client.MatchingLabels(selector))
		Expect(err).To(Succeed())

		return len(podList.Items) == 0
	}, time.Minute, time.Second*5).Should(BeTrue())
}

func waitForPodsReadiness(namespace string, selector map[string]string) {
	Eventually(func() bool {
		podList := &v1.PodList{}
		err := k8sClient.List(context.Background(), podList, client.InNamespace(namespace), client.MatchingLabels(selector))
		Expect(err).To(Succeed())

		if len(podList.Items) == 0 {
			return false
		}

		for _, pod := range podList.Items {
			for _, container := range pod.Status.ContainerStatuses {
				if !container.Ready {
					return false
				}
			}
		}
		return true
	}, time.Minute, time.Second*2).Should(BeTrue(), "pods should become ready")
}

func waitUntilVarnishClusterRemoved(name, namespace string) {
	Eventually(func() metav1.StatusReason {
		err := k8sClient.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, &icmv1alpha1.VarnishCluster{})
		if err != nil {
			if statusErr, ok := err.(*errors.StatusError); ok {
				return statusErr.ErrStatus.Reason
			}
		}

		return "Found"
	}, time.Second*5).Should(Equal(metav1.StatusReasonNotFound), "the varnishcluster should be deleted after finalizers are finished")
}

func portForwardPod(namespace string, selector map[string]string, portsToForward []string) *portforward.PortForwarder {
	rt, upgrader, err := spdy.RoundTripperFor(restConfig)
	Expect(err).To(Succeed())

	podList := &v1.PodList{}
	err = k8sClient.List(context.Background(), podList, client.InNamespace(namespace), client.MatchingLabels(selector))
	Expect(err).To(Succeed())
	Expect(len(podList.Items)).ToNot(BeEquivalentTo(0), "no pods found to port-forward")

	pod := podList.Items[0]
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", pod.Namespace, pod.Name)
	hostIP := strings.TrimLeft(restConfig.Host, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: rt}, http.MethodPost, &serverURL)

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	pf, err := portforward.New(dialer, portsToForward, stopChan, readyChan, out, errOut)
	Expect(err).To(Succeed())

	go func() {
		defer GinkgoRecover()
		for range readyChan {
		}

		if len(errOut.String()) != 0 {
			Fail(fmt.Sprintf("Port forwarding failed: %s", errOut.String()))
		} else if len(out.String()) != 0 {
			_, err := fmt.Fprintf(GinkgoWriter, "Message from port forwarder: %s", out.String())
			Expect(err).To(Succeed())
		}
	}()

	go func() {
		defer GinkgoRecover()
		err := pf.ForwardPorts() //this will block until stopped
		Expect(err).To(Succeed())
	}()

	for range readyChan {
	} //wait till ready
	return pf
}
