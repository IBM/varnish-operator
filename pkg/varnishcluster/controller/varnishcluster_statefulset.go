package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	vclabels "icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/names"
	"icm-varnish-k8s-operator/pkg/varnishcluster/compare"
	"strings"

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

func (r *ReconcileVarnishCluster) reconcileStatefulSet(ctx context.Context, instance, instanceStatus *icmapiv1alpha1.VarnishCluster, endpointSelector map[string]string) (*appsv1.StatefulSet, map[string]string, error) {
	varnishLabels := vclabels.CombinedComponentLabels(instance, icmapiv1alpha1.VarnishComponentVarnish)
	gvk := instance.GroupVersionKind()
	var varnishImage string
	if instance.Spec.Varnish.Image == "" {
		varnishImage = r.config.CoupledVarnishImage
	} else {
		varnishImage = instance.Spec.Varnish.Image
	}

	varnishdArgs := getSanitizedVarnishArgs(&instance.Spec)

	var imagePullSecrets []v1.LocalObjectReference
	if instance.Spec.Varnish.ImagePullSecret != nil {
		imagePullSecrets = []v1.LocalObjectReference{{Name: *instance.Spec.Varnish.ImagePullSecret}}
	}

	var updateStrategy appsv1.StatefulSetUpdateStrategy
	switch instance.Spec.UpdateStrategy.Type {
	case icmapiv1alpha1.OnDeleteVarnishClusterStrategyType,
		icmapiv1alpha1.DelayedRollingUpdateVarnishClusterStrategyType:
		updateStrategy.Type = appsv1.OnDeleteStatefulSetStrategyType
	case icmapiv1alpha1.RollingUpdateVarnishClusterStrategyType:
		updateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType
		updateStrategy.RollingUpdate = instance.Spec.UpdateStrategy.RollingUpdate
	}

	varnishControllerImage := imageNameGenerate(instance.Spec.Varnish.Controller.Image, varnishImage, icmapiv1alpha1.VarnishControllerImage)
	varnishMetricsImage := imageNameGenerate(instance.Spec.Varnish.MetricsExporter.Image, varnishImage, icmapiv1alpha1.VarnishMetricsExporterImage)

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
							Name: icmapiv1alpha1.VarnishSharedVolume,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: icmapiv1alpha1.VarnishSettingsVolume,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: []v1.Container{
						{
							Name:  icmapiv1alpha1.VarnishContainerName + "-init",
							Image: varnishImage,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      icmapiv1alpha1.VarnishSettingsVolume,
									MountPath: "/data",
								},
							},
							Command:                  []string{"/bin/bash", "-c"},
							Args:                     []string{"echo $VARNISH_SECRET > /data/secret"},
							ImagePullPolicy:          instance.Spec.Varnish.ImagePullPolicy,
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: v1.TerminationMessageReadFile,
						},
					},
					Containers: []v1.Container{
						//Varnish container
						{
							Name:  icmapiv1alpha1.VarnishContainerName,
							Image: varnishImage,
							Ports: []v1.ContainerPort{
								{
									Name:          icmapiv1alpha1.VarnishPortName,
									ContainerPort: icmapiv1alpha1.VarnishPort,
									Protocol:      v1.ProtocolTCP,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      icmapiv1alpha1.VarnishSharedVolume,
									MountPath: "/var/lib/varnish",
								},
								{
									Name:      icmapiv1alpha1.VarnishSettingsVolume,
									MountPath: "/etc/varnish",
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
						},
						//Varnish metrics
						{
							Name:  icmapiv1alpha1.VarnishMetricsExporterName,
							Image: varnishMetricsImage,
							Ports: []v1.ContainerPort{
								{
									Name:          icmapiv1alpha1.VarnishMetricsPortName,
									ContainerPort: icmapiv1alpha1.VarnishPrometheusExporterPort,
									Protocol:      v1.ProtocolTCP,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      icmapiv1alpha1.VarnishSharedVolume,
									MountPath: "/var/lib/varnish",
									ReadOnly:  true,
								},
								{
									Name:      icmapiv1alpha1.VarnishSettingsVolume,
									MountPath: "/etc/varnish",
									ReadOnly:  true,
								},
							},
							Resources: instance.Spec.Varnish.MetricsExporter.Resources,
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Port:   intstr.FromInt(icmapiv1alpha1.VarnishPrometheusExporterPort),
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
							Name:  icmapiv1alpha1.VarnishControllerName,
							Image: varnishControllerImage,
							Env: []v1.EnvVar{
								{Name: "ENDPOINT_SELECTOR_STRING", Value: labels.SelectorFromSet(endpointSelector).String()},
								{Name: "CONFIGMAP_NAME", Value: *instance.Spec.VCL.ConfigMapName},
								{Name: "NAMESPACE", Value: instance.Namespace},
								{Name: "POD_NAME", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"}}},
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
									Name:      icmapiv1alpha1.VarnishSettingsVolume,
									MountPath: "/etc/varnish",
								},
								{
									Name:      icmapiv1alpha1.VarnishSharedVolume,
									MountPath: "/var/lib/varnish",
									ReadOnly:  true,
								},
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Port: intstr.FromInt(icmapiv1alpha1.HealthCheckPort),
										Path: "/readyz",
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
					},
					TerminationGracePeriodSeconds: proto.Int64(30),
					DNSPolicy:                     v1.DNSClusterFirst,
					SecurityContext:               &v1.PodSecurityContext{},
					ServiceAccountName:            names.ServiceAccount(instance.Name),
					Affinity:                      instance.Spec.Affinity,
					Tolerations:                   instance.Spec.Tolerations,
					ImagePullSecrets:              imagePullSecrets,
					RestartPolicy:                 v1.RestartPolicyAlways,
				},
			},
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentVarnish)
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
