package controller

import (
	"context"
	"strings"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"github.com/gogo/protobuf/proto"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishCluster) reconcileStatefulSet(ctx context.Context, instance, instanceStatus *vcapi.VarnishCluster, endpointSelector map[string]string) (*appsv1.StatefulSet, map[string]string, error) {
	varnishLabels := vclabels.CombinedComponentLabels(instance, vcapi.VarnishComponentVarnish)
	gvk := instance.GroupVersionKind()
	var varnishImage string
	if instance.Spec.Varnish.Image == "" {
		varnishImage = r.config.CoupledVarnishImage
	} else {
		varnishImage = instance.Spec.Varnish.Image
	}
	varnishSecretName, varnishSecretKeyName := namesForInstanceSecret(instance)
	varnishdArgs := getSanitizedVarnishArgs(&instance.Spec)

	var updateStrategy appsv1.StatefulSetUpdateStrategy
	switch instance.Spec.UpdateStrategy.Type {
	case vcapi.OnDeleteVarnishClusterStrategyType,
		vcapi.DelayedRollingUpdateVarnishClusterStrategyType:
		updateStrategy.Type = appsv1.OnDeleteStatefulSetStrategyType
	case vcapi.RollingUpdateVarnishClusterStrategyType:
		updateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType
		updateStrategy.RollingUpdate = instance.Spec.UpdateStrategy.RollingUpdate
	}

	varnishControllerImage := imageNameGenerate(instance.Spec.Varnish.Controller.Image, varnishImage, vcapi.VarnishControllerImage)
	varnishMetricsImage := imageNameGenerate(instance.Spec.Varnish.MetricsExporter.Image, varnishImage, vcapi.VarnishMetricsExporterImage)

	desired := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.StatefulSet(instance.Name),
			Labels:    varnishLabels,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: names.HeadlessService(instance.Name),
			Replicas:    instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: varnishLabels,
			},
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			UpdateStrategy:       updateStrategy,
			RevisionHistoryLimit: proto.Int(10),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: varnishLabels,
				},
				Spec: v1.PodSpec{
					// Share a single process namespace between all of the containers in a pod.
					// When this is set containers will be able to view and signal processes from other containers
					// in the same pod.It is required for the pod to provide reliable way to collect metrics.
					// Otherwise metrics collection container may only collect general varnish process metrics.
					ShareProcessNamespace: proto.Bool(true),
					Volumes: []v1.Volume{
						{
							Name: vcapi.VarnishSharedVolume,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: vcapi.VarnishSettingsVolume,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: vcapi.VarnishSecretVolume,
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									Items: []v1.KeyToPath{
										{
											Key:  varnishSecretKeyName,
											Path: "secret",
											Mode: proto.Int32(0444), //octal mode read only
										},
									},
									DefaultMode: proto.Int32(v1.SecretVolumeSourceDefaultMode),
									SecretName:  varnishSecretName,
								},
							},
						},
						{
							Name: vcapi.HaproxyConfigVolume,
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: instance.Spec.HaproxySidecar.ConfigMapName,
									},
								},
							},
						},
					},

					Containers: []v1.Container{
						//Varnish container
						{
							Name:  vcapi.VarnishContainerName,
							Image: varnishImage,
							Ports: []v1.ContainerPort{
								{
									Name:          vcapi.VarnishPortName,
									ContainerPort: vcapi.VarnishPort,
									Protocol:      v1.ProtocolTCP,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      vcapi.VarnishSharedVolume,
									MountPath: "/var/lib/varnish",
								},
								{
									Name:      vcapi.VarnishSettingsVolume,
									MountPath: "/etc/varnish",
									ReadOnly:  true,
								},
								{
									Name:      vcapi.VarnishSecretVolume,
									MountPath: "/etc/varnish-secret",
									ReadOnly:  true,
								},
							},
							Args:      varnishdArgs,
							Resources: *instance.Spec.Varnish.Resources,
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									Exec: &v1.ExecAction{
										Command: []string{"/usr/bin/varnishadm", "ping"},
									},
								},
								TimeoutSeconds:   30,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: v1.TerminationMessageReadFile,
							ImagePullPolicy:          instance.Spec.Varnish.ImagePullPolicy,
							EnvFrom:                  instance.Spec.Varnish.EnvFrom,
						},
						//Varnish metrics
						{
							Name:  vcapi.VarnishMetricsExporterName,
							Image: varnishMetricsImage,
							Ports: []v1.ContainerPort{
								{
									Name:          vcapi.VarnishMetricsPortName,
									ContainerPort: vcapi.VarnishPrometheusExporterPort,
									Protocol:      v1.ProtocolTCP,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      vcapi.VarnishSharedVolume,
									MountPath: "/var/lib/varnish",
									ReadOnly:  true,
								},
								{
									Name:      vcapi.VarnishSettingsVolume,
									MountPath: "/etc/varnish",
									ReadOnly:  true,
								},
							},
							Resources: instance.Spec.Varnish.MetricsExporter.Resources,
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Port:   intstr.FromInt(vcapi.VarnishPrometheusExporterPort),
										Scheme: v1.URISchemeHTTP,
										Path:   "/",
									},
								},
								TimeoutSeconds:   30,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: v1.TerminationMessageReadFile,
							ImagePullPolicy:          instance.Spec.Varnish.MetricsExporter.ImagePullPolicy,
						},
						//Varnish controller
						{
							Name:  vcapi.VarnishControllerName,
							Image: varnishControllerImage,
							Ports: []v1.ContainerPort{
								{
									Name:          vcapi.VarnishControllerMetricsPortName,
									Protocol:      v1.ProtocolTCP,
									ContainerPort: vcapi.VarnishControllerMetricsPort,
								},
							},
							Env: []v1.EnvVar{
								{Name: "ENDPOINT_SELECTOR_STRING", Value: labels.SelectorFromSet(endpointSelector).String()},
								{Name: "NAMESPACE", Value: instance.Namespace},
								{Name: "POD_NAME", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"}}},
								{Name: "NODE_NAME", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "spec.nodeName"}}},
								{Name: "VARNISH_CLUSTER_NAME", Value: instance.Name},
								{Name: "VARNISH_CLUSTER_UID", Value: string(instance.UID)},
								{Name: "VARNISH_CLUSTER_GROUP", Value: gvk.Group},
								{Name: "VARNISH_CLUSTER_VERSION", Value: gvk.Version},
								{Name: "VARNISH_CLUSTER_KIND", Value: gvk.Kind},
								{Name: "LOG_FORMAT", Value: instance.Spec.LogFormat},
								{Name: "LOG_LEVEL", Value: instance.Spec.LogLevel},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      vcapi.VarnishSettingsVolume,
									MountPath: "/etc/varnish",
								},
								{
									Name:      vcapi.VarnishSharedVolume,
									MountPath: "/var/lib/varnish",
									ReadOnly:  true,
								},
								{
									Name:      vcapi.VarnishSecretVolume,
									MountPath: "/etc/varnish-secret",
									ReadOnly:  true,
								},
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Port:   intstr.FromInt(vcapi.HealthCheckPort),
										Path:   "/readyz",
										Scheme: v1.URISchemeHTTP,
									},
								},
								TimeoutSeconds:   10,
								PeriodSeconds:    3,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							Resources:                instance.Spec.Varnish.Controller.Resources,
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: v1.TerminationMessageReadFile,
							ImagePullPolicy:          instance.Spec.Varnish.Controller.ImagePullPolicy,
						},
						//haproxy sidecar
						{
							Name: vcapi.HaproxyContainerName,
							Image: instance.Spec.HaproxySidecar.Image,
							ImagePullPolicy: instance.Spec.HaproxySidecar.ImagePullPolicy,
							// apparently /healthz is only for haproxy-ingress
							//ReadinessProbe: &v1.Probe{
							//	Handler: v1.Handler{
							//		HTTPGet: &v1.HTTPGetAction{
							//			Port: intstr.FromInt(vcapi.HaproxyHealthCheckPort),
							//			Path: "/healthz",
							//			Scheme: v1.URISchemeHTTP,
							//		},
							//	},
							//	TimeoutSeconds: 10,
							//	PeriodSeconds: 10,
							//	SuccessThreshold: 1,
							//	FailureThreshold: 3,
							//	InitialDelaySeconds: 10,
							//},
							Resources: *instance.Spec.HaproxySidecar.Resources,
							VolumeMounts: []v1.VolumeMount{
								{
									Name: vcapi.HaproxyConfigVolume,
									MountPath: vcapi.HaproxyConfigMountPath,
									ReadOnly: true,
								},
							},
						},
					},
					TerminationGracePeriodSeconds: proto.Int64(30),
					DNSPolicy:                     v1.DNSClusterFirst,
					SecurityContext:               &v1.PodSecurityContext{},
					ServiceAccountName:            names.ServiceAccount(instance.Name),
					Affinity:                      instance.Spec.Affinity,
					Tolerations:                   instance.Spec.Tolerations,
					RestartPolicy:                 v1.RestartPolicyAlways,
				},
			},
		},
	}

	if instance.Spec.Varnish.ImagePullSecret != "" {
		desired.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			{
				Name: instance.Spec.Varnish.ImagePullSecret,
			},
		}
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentVarnish)
	logr = logr.With(logger.FieldComponentName, desired.Name)

	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return nil, nil, errors.Wrap(err, "could not set controller as the OwnerReference for statefulset")
	}
	r.scheme.Default(desired)

	found := &appsv1.StatefulSet{}

	err := r.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, found)
	// If the statefulset does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the statefulset exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating StatefulSet", "new", desired)
		err = r.Create(ctx, desired)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not create statefulset")
		}
	} else if err != nil {
		return nil, nil, errors.Wrap(err, "could not get current state of statefulset")
	} else {
		// the pod selector is immutable once set, so always enforce the same as existing
		desired.Spec.Selector = found.Spec.Selector
		desired.Spec.Template.Labels = found.Spec.Template.Labels
		desired.Spec.Template.Annotations = found.Spec.Template.Annotations
		if !compare.EqualStatefulSet(found, desired) {
			logr.Infoc("Updating StatefulSet", "diff", compare.DiffStatefulSet(found, desired))
			found.Spec = desired.Spec
			found.Labels = desired.Labels
			if err = r.Update(ctx, found); err != nil {
				return nil, nil, errors.Wrap(err, "could not update statefulset")
			}
		} else {
			logr.Debugw("No updates for StatefulSet")
		}
	}

	instanceStatus.Status.VarnishArgs = strings.Join(varnishdArgs, " ")
	instanceStatus.Status.Replicas = found.Status.Replicas

	return found, varnishLabels, nil
}

func imageNameGenerate(specified, base, suffix string) string {
	if specified != "" {
		return specified
	}
	baseName, tag := splitNameAndTag(base)
	return baseName + suffix + tag
}

func splitNameAndTag(fullName string) (image, tag string) {
	parts := strings.Split(fullName, ":")
	image = parts[0]
	if len(parts) > 1 {
		tag = ":" + parts[1]
	}
	return
}
