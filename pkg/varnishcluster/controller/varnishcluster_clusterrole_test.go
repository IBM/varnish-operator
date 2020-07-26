package controller

import (
	"context"
	"time"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/names"

	rbac "k8s.io/api/rbac/v1"

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
		},
	}

	crName := types.NamespacedName{Name: names.ClusterRole(vc.Name, vc.Namespace)}

	AfterEach(func() {
		CleanUpCreatedResources(vcName, vcNamespace)
	})

	Context("when varnishcluster is created", func() {
		It("should be created", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			cr := &rbac.ClusterRole{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), crName, cr)
			}, time.Second*5).Should(Succeed())

			Expect(cr.Labels).To(Equal(map[string]string{
				vcapi.LabelVarnishOwner:     vcName,
				vcapi.LabelVarnishComponent: vcapi.VarnishComponentClusterRole,
				vcapi.LabelVarnishUID:       string(newVC.UID),
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

			cr := &rbac.ClusterRole{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), crName, cr)
			}, time.Second*5).Should(Succeed())

			cr.Annotations = map[string]string{"rewritten": "annotations"}
			err = k8sClient.Update(context.Background(), cr)
			Expect(err).To(Succeed())

			//Expect the operator to notice the changes and restore the desired state
			Eventually(func() map[string]string {
				updatedCr := &rbac.ClusterRole{}
				err = k8sClient.Get(context.Background(), types.NamespacedName{Name: names.ClusterRole(vc.Name, vc.Namespace)}, updatedCr)
				Expect(err).To(Succeed())
				return updatedCr.Annotations
			}, time.Second*5).Should(Equal(map[string]string{
				"varnish-cluster-name":      vcName,
				"varnish-cluster-namespace": vcNamespace,
			}))
		})
	})
})
