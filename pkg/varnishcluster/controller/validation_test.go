package controller

import (
	"context"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"

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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).NotTo(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.vcl field", func() {
		It("should be required", func() {
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
					//VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					//Backend: &vcapi.VarnishClusterBackend{
					//	Selector: map[string]string{"app": "nginx"},
					//	Port:     &validBackendPort,
					//},
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
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.backend.selector field", func() {
		It("should be required", func() {
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						//	Port:     &validBackendPort,
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
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.backend.port field", func() {
		It("should be required", func() {
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						//Selector: map[string]string{"app": "nginx"},
						Port: &validBackendPort,
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
			Expect(err).ToNot(BeNil())
			Expect(err.(*errors.StatusError).ErrStatus.Reason == metav1.StatusReasonInvalid).To(BeTrue())
		})
	})

	Describe(".spec.service field", func() {
		It("should be required", func() {
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					//Service: &vcapi.VarnishClusterService{
					//	Port: 8081,
					//},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						//	Port: 8081,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
						EntrypointFileName: proto.String("test-invalid-file-name"),
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeClusterIP,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeNodePort,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeLoadBalancer,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
						Type: v1.ServiceTypeExternalName,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &vcapi.VarnishClusterVarnish{
						ImagePullPolicy: v1.PullNever,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &vcapi.VarnishClusterVarnish{
						ImagePullPolicy: v1.PullAlways,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &vcapi.VarnishClusterVarnish{
						ImagePullPolicy: v1.PullIfNotPresent,
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
					},
					Varnish: &vcapi.VarnishClusterVarnish{
						ImagePullPolicy: "invalid",
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					UpdateStrategy: &vcapi.VarnishClusterUpdateStrategy{
						Type: "NotSupported",
					},
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &validBackendPort,
					},
					Service: &vcapi.VarnishClusterService{
						Port: proto.Int32(80),
					},
					VCL: &vcapi.VarnishClusterVCL{
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
			Expect(err).To(BeNil())
		})
	})

	Describe("Backend port in string format", func() {
		It("should be allowed", func() {
			backendPortStr := intstr.FromString("varnish")
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &backendPortStr,
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
			Expect(err).To(BeNil())
		})
	})

	Describe("grafana datasourceName field with dashboard enabled", func() {
		It("should be required", func() {
			backendPortStr := intstr.FromString("varnish")
			vc := &vcapi.VarnishCluster{
				ObjectMeta: objMeta,
				Spec: vcapi.VarnishClusterSpec{
					Backend: &vcapi.VarnishClusterBackend{
						Selector: map[string]string{"app": "nginx"},
						Port:     &backendPortStr,
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
							Enabled: true,
						},
					},
				},
			}

			err := k8sClient.Create(context.Background(), vc)
			Expect(err).To(HaveOccurred())
		})
	})
})
