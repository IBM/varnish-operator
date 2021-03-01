package controller

import (
	"context"
	"fmt"
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/names"
	"k8s.io/apimachinery/pkg/util/json"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("grafana dashboard", func() {
	validBackendPort := intstr.FromInt(8080)
	vcNamespace := "default"
	vcName := "test"
	objMeta := metav1.ObjectMeta{
		Namespace: vcNamespace,
		Name:      vcName,
	}

	vc := &vcapi.VarnishCluster{
		ObjectMeta: objMeta,
		Spec: vcapi.VarnishClusterSpec{
			Backend: &vcapi.VarnishClusterBackend{
				Selector: map[string]string{"app": "nginx"},
				Port:     &validBackendPort,
			},
			Service: &vcapi.VarnishClusterService{
				Port: proto.Int32(8081),
			},
			VCL: &vcapi.VarnishClusterVCL{
				ConfigMapName:      proto.String("test"),
				EntrypointFileName: proto.String("test.vcl"),
			},
			Monitoring: &vcapi.VarnishClusterMonitoring{
				GrafanaDashboard: &vcapi.VarnishClusterMonitoringGrafanaDashboard{
					Enabled:        true,
					Title:          "",
					Namespace:      "",
					Labels:         map[string]string{"foo": "bar"},
					DatasourceName: proto.String("Prometheus"),
				},
			},
		},
	}

	dashboardName := types.NamespacedName{Name: names.GrafanaDashboard(vc.Name), Namespace: vcNamespace}

	AfterEach(func() {
		CleanUpCreatedResources(vcName, vcNamespace)
	})

	Context("when varnishcluster is created and dashboard is enabled", func() {
		It("should be created with default config", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			dashboardCM := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), dashboardName, dashboardCM)
			}, time.Second*5).Should(Succeed())

			By("Labels should consist of standard component labels and additional user specified labels")
			Expect(dashboardCM.Labels).To(Equal(map[string]string{
				vcapi.LabelVarnishOwner:     vcName,
				vcapi.LabelVarnishComponent: vcapi.VarnishComponentGrafanaDashboard,
				vcapi.LabelVarnishUID:       string(newVC.UID),
				"foo":                       "bar",
			}))

			By("Dashboard title should default to <cluster name> varnish")
			dashboardString := dashboardCM.Data[names.GrafanaDashboardFile(newVC.Name)]
			var data map[string]interface{}
			_ = json.Unmarshal([]byte(dashboardString), &data)
			Expect(data["title"].(string)).To(Equal(fmt.Sprintf("Varnish (%s)", newVC.Name)))

			By("Owner reference should be set if the dashboard installed in the same namespace as VarnishCluster")
			ownerReference := []metav1.OwnerReference{
				{
					APIVersion:         "caching.ibm.com/v1alpha1",
					Kind:               "VarnishCluster",
					Name:               newVC.Name,
					UID:                newVC.UID,
					Controller:         proto.Bool(true),
					BlockOwnerDeletion: proto.Bool(true),
				},
			}
			Expect(dashboardCM.OwnerReferences).To(Equal(ownerReference))

			By("Create the namespace to install the dashboard to")
			overrideNamespace := &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "different-namespace",
				},
			}
			Expect(k8sClient.Create(context.Background(), overrideNamespace)).To(Succeed())

			By("Overriding the namespace to install the dashboard to")
			Eventually(func() error {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: vcName, Namespace: vcNamespace}, newVC)
				if err != nil {
					return err
				}
				newVC.Spec.Monitoring.GrafanaDashboard.Namespace = overrideNamespace.Name
				err = k8sClient.Update(context.Background(), newVC)
				return err
			}, time.Second*10).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: names.GrafanaDashboard(vcName), Namespace: overrideNamespace.Name}, &v1.ConfigMap{})
			}, time.Second*10).Should(Succeed(), "The dashboard should be installed in the specified namespace")

			Eventually(func() metav1.StatusReason {
				err := k8sClient.Get(context.Background(), dashboardName, &v1.ConfigMap{})
				if err != nil {
					if statusErr, ok := err.(*errors.StatusError); ok {
						return statusErr.ErrStatus.Reason
					}
				}

				return "Found"
			}, time.Second*10).Should(Equal(metav1.StatusReasonNotFound), "The dashboard in default namespace should be deleted")

			By("Disabling the dashboard")
			Eventually(func() error {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: vcName, Namespace: vcNamespace}, newVC)
				if err != nil {
					return err
				}
				newVC.Spec.Monitoring.GrafanaDashboard.Enabled = false
				err = k8sClient.Update(context.Background(), newVC)
				return err
			}, time.Second*10).Should(Succeed())

			Eventually(func() metav1.StatusReason {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: names.GrafanaDashboard(vcName), Namespace: overrideNamespace.Name}, &v1.ConfigMap{})
				if err != nil {
					if statusErr, ok := err.(*errors.StatusError); ok {
						return statusErr.ErrStatus.Reason
					}
				}

				return "Found"
			}, time.Second*10).Should(Equal(metav1.StatusReasonNotFound), "The dashboard should be deleted")
		})

		It("should be created if dashboard title override specified", func() {
			newVC := vc.DeepCopy()
			newVC.Spec.Monitoring.GrafanaDashboard.Title = "Test Varnish Dashboard"
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			dashboardCM := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), dashboardName, dashboardCM)
			}, time.Second*5).Should(Succeed())

			By("Dashboard title should be overridden")
			dashboardString := dashboardCM.Data[names.GrafanaDashboardFile(newVC.Name)]
			var data map[string]interface{}
			_ = json.Unmarshal([]byte(dashboardString), &data)
			Expect(data["title"].(string)).To(Equal(newVC.Spec.Monitoring.GrafanaDashboard.Title))
		})
	})
})
