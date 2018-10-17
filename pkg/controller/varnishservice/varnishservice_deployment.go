package varnishservice

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/compare"
	"icm-varnish-k8s-operator/pkg/config"
	"icm-varnish-k8s-operator/pkg/logger"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	kappsv1 "k8s.io/kubernetes/pkg/apis/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileDeployment(instance, instanceStatus *icmapiv1alpha1.VarnishService, serviceAccountName string, applicationPort *v1.ServicePort, endpointSelector map[string]string) (map[string]string, error) {
	componentName := "varnishes"
	podSelector := generateLabels(instance, componentName)
	desired := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-varnish-deployment",
			Labels:    combinedLabels(instance, componentName),
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
							Image: config.GlobalConf.VarnishImageFullPath,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: config.GlobalConf.VarnishPort,
								},
								{
									ContainerPort: config.GlobalConf.VarnishExporterTargetPort,
								},
							},
							Env: []v1.EnvVar{
								{Name: "ENDPOINT_SELECTOR_STRING", Value: labels.SelectorFromSet(endpointSelector).String()},
								{Name: "BACKENDS_FILE", Value: instance.Spec.Deployment.BackendsFile},
								{Name: "DEFAULT_FILE", Value: instance.Spec.Deployment.DefaultFile},
								{Name: "NAMESPACE", Value: instance.Namespace},
								{Name: "TARGET_PORT", Value: applicationPort.TargetPort.String()},
								{Name: "VARNISH_PORT", Value: strconv.FormatInt(int64(config.GlobalConf.VarnishPort), 10)},
								{Name: "VARNISH_MEMORY", Value: instance.Spec.Deployment.VarnishMemory},
								{Name: "VCL_DIR", Value: config.GlobalConf.VCLDir},
							},
							Resources:       *instance.Spec.Deployment.VarnishResources,
							LivenessProbe:   instance.Spec.Deployment.LivenessProbe,
							ReadinessProbe:  instance.Spec.Deployment.ReadinessProbe,
							ImagePullPolicy: config.GlobalConf.VarnishImagePullPolicy,
						},
					},
					RestartPolicy:      instance.Spec.Deployment.VarnishRestartPolicy,
					ServiceAccountName: serviceAccountName,
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: instance.Spec.Deployment.ImagePullSecretName,
						},
					},
					Affinity:    instance.Spec.Deployment.Affinity,
					Tolerations: instance.Spec.Deployment.Tolerations,
				},
			},
		},
	}

	logr := logger.WithValues("name", desired.Name, "namespace", desired.Namespace)

	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return nil, logr.RErrorw(err, "could not set controller as the OwnerReference for deployment")
	}
	kappsv1.SetObjectDefaults_Deployment(desired)

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
			logr.Debugw("Updating Deployment", "diff", compare.DiffDeployment(found, desired))
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
