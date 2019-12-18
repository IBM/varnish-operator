package tests

import (
	"context"
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Validating webhook", func() {
	validBackendPort := intstr.FromInt(8080)
	vcNamespace := "default"
	vcName := "test"
	objMeta := metav1.ObjectMeta{
		Namespace: vcNamespace,
		Name:      vcName,
	}

	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &icmv1alpha1.VarnishCluster{}, client.InNamespace(vcNamespace))
		Expect(err).To(Succeed())
	})

	It("is working", func() {
		vc := &icmv1alpha1.VarnishCluster{
			ObjectMeta: objMeta,
			Spec: icmv1alpha1.VarnishClusterSpec{
				Backend: &icmv1alpha1.VarnishClusterBackend{
					Selector: map[string]string{"app": "nginx"},
					Port:     &validBackendPort,
				},
				Varnish: &icmv1alpha1.VarnishClusterVarnish{
					Args: []string{"@$invalid", "argument"},
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

		err := k8sClient.Create(context.Background(), vc)
		Expect(err).To(HaveOccurred())
		statusErr, ok := err.(*errors.StatusError)
		Expect(ok).To(BeTrue())
		Expect(statusErr.ErrStatus.Code).To(BeEquivalentTo(403))
		Expect(statusErr.ErrStatus.Status).To(Equal("Failure"))
	})
})

var _ = Describe("Mutating webhook", func() {
	validBackendPort := intstr.FromInt(8080)
	vcNamespace := "default"
	vcName := "test"
	objMeta := metav1.ObjectMeta{
		Namespace: vcNamespace,
		Name:      vcName,
	}

	AfterEach(func() {
		err := k8sClient.DeleteAllOf(context.Background(), &icmv1alpha1.VarnishCluster{}, client.InNamespace(vcNamespace))
		Expect(err).To(Succeed())
	})

	It("is working", func() {
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

		err := k8sClient.Create(context.Background(), vc)
		Expect(err).To(Succeed())

		persistedVC := &icmv1alpha1.VarnishCluster{}
		err = k8sClient.Get(context.Background(), types.NamespacedName{Name: vcName, Namespace: vcNamespace}, persistedVC)
		Expect(err).To(Succeed())
		Expect(persistedVC.Spec.Replicas).To(Equal(proto.Int(1)), "Replicas count should be set to 1 in mutating webhook")
	})
})
