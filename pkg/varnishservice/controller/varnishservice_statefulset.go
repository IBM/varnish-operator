package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	vslabels "icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"
	"strings"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileStatefulSet(ctx context.Context, instance, instanceStatus *icmapiv1alpha1.VarnishService, serviceAccountName string, endpointSelector map[string]string, svcName string) (*appsv1.StatefulSet, map[string]string, error) {
	varnishLabels := vslabels.CombinedComponentLabels(instance, icmapiv1alpha1.VarnishComponentVarnish)
	gvk := instance.GroupVersionKind()
	var varnishImage string
	if instance.Spec.StatefulSet.Container.Image == "" {
		varnishImage = r.config.CoupledVarnishImage
	} else {
		varnishImage = instance.Spec.StatefulSet.Container.Image
	}

	varnishdArgs := strings.Join(getSanitizedVarnishArgs(&instance.Spec), " ")

	var imagePullSecrets []v1.LocalObjectReference
	if instance.Spec.StatefulSet.Container.ImagePullSecret != nil {
		imagePullSecrets = []v1.LocalObjectReference{{Name: *instance.Spec.StatefulSet.Container.ImagePullSecret}}
	}

	var updateStrategy appsv1.StatefulSetUpdateStrategy
	if instance.Spec.StatefulSet.UpdateStrategy.Type == icmapiv1alpha1.VarnishUpdateStrategyDelayedRollingUpdate {
		updateStrategy = appsv1.StatefulSetUpdateStrategy{
			Type: appsv1.OnDeleteStatefulSetStrategyType,
		}
	} else {
		updateStrategy = instance.Spec.StatefulSet.UpdateStrategy.StatefulSetUpdateStrategy
	}

	desired := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-varnish",
			Labels:    varnishLabels,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: svcName,
			Replicas:    instance.Spec.StatefulSet.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: varnishLabels,
			},
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			UpdateStrategy:       updateStrategy,
			RevisionHistoryLimit: func(in int32) *int32 { return &in }(10),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: varnishLabels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  icmapiv1alpha1.VarnishContainerName,
							Image: varnishImage,
							Ports: []v1.ContainerPort{
								{
									Name:          instance.Spec.Service.VarnishPort.Name,
									ContainerPort: icmapiv1alpha1.VarnishPort,
									Protocol:      v1.ProtocolTCP,
								},
								{
									Name:          instance.Spec.Service.VarnishExporterPort.Name,
									ContainerPort: icmapiv1alpha1.VarnishPrometheusExporterPort,
									Protocol:      v1.ProtocolTCP,
								},
							},
							Env: []v1.EnvVar{
								{Name: "ENDPOINT_SELECTOR_STRING", Value: labels.SelectorFromSet(endpointSelector).String()},
								{Name: "CONFIGMAP_NAME", Value: instance.Spec.VCLConfigMap.Name},
								{Name: "NAMESPACE", Value: instance.Namespace},
								{Name: "POD_NAME", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "metadata.name"}}},
								{Name: "VARNISH_SERVICE_NAME", Value: instance.Name},
								{Name: "VARNISH_SERVICE_UID", Value: string(instance.UID)},
								{Name: "VARNISH_SERVICE_GROUP", Value: gvk.Group},
								{Name: "VARNISH_SERVICE_VERSION", Value: gvk.Version},
								{Name: "VARNISH_SERVICE_KIND", Value: gvk.Kind},
								{Name: "LOG_FORMAT", Value: instance.Spec.LogFormat},
								{Name: "LOG_LEVEL", Value: instance.Spec.LogLevel},
								{Name: "VARNISH_ARGS", Value: varnishdArgs},
							},
							Resources: instance.Spec.StatefulSet.Container.Resources,
							// TODO: get working liveness probe
							//LivenessProbe:   &v1.Probe{},
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
							ImagePullPolicy:          instance.Spec.StatefulSet.Container.ImagePullPolicy,
						},
					},
					RestartPolicy:                 instance.Spec.StatefulSet.Container.RestartPolicy,
					TerminationGracePeriodSeconds: func(in int64) *int64 { return &in }(30),
					DNSPolicy:                     v1.DNSClusterFirst,
					SecurityContext:               &v1.PodSecurityContext{},
					ServiceAccountName:            serviceAccountName,
					Affinity:                      instance.Spec.StatefulSet.Affinity,
					Tolerations:                   instance.Spec.StatefulSet.Tolerations,
					ImagePullSecrets:              imagePullSecrets,
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

	instanceStatus.Status.StatefulSet.StatefulSetStatus = found.Status
	instanceStatus.Status.StatefulSet.Name = found.Name
	instanceStatus.Status.StatefulSet.VarnishArgs = varnishdArgs

	return found, varnishLabels, nil
}
