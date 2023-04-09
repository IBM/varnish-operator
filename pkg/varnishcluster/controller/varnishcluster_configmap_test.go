package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	vcapi "github.com/cin/varnish-operator/api/v1alpha1"
	vclabels "github.com/cin/varnish-operator/pkg/labels"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("the ConfigMap", func() {
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

	cmName := types.NamespacedName{Name: "test", Namespace: vcNamespace}

	AfterEach(func() {
		CleanUpCreatedResources(vcName, vcNamespace)
	})

	Context("when varnishcluster is created", func() {
		It("should be created", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), cmName, cm)
			}, time.Second*5).Should(Succeed())
			cmLabels := vclabels.CombinedComponentLabels(newVC, vcapi.VarnishComponentVCLFileConfigMap)

			Expect(cm.Labels).To(Equal(cmLabels))
		})
	})

	Context("when contents updated manually", func() {
		It("Should be reconciled and Status updated", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), cmName, cm)
			}, time.Second*5).Should(Succeed())

			oldVersion := newVC.Status.VCL.ConfigMapVersion
			By("after a change in the config")
			Eventually(func() error {
				cm.Data["test.vcl"] = strings.Replace(cm.Data["test.vcl"], "set resp.http.X-Varnish-Cache = \"HIT\";", "set resp.http.X-Varnish-Cache = \"HITit\";", 1)
				return k8sClient.Update(context.Background(), cm)
			}, time.Second*10).Should(Succeed())

			By("ConfigMapVersion should be updated")
			Eventually(func() string {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: vcName, Namespace: vcNamespace}, newVC)
				Expect(err).ToNot(HaveOccurred())
				return newVC.Status.VCL.ConfigMapVersion
			}, time.Second*10).ShouldNot(Equal(oldVersion))
		})
	})

	Context("when name updated manually", func() {
		It("Should be reconciled and Status updated", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())
			cm := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), cmName, cm)
			}, time.Second*5).Should(Succeed())

			newCMName := "newtest"
			By("the ConfigMapVersion should be changed")
			Eventually(func() error {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: vcName, Namespace: vcNamespace}, newVC)
				if err != nil {
					return err
				}
				newVC.Spec.VCL.ConfigMapName = &newCMName
				err = k8sClient.Update(context.Background(), newVC)
				return err
			}, time.Second*10).Should(Succeed())

			newCM := &v1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), types.NamespacedName{Name: newCMName, Namespace: vcNamespace}, newCM)
			}, time.Second*5).Should(Succeed())

			Eventually(func() string {
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: vcName, Namespace: vcNamespace}, newVC)
				if err != nil {
					return fmt.Sprintf("can't get the ConfigMap: %s", err.Error())
				}

				return newVC.Status.VCL.ConfigMapVersion
			}, time.Second*10).Should(Equal(newCM.ResourceVersion))
		})
	})

})
