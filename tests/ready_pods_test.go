package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"

	"github.com/prometheus/common/expfmt"
	appsv1 "k8s.io/api/apps/v1"

	v1 "k8s.io/api/core/v1"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Varnish cluster", func() {
	vcNamespace := "default"
	vcName := "test"
	objMeta := metav1.ObjectMeta{
		Namespace: vcNamespace,
		Name:      vcName,
	}
	backendResponse := "TEST"
	backendLabels := map[string]string{"app": "test-backend"}
	backendDeploymentName := "test-backend"
	varnishPodLabels := map[string]string{vcapi.LabelVarnishComponent: vcapi.VarnishComponentVarnish}

	backendsDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendDeploymentName,
			Namespace: vcNamespace,
			Labels:    backendLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: proto.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: backendLabels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: backendLabels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "backend",
							Image: "hashicorp/http-echo",
							Ports: []v1.ContainerPort{
								{
									Name:          "web",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 5678,
								},
							},
							Args: []string{fmt.Sprintf("-text=%s", backendResponse)},
						},
					},
				},
			},
		},
	}

	backendPort := intstr.FromInt(5678)
	vc := &vcapi.VarnishCluster{
		ObjectMeta: objMeta,
		Spec: vcapi.VarnishClusterSpec{
			Backend: &vcapi.VarnishClusterBackend{
				Selector: backendLabels,
				Port:     &backendPort,
			},
			Service: &vcapi.VarnishClusterService{
				Port: proto.Int32(9090),
			},
			Varnish: &vcapi.VarnishClusterVarnish{
				ImagePullPolicy: v1.PullNever,
				Controller: &vcapi.VarnishClusterVarnishController{
					ImagePullPolicy: v1.PullNever,
				},
				MetricsExporter: &vcapi.VarnishClusterVarnishMetricsExporter{
					ImagePullPolicy: v1.PullNever,
				},
			},
			VCL: &vcapi.VarnishClusterVCL{
				ConfigMapName:      proto.String("test"),
				EntrypointFileName: proto.String("test.vcl"),
			},
		},
	}

	AfterEach(func() {
		By("deleting created resources")
		Expect(k8sClient.DeleteAllOf(context.Background(), &vcapi.VarnishCluster{}, client.InNamespace(vcNamespace))).To(Succeed())
		Expect(k8sClient.DeleteAllOf(context.Background(), &appsv1.Deployment{}, client.InNamespace(vcNamespace), client.MatchingLabels(backendLabels))).To(Succeed())
		waitForPodsTermination(vcNamespace, varnishPodLabels)
		waitForPodsTermination(vcNamespace, backendLabels)
		waitUntilVarnishClusterRemoved(vcName, vcNamespace)
	})

	It("pods respond with backend responses and metrics", func() {
		Expect(k8sClient.Create(context.Background(), backendsDeployment)).To(Succeed())
		Expect(k8sClient.Create(context.Background(), vc)).To(Succeed())
		By("backend pods become ready")
		waitForPodsReadiness(vcNamespace, backendLabels)
		By("varnish pods become ready")
		waitForPodsReadiness(vcNamespace, varnishPodLabels)
		pf := portForwardPod(vcNamespace, varnishPodLabels, []string{"6081:6081", "9131:9131"})
		defer pf.Close()

		By("varnish pod responds with the backend response")
		var resp *http.Response
		Eventually(func() (int, error) {
			var err error
			resp, err = http.Get("http://localhost:6081/test")
			if err != nil {
				return 0, err
			}
			return resp.StatusCode, nil
		}, time.Second*30, time.Second*2).Should(Equal(200))
		Expect(resp.Header.Get("X-Varnish-Cache")).To(Equal("MISS"))
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.TrimSpace(string(body))).To(Equal("TEST"))

		By("varnish pod responds with the cached response")
		resp, err = http.Get("http://localhost:6081/test")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))
		Expect(resp.Header.Get("X-Varnish-Cache")).To(Equal("HIT"))
		body, err = ioutil.ReadAll(resp.Body)
		Expect(err).To(Succeed())
		Expect(strings.TrimSpace(string(body))).To(Equal("TEST"))

		By("varnish pod respond with prometheus metrics")
		resp, err = http.Get("http://localhost:9131/metrics")
		Expect(err).To(Succeed())
		Expect(resp.StatusCode).To(Equal(200))

		var parser expfmt.TextParser
		metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		metric, found := getMetricByLabel(metricFamilies, "varnish_backend_req", "backend", backendDeploymentName)
		Expect(found).To(BeTrue())
		Expect(*metric.Counter.Value).To(BeNumerically(">=", 1))
		metric, found = getMetricByLabel(metricFamilies, "varnish_backend_bereq_hdrbytes", "backend", backendDeploymentName)
		Expect(found).To(BeTrue())
		Expect(*metric.Counter.Value).To(BeNumerically(">=", 1))
		metric, found = getMetric(metricFamilies, "varnish_main_n_object")
		Expect(found).To(BeTrue())
		Expect(*metric.Gauge.Value).To(BeNumerically(">=", 1))
		metric, found = getMetric(metricFamilies, "varnish_main_n_vcl")
		Expect(found).To(BeTrue())
		Expect(*metric.Gauge.Value).To(BeNumerically(">=", 2)) //should be the `default` and the one we loaded
		metric, found = getMetric(metricFamilies, "varnish_main_uptime")
		Expect(found).To(BeTrue())
		Expect(*metric.Counter.Value).To(BeNumerically(">=", 1))
	})
})
