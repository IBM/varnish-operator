package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	componentNameVarnishes = "varnishes"
)

func (r *ReconcileVarnishService) reconcileDeployment(instance, instanceStatus *icmapiv1alpha1.VarnishService, serviceAccountName string, endpointSelector map[string]string) (map[string]string, error) {
	podSelector := generateLabels(instance, componentNameVarnishes)
	gvk := instance.GroupVersionKind()
	desired := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-varnish-deployment",
			Labels:    combinedLabels(instance, componentNameVarnishes),
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: instance.Spec.Deployment.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: podSelector,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podSelector,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "varnish",
							Image: instance.Spec.Deployment.Container.Image,
							Ports: []v1.ContainerPort{
								{
									Name:          instance.Spec.Service.VarnishPort.Name,
									HostPort:      instance.Spec.Service.VarnishPort.Port,
									ContainerPort: instance.Spec.Service.VarnishPort.Port,
								},
								{
									Name:          instance.Spec.Service.VarnishExporterPort.Name,
									HostPort:      instance.Spec.Service.VarnishExporterPort.Port,
									ContainerPort: instance.Spec.Service.VarnishExporterPort.Port,
								},
							},
							Env: []v1.EnvVar{
								{Name: "ENDPOINT_SELECTOR_STRING", Value: labels.SelectorFromSet(endpointSelector).String()},
								{Name: "CONFIGMAP_NAME", Value: instance.Spec.VCLConfigMap.Name},
								{Name: "BACKENDS_FILE", Value: instance.Spec.VCLConfigMap.BackendsFile},
								{Name: "BACKENDS_TMPL_FILE", Value: instance.Spec.VCLConfigMap.BackendsTmplFile},
								{Name: "DEFAULT_FILE", Value: instance.Spec.VCLConfigMap.DefaultFile},
								{Name: "NAMESPACE", Value: instance.Namespace},
								{Name: "POD_NAME", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.name"}}},
								{Name: "VARNISH_SERVICE_NAME", Value: instance.Name},
								{Name: "VARNISH_SERVICE_UID", Value: string(instance.UID)},
								{Name: "VARNISH_SERVICE_GROUP", Value: gvk.Group},
								{Name: "VARNISH_SERVICE_VERSION", Value: gvk.Version},
								{Name: "VARNISH_SERVICE_KIND", Value: gvk.Kind},
								{Name: "TARGET_PORT", Value: instance.Spec.Service.VarnishPort.TargetPort.String()},
								{Name: "LOG_FORMAT", Value: instance.Spec.LogFormat},
								{Name: "LOG_LEVEL", Value: instance.Spec.LogLevel},
								{Name: "VARNISH_ARGS", Value: strings.Join(instance.Spec.Deployment.Container.VarnishArgs, " ")},
							},
							Resources:       *instance.Spec.Deployment.Container.Resources,
							LivenessProbe:   instance.Spec.Deployment.Container.LivenessProbe,
							ReadinessProbe:  instance.Spec.Deployment.Container.ReadinessProbe,
							ImagePullPolicy: *instance.Spec.Deployment.Container.ImagePullPolicy,
						},
					},
					RestartPolicy:      instance.Spec.Deployment.Container.RestartPolicy,
					ServiceAccountName: serviceAccountName,
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: *instance.Spec.Deployment.Container.ImagePullSecret,
						},
					},
					Affinity:    instance.Spec.Deployment.Affinity,
					Tolerations: instance.Spec.Deployment.Tolerations,
				},
			},
		},
	}

	logr := r.logger.With("name", desired.Name, "namespace", desired.Namespace)

	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return nil, logr.RErrorw(err, "could not set controller as the OwnerReference for deployment")
	}
	r.scheme.Default(desired)

	found := &appsv1.Deployment{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, found)
	// If the deployment does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the deployment exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Deployment", "new", desired)
		err = r.Create(context.TODO(), desired)
		if err != nil {
			return nil, logr.RErrorw(err, "could not create deployment")
		}
	} else if err != nil {
		return nil, logr.RErrorw(err, "could not get current state of deployment")
	} else {
		// the pod selector is immutable once set, so always enforce the same as existing
		desired.Spec.Selector = found.Spec.Selector
		desired.Spec.Template.Labels = found.Spec.Template.Labels
		if !compare.EqualDeployment(found, desired) {
			logr.Infoc("Updating Deployment", "diff", compare.DiffDeployment(found, desired))
			found.Spec = desired.Spec
			found.Labels = desired.Labels
			if err = r.Update(context.TODO(), found); err != nil {
				return nil, logr.RErrorw(err, "could not update deployment")
			}
		} else {
			logr.Debugw("No updates for Deployment")
		}
	}

	instanceStatus.Status.Deployment = found.Status

	return podSelector, nil
}
