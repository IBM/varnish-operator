package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/names"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("the varnish secret", func() {
	validBackendPort := intstr.FromInt(8080)
	vcNamespace := "varnish-secret"
	vcName := "test-secret"
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

	secretName := types.NamespacedName{Name: names.VarnishSecret(vc.Name), Namespace: vcNamespace}
	AfterEach(func() {
		CleanUpCreatedResources(vcName, vcNamespace)
	})

	Context("when varnishcluster is created", func() {
		It("should be created", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			secret := &v1.Secret{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), secretName, secret)
			}, time.Second*5).Should(Succeed())

			Expect(secret.Labels).To(Equal(map[string]string{
				icmv1alpha1.LabelVarnishComponent: "secret",
				icmv1alpha1.LabelVarnishOwner:     vcName,
				icmv1alpha1.LabelVarnishUID:       string(newVC.UID),
			}))

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
			Expect(secret.OwnerReferences).To(Equal(ownerReference))
		})
	})

	Context("when the secret already exists and has not empty password", func() {
		It("should not be touched", func() {

			customSecret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      names.VarnishSecret(vc.Name),
					Namespace: vcNamespace,
					Labels:    map[string]string{"label": "custom-secret"},
				},
				Data: map[string][]byte{
					varnishDefaultSecretKeyName: []byte("some preshared secret"),
				},
			}

			Eventually(func() error {
				return k8sClient.Create(context.Background(), customSecret)
			}, time.Second*5).Should(Succeed())

			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			secret := &v1.Secret{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), secretName, secret)
			}, time.Second*5).Should(Succeed())

			Expect(secret).To(Equal(customSecret))
		})
	})

	Context("when the secret already exists and has empty password", func() {
		It("should update the password", func() {
			customSecret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      names.VarnishSecret(vc.Name),
					Namespace: vcNamespace,
					Labels:    map[string]string{"label": "custom-secret"},
				},
				Data: map[string][]byte{},
			}

			Eventually(func() error {
				return k8sClient.Create(context.Background(), customSecret)
			}, time.Second*5).Should(Succeed())

			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() []byte {
				secret := &v1.Secret{}
				err := k8sClient.Get(context.Background(), secretName, secret)
				Expect(err).To(Succeed())
				return secret.Data[varnishDefaultSecretKeyName]
			}, time.Second*5).ShouldNot(BeEmpty())

		})
	})

	Context("when the secret already exists but data uninitialized", func() {
		It("should update the password", func() {
			customSecret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      names.VarnishSecret(vc.Name),
					Namespace: vcNamespace,
					Labels:    map[string]string{"label": "custom-secret"},
				},
			}

			Eventually(func() error {
				return k8sClient.Create(context.Background(), customSecret)
			}, time.Second*5).Should(Succeed())

			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() []byte {
				secret := &v1.Secret{}
				err := k8sClient.Get(context.Background(), secretName, secret)
				Expect(err).To(Succeed())
				return secret.Data[varnishDefaultSecretKeyName]
			}, time.Second*5).ShouldNot(BeEmpty())

		})
	})

	Context("when the secret already exists and has empty password", func() {
		It("should update the password and do not touch another key", func() {
			customSecret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      names.VarnishSecret(vc.Name),
					Namespace: vcNamespace,
					Labels:    map[string]string{"label": "custom-secret"},
				},
				Data: map[string][]byte{
					"horse": []byte("redish"),
				},
			}

			Eventually(func() error {
				return k8sClient.Create(context.Background(), customSecret)
			}, time.Second*5).Should(Succeed())

			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			secret := &v1.Secret{}
			Eventually(func() []byte {
				err := k8sClient.Get(context.Background(), secretName, secret)
				Expect(err).To(Succeed())
				return secret.Data[varnishDefaultSecretKeyName]
			}, time.Second*5).ShouldNot(BeEmpty())

			Expect(secret.Data["horse"]).To(Equal([]byte("redish")))

		})
	})
})

func TestNamesForInstanceSecret(t *testing.T) {
	cases := []struct {
		desc     string
		secret   *icmv1alpha1.VarnishClusterVarnishSecret
		expected []string
	}{
		{
			"default",
			nil,
			[]string{"varnish-test-varnish-secret", "secret"},
		},
		{
			"empty",
			&icmv1alpha1.VarnishClusterVarnishSecret{},
			[]string{"varnish-test-varnish-secret", "secret"},
		},
		{
			"custom secret",
			&icmv1alpha1.VarnishClusterVarnishSecret{
				SecretName: proto.String("credentials"),
				Key:        proto.String("varnish"),
			},
			[]string{"credentials", "varnish"},
		},
		{
			"custom secret no key",
			&icmv1alpha1.VarnishClusterVarnishSecret{
				SecretName: proto.String("credentials"),
			},
			[]string{"credentials", "secret"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(tt *testing.T) {
			instance := &icmapiv1alpha1.VarnishCluster{}
			instance.Name = "varnish-test"
			instance.Spec = icmapiv1alpha1.VarnishClusterSpec{
				Varnish: &icmapiv1alpha1.VarnishClusterVarnish{
					Secret: tc.secret,
				},
			}
			secretName, key := namesForInstanceSecret(instance)
			actual := []string{secretName, key}
			if !cmp.Equal(tc.expected, []string{secretName, key}) {
				t.Logf("Failed.\nDiff: \n%#v\n%#v", tc.expected, actual)
				t.Fail()
			}
		})
	}

}
