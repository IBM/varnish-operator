package controller

import (
	"context"
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"time"

	rbac "k8s.io/api/rbac/v1"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("the clusterrole", func() {
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
		},
	}

	crName := types.NamespacedName{Name: vc.Name + "-varnish-clusterrole-" + vc.Namespace, Namespace: vcNamespace}

	AfterEach(func() {
		CleanUpCreatedResources(vcName, vcNamespace)
	})

	Context("when varnishcluster is created", func() {
		It("should be created", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			Eventually(requestsChan, time.Second*5).Should(Receive(Equal(reconcile.Request{NamespacedName: types.NamespacedName{Name: vcName, Namespace: vcNamespace}})))

			cr := &rbac.ClusterRole{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), crName, cr)
			}).Should(Succeed())

			Expect(cr.Labels).To(Equal(map[string]string{
				icmv1alpha1.LabelVarnishOwner:     vcName,
				icmv1alpha1.LabelVarnishComponent: "clusterrole",
				icmv1alpha1.LabelVarnishUID:       string(newVC.UID),
			}))
			Expect(cr.Annotations).To(Equal(map[string]string{
				"varnish-cluster-name":      vcName,
				"varnish-cluster-namespace": vcNamespace,
			}))
		})
	})

	Context("when updated manually", func() {
		It("should be restored to desired state", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			Eventually(requestsChan, time.Second*5).Should(Receive(Equal(reconcile.Request{NamespacedName: types.NamespacedName{Name: vcName, Namespace: vcNamespace}})))

			cr := &rbac.ClusterRole{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), crName, cr)
			}).Should(Succeed())

			cr.Annotations = map[string]string{"rewritten": "annotations"}
			err = k8sClient.Update(context.Background(), cr)
			Expect(err).To(Succeed())

			Eventually(requestsChan, time.Second*5).Should(Receive(Equal(reconcile.Request{NamespacedName: types.NamespacedName{Name: vcName, Namespace: vcNamespace}})))

			//Expect the operator to notice the changes and restore the desired state
			Eventually(func() map[string]string {
				updatedCr := &rbac.ClusterRole{}
				err = k8sClient.Get(context.Background(), types.NamespacedName{Name: vc.Name + "-varnish-clusterrole-" + vc.Namespace, Namespace: vcNamespace}, updatedCr)
				Expect(err).To(Succeed())
				return updatedCr.Annotations
			}, time.Second*5).Should(Equal(map[string]string{
				"varnish-cluster-name":      vcName,
				"varnish-cluster-namespace": vcNamespace,
			}))
		})
	})
})