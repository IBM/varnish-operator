package tests

import (
	"context"
	"github.com/gogo/protobuf/proto"
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
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
		err := k8sClient.DeleteAllOf(context.Background(), &vcapi.VarnishCluster{}, client.InNamespace(vcNamespace))
		Expect(err).To(Succeed())
		waitUntilVarnishClusterRemoved(vcName, vcNamespace)
	})

	It("is working", func() {
		vc := &vcapi.VarnishCluster{
			ObjectMeta: objMeta,
			Spec: vcapi.VarnishClusterSpec{
				Backend: &vcapi.VarnishClusterBackend{
					Selector: map[string]string{"app": "nginx"},
					Port:     &validBackendPort,
				},
				Varnish: &vcapi.VarnishClusterVarnish{
					Args: []string{"@$invalid", "argument"},
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
		err := k8sClient.DeleteAllOf(context.Background(), &vcapi.VarnishCluster{}, client.InNamespace(vcNamespace))
		Expect(err).To(Succeed())
		waitUntilVarnishClusterRemoved(vcName, vcNamespace)
	})

	It("is working", func() {
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

		err := k8sClient.Create(context.Background(), vc)
		Expect(err).To(Succeed())

		persistedVC := &vcapi.VarnishCluster{}
		err = k8sClient.Get(context.Background(), types.NamespacedName{Name: vcName, Namespace: vcNamespace}, persistedVC)
		Expect(err).To(Succeed())
		Expect(persistedVC.Spec.Replicas).To(Equal(proto.Int(1)), "Replicas count should be set to 1 in mutating webhook")
	})
})
