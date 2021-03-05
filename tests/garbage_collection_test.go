package tests

import (
	"context"
	"fmt"
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/names"
	"time"

	apps "k8s.io/api/apps/v1"
	rbac "k8s.io/api/rbac/v1"

	v1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/gogo/protobuf/proto"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Varnish Cluster", func() {
	validBackendPort := intstr.FromInt(8080)
	vcNamespace := "default"
	vcName := "test"
	objMeta := metav1.ObjectMeta{
		Namespace: vcNamespace,
		Name:      vcName,
	}

	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &vcapi.VarnishCluster{}, client.InNamespace(vcNamespace))
		Expect(err).To(Succeed())
		waitUntilVarnishClusterRemoved(vcName, vcNamespace)
	})

	Context("when deleted", func() {
		It("should also garbage collect created resources", func() {
			minAvailabe := intstr.FromInt(1)
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Varnish: &vcapi.VarnishClusterVarnish{},
					PodDisruptionBudget: &policyv1beta1.PodDisruptionBudgetSpec{
						MinAvailable: &minAvailabe,
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
							DatasourceName: proto.String("Prometheus-datasource"),
						},
					},
				},
			}

			By("Creating VarnishCluster")
			Expect(k8sClient.Create(context.Background(), vc)).To(Succeed())

			By("Checking if all resources are created")
			expectResourceIsCreated(types.NamespacedName{Name: *vc.Spec.VCL.ConfigMapName, Namespace: vcNamespace}, &v1.ConfigMap{})
			expectResourceIsCreated(types.NamespacedName{Name: names.GrafanaDashboard(vcName), Namespace: vcNamespace}, &v1.ConfigMap{})
			expectResourceIsCreated(types.NamespacedName{Name: names.HeadlessService(vcName), Namespace: vcNamespace}, &v1.Service{})
			expectResourceIsCreated(types.NamespacedName{Name: names.NoCacheService(vcName), Namespace: vcNamespace}, &v1.Service{})
			expectResourceIsCreated(types.NamespacedName{Name: vcName, Namespace: vcNamespace}, &v1.Service{})
			expectResourceIsCreated(types.NamespacedName{Name: names.VarnishSecret(vcName), Namespace: vcNamespace}, &v1.Secret{})
			expectResourceIsCreated(types.NamespacedName{Name: names.PodDisruptionBudget(vcName), Namespace: vcNamespace}, &policyv1beta1.PodDisruptionBudget{})
			expectResourceIsCreated(types.NamespacedName{Name: names.Role(vcName), Namespace: vcNamespace}, &rbac.Role{})
			expectResourceIsCreated(types.NamespacedName{Name: names.RoleBinding(vcName), Namespace: vcNamespace}, &rbac.RoleBinding{})
			expectResourceIsCreated(types.NamespacedName{Name: names.ClusterRole(vcName, vcNamespace)}, &rbac.ClusterRole{})
			expectResourceIsCreated(types.NamespacedName{Name: names.ClusterRoleBinding(vcName, vcNamespace)}, &rbac.ClusterRoleBinding{})
			expectResourceIsCreated(types.NamespacedName{Name: names.ServiceAccount(vcName), Namespace: vcNamespace}, &v1.ServiceAccount{})
			expectResourceIsCreated(types.NamespacedName{Name: names.StatefulSet(vcName), Namespace: vcNamespace}, &apps.StatefulSet{})

			By("Deleting VarnishCluster")
			Expect(k8sClient.Delete(context.Background(), vc)).To(Succeed())

			By("Waiting for VarnishCluster to be removed")
			waitUntilVarnishClusterRemoved(vcName, vcNamespace)

			By("Checking if all created resources are deleted")
			expectResourceIsDeleted(types.NamespacedName{Name: *vc.Spec.VCL.ConfigMapName, Namespace: vcNamespace}, &v1.ConfigMap{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.GrafanaDashboard(vcName), Namespace: vcNamespace}, &v1.ConfigMap{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.HeadlessService(vcName), Namespace: vcNamespace}, &v1.Service{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.NoCacheService(vcName), Namespace: vcNamespace}, &v1.Service{})
			expectResourceIsDeleted(types.NamespacedName{Name: vcName, Namespace: vcNamespace}, &v1.Service{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.PodDisruptionBudget(vcName), Namespace: vcNamespace}, &policyv1beta1.PodDisruptionBudget{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.Role(vcName), Namespace: vcNamespace}, &rbac.Role{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.RoleBinding(vcName), Namespace: vcNamespace}, &rbac.RoleBinding{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.ClusterRole(vcName, vcNamespace)}, &rbac.ClusterRole{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.ClusterRoleBinding(vcName, vcNamespace)}, &rbac.ClusterRoleBinding{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.ServiceAccount(vcName), Namespace: vcNamespace}, &v1.ServiceAccount{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.StatefulSet(vcName), Namespace: vcNamespace}, &apps.StatefulSet{})
			expectResourceIsDeleted(types.NamespacedName{Name: names.VarnishSecret(vcName), Namespace: vcNamespace}, &v1.Secret{})
		})
	})
})

func expectResourceIsCreated(name types.NamespacedName, obj client.Object) {
	Eventually(func() error {
		return k8sClient.Get(context.Background(), name, obj)
	}, time.Second*5).Should(Succeed(), fmt.Sprintf("%T %s expected to exist", obj, name))
}

func expectResourceIsDeleted(name types.NamespacedName, obj client.Object) {
	Eventually(func() metav1.StatusReason {
		err := k8sClient.Get(context.Background(), name, obj)
		if err != nil {
			if statusErr, ok := err.(*errors.StatusError); ok {
				return statusErr.ErrStatus.Reason
			}
		}

		return "Found"
	}, time.Second*5).Should(Equal(metav1.StatusReasonNotFound), fmt.Sprintf("%T %s should be deleted", obj, name))
}
