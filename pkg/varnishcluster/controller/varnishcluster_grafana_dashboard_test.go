package controller

import (
	"context"
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/names"
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

	vc := &icmv1alpha1.VarnishCluster{
		ObjectMeta: objMeta,
		Spec: icmv1alpha1.VarnishClusterSpec{
			Backend: &icmv1alpha1.VarnishClusterBackend{
				Selector: map[string]string{"app": "nginx"},
				Port:     &validBackendPort,
			},
			Service: &icmv1alpha1.VarnishClusterService{
				Port: proto.Int32(8081),
			},
			VCL: &icmv1alpha1.VarnishClusterVCL{
				ConfigMapName:      proto.String("test"),
				EntrypointFileName: proto.String("test.vcl"),
			},
			Monitoring: &icmv1alpha1.VarnishClusterMonitoring{
				GrafanaDashboard: &icmv1alpha1.VarnishClusterMonitoringGrafanaDashboard{
					Enabled:        true,
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
		It("should be created", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			dashboardCM := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), dashboardName, dashboardCM)
			}, time.Second*5).Should(Succeed())

			By("Labels should consist of standard component labels and additional user specified labels")
			Expect(dashboardCM.Labels).To(Equal(map[string]string{
				icmv1alpha1.LabelVarnishOwner:     vcName,
				icmv1alpha1.LabelVarnishComponent: icmv1alpha1.VarnishComponentGrafanaDashboard,
				icmv1alpha1.LabelVarnishUID:       string(newVC.UID),
				"foo":                             "bar",
			}))

			By("Owner reference should be set if the dashboard installed in the same namespace as VarnishCluster")
			ownerReference := []metav1.OwnerReference{
				{
					APIVersion:         "icm.ibm.com/v1alpha1",
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
	})
})
