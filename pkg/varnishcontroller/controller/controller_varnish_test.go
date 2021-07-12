package controller

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ibm/varnish-operator/api/v1alpha1"
	varnishEvents "github.com/ibm/varnish-operator/pkg/varnishcontroller/events"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/metrics"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/varnishadm"

	"github.com/onsi/gomega"
	prometheusClient "github.com/prometheus/client_model/go"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type varnishMock struct {
	reloadResponse       string
	pingError            error
	listResponse         []varnishadm.VCLConfig
	listError            error
	discardError         error
	reloadError          error
	activeVCLConfigName  string
	activeVCLConfigError error
}

func (v *varnishMock) Ping() error {
	return v.pingError
}

func (v *varnishMock) List() ([]varnishadm.VCLConfig, error) {
	return v.listResponse, v.listError
}

func (v *varnishMock) Reload(version, entry string) ([]byte, error) {
	return []byte(v.reloadResponse), v.reloadError
}

func (v *varnishMock) GetActiveConfigurationName() (string, error) {
	return v.activeVCLConfigName, v.activeVCLConfigError
}

func (v *varnishMock) Discard(vclConfigName string) error {
	return v.discardError
}

type eventsObserver struct {
	eventsObserved bool
}

func (e *eventsObserver) Event(object runtime.Object, eventtype, reason, message string) {
	e.eventsObserved = true
}

func (e *eventsObserver) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	e.eventsObserved = true
}

func (e *eventsObserver) PastEventf(object runtime.Object, timestamp metav1.Time, eventtype, reason, messageFmt string, args ...interface{}) {
	e.eventsObserved = true
}

func (e *eventsObserver) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	e.eventsObserved = true
}

func TestReconcileVarnish(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	entrypointFileName := "test_entrypoint_file_name.vcl"
	cmResourceVersion := "test_resource_version"
	varnishAdmError := errors.New("error from varnish")
	cases := []struct {
		varnish                           *varnishMock
		pod                               *v1.Pod
		varnishcluster                    *v1alpha1.VarnishCluster
		configMap                         *v1.ConfigMap
		expectEventSent                   bool
		expectedError                     error
		expectedVCLCompilationErrorMetric int
	}{
		{
			varnish: &varnishMock{
				reloadError:    varnishAdmError,
				reloadResponse: `Something something VCL compilation failed`,
			},
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion: cmResourceVersion,
				},
			},
			varnishcluster: &v1alpha1.VarnishCluster{
				Spec: v1alpha1.VarnishClusterSpec{
					VCL: &v1alpha1.VarnishClusterVCL{
						EntrypointFileName: &entrypointFileName,
					},
					HaproxySidecar: &v1alpha1.HaproxySidecar{
						Enabled: true,
					},
				},
			},
			pod:                               &v1.Pod{},
			expectedError:                     nil,
			expectEventSent:                   true,
			expectedVCLCompilationErrorMetric: 1,
		},
		{
			varnish: &varnishMock{
				reloadError:    nil,
				reloadResponse: "success",
			},
			configMap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion: cmResourceVersion,
				},
			},
			varnishcluster: &v1alpha1.VarnishCluster{
				Spec: v1alpha1.VarnishClusterSpec{
					VCL: &v1alpha1.VarnishClusterVCL{
						EntrypointFileName: &entrypointFileName,
					},
				},
			},
			pod:                               &v1.Pod{},
			expectedError:                     nil,
			expectEventSent:                   false,
			expectedVCLCompilationErrorMetric: 0,
		},
	}

	for _, c := range cases {
		events := &eventsObserver{}
		controllerMetrics := metrics.NewVarnishControllerMetrics()
		testReconciler := &ReconcileVarnish{
			varnish:      c.varnish,
			eventHandler: &varnishEvents.EventHandler{Recorder: events},
			metrics:      controllerMetrics,
		}
		err := testReconciler.reconcileVarnish(context.Background(), c.varnishcluster, c.pod, c.configMap)
		g.Expect(reflect.DeepEqual(err, c.expectedError)).To(gomega.BeTrue())
		g.Expect(events.eventsObserved).To(gomega.Equal(c.expectEventSent))
		m := &prometheusClient.Metric{}
		err = controllerMetrics.VCLCompilationError.Write(m)
		g.Expect(err).To(gomega.Succeed())
		g.Expect(*m.Gauge.Value).To(gomega.BeNumerically("==", c.expectedVCLCompilationErrorMetric))
	}
}
