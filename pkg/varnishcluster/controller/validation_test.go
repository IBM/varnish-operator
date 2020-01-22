package controller

import (
	"context"
	icmv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"

	v1 "k8s.io/api/core/v1"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("Validation", func() {
	validBackendPort := intstr.FromInt(8080)
	vcNamespace := "default"
	vcName := "test"
	objMeta := metav1.ObjectMeta{
		Namespace: vcNamespace,
		Name:      vcName,
	}

	AfterEach(func() {
		CleanUpCreatedResources(vcName, vcNamespace)
	})

	Describe(".spec field", func() {
		It("should be required", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).NotTo(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.vcl field", func() {
		It("should be required", func() {
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
					//VCL: &icmv1alpha1.VarnishClusterVCL{
					//	ConfigMapName: proto.String("test"),
					//	EntrypointFileName: proto.String("test.vcl"),
					//},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.vcl.configMapName field", func() {
		It("should be required", func() {
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
						//ConfigMapName: proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.vcl.entrypointFileName field", func() {
		It("should be required", func() {
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
						ConfigMapName: proto.String("test"),
						//EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.backend field", func() {
		It("should be required", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					//Backend: &icmv1alpha1.VarnishClusterBackend{
					//	Selector: map[string]string{"app": "nginx"},
					//	Port:     &validBackendPort,
					//},
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
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.backend.selector field", func() {
		It("should be required", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						//	Port:     &validBackendPort,
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
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.backend.port field", func() {
		It("should be required", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						//Selector: map[string]string{"app": "nginx"},
						Port: &validBackendPort,
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
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service field", func() {
		It("should be required", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					//Service: &icmv1alpha1.VarnishClusterService{
					//	Port: 8081,
					//},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service.port field", func() {
		It("should be required", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						//	Port: 8081,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.vcl.entrypointFile field", func() {
		It("should match pattern", func() {
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
						EntrypointFileName: proto.String("test-invalid-file-name"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service.port field", func() {
		It("should be greater than 0", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(0),
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service.port field", func() {
		It("should be not greater than or equal 65535", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(65536),
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service.metricsPort field", func() {
		It("should be greater than 0", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port:        proto.Int32(80),
						MetricsPort: -1,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service.port field", func() {
		It("should be not greater than or equal 65535", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port:        proto.Int32(80),
						MetricsPort: 65536,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service.type with value ClusterIP", func() {
		It("should be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeClusterIP,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).To(BeNil())
		})
	})

	Describe(".spec.service.type with value NodePort", func() {
		It("should be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeNodePort,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).To(BeNil())
		})
	})

	Describe(".spec.service.type with value LoadBalancer", func() {
		It("should be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeLoadBalancer,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).To(BeNil())
		})
	})

	Describe(".spec.service.type with not supported values", func() {
		It("should not be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeExternalName,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.varnish.imagePullPolicy with `Never`", func() {
		It("should be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &icmv1alpha1.VarnishClusterVarnish{
						ImagePullPolicy: v1.PullNever,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).To(BeNil())
		})
	})

	Describe(".spec.varnish.imagePullPolicy with `Always`", func() {
		It("should be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &icmv1alpha1.VarnishClusterVarnish{
						ImagePullPolicy: v1.PullAlways,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).To(BeNil())
		})
	})

	Describe(".spec.varnish.imagePullPolicy with `IfNotPresent`", func() {
		It("should be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &icmv1alpha1.VarnishClusterVarnish{
						ImagePullPolicy: v1.PullIfNotPresent,
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).To(BeNil())
		})
	})

	Describe(".spec.varnish.imagePullPolicy field with not supported values", func() {
		It("should not be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &icmv1alpha1.VarnishClusterVarnish{
						ImagePullPolicy: "invalid",
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.varnish.logLevel with not supported value", func() {
		It("should not be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
					LogLevel: "invalid",
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.varnish.logFormat with not supported value", func() {
		It("should not be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
					LogFormat: "invalid",
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.updateStrategy.type with not supported value", func() {
		It("should not be allowed", func() {
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					UpdateStrategy: &icmv1alpha1.VarnishClusterUpdateStrategy{
						Type: "NotSupported",
					},
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &icmv1alpha1.VarnishClusterService{
						Port: proto.Int32(80),
					},
					VCL: &icmv1alpha1.VarnishClusterVCL{
						ConfigMapName:      proto.String("test"),
						EntrypointFileName: proto.String("test.vcl"),
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe("Spec with all valid fields", func() {
		It("should be created", func() {
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
			Expect(err).To(BeNil())
		})
	})

	Describe("Backend port in string format", func() {
		It("should be allowed", func() {
			backendPortStr := intstr.FromString("varnish")
			vc := &icmv1alpha1.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: icmv1alpha1.VarnishClusterSpec{
					Backend: &icmv1alpha1.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &backendPortStr,
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
			Expect(err).To(BeNil())
		})
	})
})
