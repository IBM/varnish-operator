package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"

	vcapi "github.com/cin/varnish-operator/api/v1alpha1"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ = Describe("statefulset", func() {
	validBackendPort := intstr.FromInt(8080)
	vcNamespace := "default"
	vcName := "test"
	objMeta := metav1.ObjectMeta{
		Namespace: vcNamespace,
		Name:      vcName,
		Labels:    map[string]string{"custom": "label"},
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

	stsName := types.NamespacedName{Name: vcName + "-varnish", Namespace: vcNamespace}

	AfterEach(func() {
		CleanUpCreatedResources(vcName, vcNamespace)
	})

	Context("when varnishcluster is created with minimal configuration", func() {
		It("should be created with correct default values", func() {
			newVC := vc.DeepCopy()
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			expectedLabels := map[string]string{
				"custom":                    "label",
				vcapi.LabelVarnishOwner:     vcName,
				vcapi.LabelVarnishComponent: vcapi.VarnishComponentVarnish,
				vcapi.LabelVarnishUID:       string(newVC.UID),
			}

			sts := &apps.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), stsName, sts)
			}, time.Second*5).Should(Succeed())

			Expect(sts.Labels).To(Equal(expectedLabels))
			Expect(sts.OwnerReferences[0].UID).To(Equal(newVC.UID))
			Expect(sts.OwnerReferences[0].Controller).To(Equal(proto.Bool(true)))
			Expect(sts.OwnerReferences[0].Kind).To(Equal("VarnishCluster"))
			Expect(sts.OwnerReferences[0].APIVersion).To(Equal("caching.ibm.com/v1alpha1"))
			Expect(sts.OwnerReferences[0].Name).To(Equal(newVC.Name))
			Expect(sts.OwnerReferences[0].BlockOwnerDeletion).To(Equal(proto.Bool(true)))

			Expect(sts.Spec.UpdateStrategy.Type).To(Equal(apps.OnDeleteStatefulSetStrategyType))
			Expect(sts.Spec.Replicas).To(Equal(proto.Int32(1)))
			Expect(sts.Spec.ServiceName).To(Equal("test-headless-service"))
			Expect(sts.Spec.PodManagementPolicy).To(Equal(apps.ParallelPodManagement))
			Expect(sts.Spec.RevisionHistoryLimit).To(Equal(proto.Int32(10)))
			Expect(sts.Spec.Selector.MatchLabels).To(Equal(expectedLabels))

			Expect(sts.Spec.Template.Labels).To(Equal(expectedLabels))

			podSpec := sts.Spec.Template.Spec
			Expect(podSpec.ShareProcessNamespace).To(Equal(proto.Bool(true)))
			Expect(podSpec.TerminationGracePeriodSeconds).To(Equal(proto.Int64(30)))
			Expect(podSpec.DNSPolicy).To(Equal(v1.DNSClusterFirst))
			Expect(podSpec.ServiceAccountName).To(Equal("test-varnish-serviceaccount"))
			Expect(podSpec.Affinity).To(BeNil())
			Expect(podSpec.Tolerations).To(BeNil())
			Expect(podSpec.RestartPolicy).To(Equal(v1.RestartPolicyAlways))

			varnishContainer, err := getContainerByName(podSpec, vcapi.VarnishContainerName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishContainer.Args).To(Equal([]string{
				"-F",
				"-S", "/etc/varnish-secret/secret",
				"-T", "0.0.0.0:6082",
				"-a", "0.0.0.0:6081",
				"-b", "127.0.0.1:0",
			}))
			Expect(varnishContainer.Image).To(Equal(testCoupledVarnishImage))
			Expect(varnishContainer.Resources).ToNot(BeNil(), "kubernetes will set to empty struct if nil and we will infinitely fight with kubernetes by resetting it to nil")
			varnishPort, err := getContainerPortByName(varnishContainer, vcapi.VarnishPortName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishPort.ContainerPort).To(Equal(int32(vcapi.VarnishPort)))
			Expect(varnishPort.Protocol).To(Equal(v1.ProtocolTCP))

			varnishControllerContainer, err := getContainerByName(podSpec, vcapi.VarnishControllerName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishControllerContainer.Image).To(Equal("cinple/varnish-controller:test"))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "NAMESPACE", Value: vcNamespace}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "POD_NAME", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"}}}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "VARNISH_CLUSTER_NAME", Value: vcName}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "VARNISH_CLUSTER_UID", Value: string(newVC.UID)}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "VARNISH_CLUSTER_GROUP", Value: "caching.ibm.com"}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "VARNISH_CLUSTER_VERSION", Value: "v1alpha1"}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "VARNISH_CLUSTER_KIND", Value: "VarnishCluster"}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "LOG_FORMAT", Value: "json"}))
			Expect(varnishControllerContainer.Env).To(ContainElement(v1.EnvVar{Name: "LOG_LEVEL", Value: "info"}))

			metricsContainer, err := getContainerByName(podSpec, vcapi.VarnishMetricsExporterName)
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsContainer.Image).To(Equal("cinple/varnish-metrics-exporter:test"))
			varnishMetricsExporterPort, err := getContainerPortByName(metricsContainer, vcapi.VarnishMetricsPortName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishMetricsExporterPort.ContainerPort).To(Equal(int32(vcapi.VarnishPrometheusExporterPort)))
			Expect(varnishMetricsExporterPort.Protocol).To(Equal(v1.ProtocolTCP))
		})
	})

	Context("when varnishcluster is created with additional varnish args", func() {
		It("should be created with additional varnish args included", func() {
			newVC := vc.DeepCopy()
			newVC.Spec.Varnish = &vcapi.VarnishClusterVarnish{
				Args: []string{"-p", "default_ttl=3600", "-p", "default_grace=3600"},
			}
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			sts := &apps.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), stsName, sts)
			}, time.Second*5).Should(Succeed())

			varnishContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishContainerName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishContainer.Args).To(Equal([]string{
				"-F",
				"-S",
				"/etc/varnish-secret/secret",
				"-T",
				"0.0.0.0:6082",
				"-a",
				"0.0.0.0:6081",
				"-b",
				"127.0.0.1:0",
				"-p",
				"default_grace=3600",
				"-p",
				"default_ttl=3600",
			}))
		})
	})

	Context("when varnishcluster is created with non default container images", func() {
		It("should be created with overridden container images", func() {
			newVC := vc.DeepCopy()
			varnishImage := "us.icr.io/different-location/varnish:test"
			varnishControllerImage := "us.icr.io/other-location/varnish-controller:test"
			varnishMetricsExporterImage := "us.icr.io/an-another-location/varnish-metrics-exporter:test"
			newVC.Spec.Varnish = &vcapi.VarnishClusterVarnish{
				Image: varnishImage,
			}
			newVC.Spec.Varnish.Controller = &vcapi.VarnishClusterVarnishController{Image: varnishControllerImage}
			newVC.Spec.Varnish.MetricsExporter = &vcapi.VarnishClusterVarnishMetricsExporter{Image: varnishMetricsExporterImage}
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			sts := &apps.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), stsName, sts)
			}, time.Second*5).Should(Succeed())

			varnishContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishContainerName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishContainer.Image).To(Equal(varnishImage))

			varnishControllerContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishControllerName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishControllerContainer.Image).To(Equal(varnishControllerImage))

			metricsContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishMetricsExporterName)
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsContainer.Image).To(Equal(varnishMetricsExporterImage))
		})
	})

	Context("when varnishcluster is created with non default varnish container images", func() {
		It("should be created with overridden container images derived from varnish image", func() {
			newVC := vc.DeepCopy()
			varnishImage := "us.icr.io/different-location/varnish:test"
			newVC.Spec.Varnish = &vcapi.VarnishClusterVarnish{
				Image: varnishImage,
			}
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			sts := &apps.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), stsName, sts)
			}, time.Second*5).Should(Succeed())

			varnishContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishContainerName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishContainer.Image).To(Equal(varnishImage))

			varnishControllerContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishControllerName)
			Expect(err).ToNot(HaveOccurred())
			Expect(varnishControllerContainer.Image).To(Equal("us.icr.io/different-location/varnish-controller:test"))

			metricsContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishMetricsExporterName)
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsContainer.Image).To(Equal("us.icr.io/different-location/varnish-metrics-exporter:test"))
		})
	})
	Context("when varnishcluster is created with priorityClassName enable", func() {
		It("should be created with corresponding PriorityClassName set", func() {
			newVC := vc.DeepCopy()
			priorityClass := "test-priorityClass"
			newVC.Spec.PriorityClassName = priorityClass

			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			sts := &apps.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), stsName, sts)
			}, time.Second*5).Should(Succeed())

			Expect(sts.Spec.Template.Spec.PriorityClassName).To(Equal(priorityClass))

		})

	})

	Context("when varnishcluster is created with persistence enabled", func() {
		It("should be created with corresponding volume mounts and volume claim templates", func() {
			newVC := vc.DeepCopy()
			volumeMode := v1.PersistentVolumeFilesystem
			textExtraVolumeClaimTemplates := []vcapi.PVC{
				{
					Metadata: vcapi.ObjectMetadata{
						Name: "data-volume",
					},
					Spec: v1.PersistentVolumeClaimSpec{
						StorageClassName: proto.String("storage-class-example"),
						AccessModes: []v1.PersistentVolumeAccessMode{
							v1.ReadWriteOnce,
						},
						Resources: v1.ResourceRequirements{
							Requests: map[v1.ResourceName]resource.Quantity{
								v1.ResourceStorage: resource.MustParse("30Gi"),
							},
						},
						VolumeMode: &volumeMode,
					},
				},
			}

			testExtraVolume := v1.Volume{
				Name: "extra-volume",
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{},
				},
			}

			testExtraVolumeMount := v1.VolumeMount{
				Name:      "extra-volume",
				MountPath: "/extra/volume/mount",
			}

			testExtraInitContainers := []v1.Container{
				{
					Name:                     "init",
					Image:                    "ubuntu",
					ImagePullPolicy:          v1.PullIfNotPresent,
					TerminationMessagePath:   "/message/path",
					TerminationMessagePolicy: v1.TerminationMessageReadFile,
					Args:                     []string{"bash", "-c", "sleep 30s"},
				},
			}
			nodeSelector := map[string]string{
				"os": "linux",
			}

			newVC.Spec.Varnish = &vcapi.VarnishClusterVarnish{
				ExtraVolumeMounts: []v1.VolumeMount{
					testExtraVolumeMount,
				},
				ExtraVolumes: []v1.Volume{
					testExtraVolume,
				},
				ExtraVolumeClaimTemplates: textExtraVolumeClaimTemplates,
				ExtraInitContainers:       testExtraInitContainers,
			}
			newVC.Spec.NodeSelector = nodeSelector
			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			sts := &apps.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), stsName, sts)
			}, time.Second*5).Should(Succeed())

			varnishContainer, err := getContainerByName(sts.Spec.Template.Spec, vcapi.VarnishContainerName)
			Expect(err).ToNot(HaveOccurred())

			Expect(sts.Spec.VolumeClaimTemplates[0].Spec).To(Equal(textExtraVolumeClaimTemplates[0].Spec))
			Expect(sts.Spec.VolumeClaimTemplates[0].Name).To(Equal(textExtraVolumeClaimTemplates[0].Metadata.Name))
			Expect(sts.Spec.Template.Spec.Volumes).To(ContainElement(testExtraVolume))
			Expect(sts.Spec.Template.Spec.InitContainers).To(Equal(testExtraInitContainers))
			Expect(sts.Spec.Template.Spec.NodeSelector).To(Equal(nodeSelector))
			Expect(varnishContainer.VolumeMounts).To(ContainElement(testExtraVolumeMount))
		})
	})

	Context("when varnishcluster is created with podAnnotations enable", func() {
		It("should be created with corresponding podAnnotations set", func() {
			newVC := vc.DeepCopy()
			podAnnotations := map[string]string{"foo": "bar"}
			newVC.Spec.PodAnnotations = podAnnotations

			err := k8sClient.Create(context.Background(), newVC)
			Expect(err).ToNot(HaveOccurred())

			sts := &apps.StatefulSet{}
			Eventually(func() error {
				return k8sClient.Get(context.Background(), stsName, sts)
			}, time.Second*5).Should(Succeed())

			Expect(sts.Spec.Template.ObjectMeta.Annotations).To(Equal(map[string]string{"foo": "bar"}))
		})

	})
})

func getContainerByName(spec v1.PodSpec, name string) (v1.Container, error) {
	for i, container := range spec.Containers {
		if container.Name == name {
			return spec.Containers[i], nil
		}
	}

	return v1.Container{}, fmt.Errorf("container %q not found", name)
}

func getContainerPortByName(container v1.Container, name string) (v1.ContainerPort, error) {
	for i, port := range container.Ports {
		if port.Name == name {
			return container.Ports[i], nil
		}
	}

	return v1.ContainerPort{}, fmt.Errorf("container port %q not found", name)
}
