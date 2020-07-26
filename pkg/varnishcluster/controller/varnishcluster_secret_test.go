package controller

import (
	"context"
	"testing"
	"time"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/names"

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
				vcapi.LabelVarnishComponent: "secret",
				vcapi.LabelVarnishOwner:     vcName,
				vcapi.LabelVarnishUID:       string(newVC.UID),
			}))

			ownerReference := []metav1.OwnerReference{
				{
					APIVersion:         "ibm.com/v1alpha1",
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
		secret   *vcapi.VarnishClusterVarnishSecret
		expected []string
	}{
		{
			"default",
			nil,
			[]string{"varnish-test-varnish-secret", "secret"},
		},
		{
			"empty",
			&vcapi.VarnishClusterVarnishSecret{},
			[]string{"varnish-test-varnish-secret", "secret"},
		},
		{
			"custom secret",
			&vcapi.VarnishClusterVarnishSecret{
				SecretName: proto.String("credentials"),
				Key:        proto.String("varnish"),
			},
			[]string{"credentials", "varnish"},
		},
		{
			"custom secret no key",
			&vcapi.VarnishClusterVarnishSecret{
				SecretName: proto.String("credentials"),
			},
			[]string{"credentials", "secret"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(tt *testing.T) {
			instance := &vcapi.VarnishCluster{}
			instance.Name = "varnish-test"
			instance.Spec = vcapi.VarnishClusterSpec{
				Varnish: &vcapi.VarnishClusterVarnish{
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
