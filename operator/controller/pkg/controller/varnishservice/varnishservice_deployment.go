package varnishservice

import (
	"context"
	"errors"
	icmapiv1alpha1 "icm-varnish-k8s-operator/operator/controller/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/operator/controller/pkg/config"
	"icm-varnish-k8s-operator/operator/controller/pkg/logger"
	"reflect"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileDeployment(instance *icmapiv1alpha1.VarnishService, serviceAccountName string, applicationPort *v1.ServicePort) (*map[string]string, error) {
	deployConfig, err := newVarnishDeploymentConfig(r.globalConf, instance, serviceAccountName, applicationPort)
	if err != nil {
		return nil, logger.RError(err, "could not generate deployment config")
	}
	deploy, err := newVarnishDeployment(r.globalConf, deployConfig)
	if err != nil {
		return nil, logger.RError(err, "could not generate deployment")
	}
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return nil, logger.RError(err, "could not set controller as the OwnerReference for deployment", "name", deploy.Name, "namespace", deploy.Namespace)
	}

	found := &appsv1.Deployment{}
	if err = r.reconcileDeploymentStatus(instance, found); err != nil {
		return nil, err
	}

	err = r.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && kerrors.IsNotFound(err) {
		// logger.Info("Creating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		logger.Info("Creating Deployment", "config", deploy)
		err = r.Create(context.TODO(), deploy)
		if err != nil {
			return nil, logger.RError(err, "could not create deployment", "name", deploy.Name, "namespace", deploy.Namespace)
		}
	} else if err != nil {
		return nil, logger.RError(err, "could not get current state of deployment", "name", deploy.Name, "namespace", deploy.Namespace)
	} else if !reflect.DeepEqual(deploy.Spec, found.Spec) {
		found.Spec = deploy.Spec
		// logger.Info("Updating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
		logger.Info("Updating Deployment", "config", found)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return nil, logger.RError(err, "could not update deployment", "name", deploy.Name, "namespace", deploy.Namespace)
		}
	}
	logger.Info("No updates for Deployment")

	return &deployConfig.Labels, nil
}

type varnishDeploymentConfig struct {
	Name                 string
	Namespace            string
	Labels               map[string]string
	AppSelector          map[string]string
	VarnishRestartPolicy *v1.RestartPolicy
	VarnishReplicas      int32
	VarnishMemory        string
	BackendsFile         string
	DefaultFile          string
	Resources            *v1.ResourceRequirements
	ExporterResources    *v1.ResourceRequirements
	VolumeMountName      string
	VolumeMountPath      string
	LivenessProbe        *v1.Probe
	ReadinessProbe       *v1.Probe
	ServiceAccountName   string
	Port                 v1.ServicePort
	ImagePullSecretName  string
	Affinity             *v1.Affinity
	Tolerations          []v1.Toleration
}

// it appears this function does nothing in ibm cloud. Meaning, they have disabled status for custom resources. Leaving it here for now, however.
func (r *ReconcileVarnishService) reconcileDeploymentStatus(instance *icmapiv1alpha1.VarnishService, currentDeployment *appsv1.Deployment) error {
	if !reflect.DeepEqual(instance.Status.Deployment, currentDeployment.Status) {
		instance.Status.Deployment = currentDeployment.Status
		if err := r.Status().Update(context.TODO(), instance); err != nil {
			return logger.RError(err, "Could not update deployment status", "name", instance.Name, "namespace", instance.Namespace)
		}
	}
	return nil
}

func newVarnishDeploymentConfig(globalConf *config.Config, vs *icmapiv1alpha1.VarnishService, serviceAccountName string, applicationPort *v1.ServicePort) (*varnishDeploymentConfig, error) {
	vdc := varnishDeploymentConfig{
		ServiceAccountName: serviceAccountName,
		Port:               *applicationPort,
		Affinity:           vs.Spec.Deployment.Affinity,
		Tolerations:        vs.Spec.Deployment.Tolerations,
	}
	// required fields
	if vdc.Name = vs.Name + "-deployment"; vdc.Name == "-deployment" {
		return &vdc, errors.New("name not defined")
	}

	vdc.Labels = map[string]string{"component": vs.Name + "-varnish"}

	if vdc.AppSelector = vs.Spec.Service.Selector; len(vdc.AppSelector) == 0 {
		return &vdc, errors.New("must have selector to target application backed by Varnish")
	}

	if vdc.VarnishReplicas = vs.Spec.Deployment.Replicas; vdc.VarnishReplicas == 0 {
		return &vdc, errors.New("replicas not defined")
	}
	if vdc.Namespace = vs.Namespace; vdc.Namespace == "" {
		return &vdc, errors.New("namespace not defined")
	}

	if vdc.ImagePullSecretName = vs.Spec.Deployment.ImagePullSecretName; vdc.ImagePullSecretName == "" {
		return &vdc, errors.New("ImagePullSecretName not defined")
	}

	// optional fields
	if vdc.VarnishMemory = vs.Spec.Deployment.VarnishMemory; vdc.VarnishMemory == "" {
		vdc.VarnishMemory = globalConf.DefaultVarnishMemory
	}
	if vdc.BackendsFile = vs.Spec.Deployment.BackendsFile; vdc.BackendsFile == "" {
		vdc.BackendsFile = globalConf.DefaultBackendsFile
	}
	if vdc.DefaultFile = vs.Spec.Deployment.DefaultFile; vdc.DefaultFile == "" {
		vdc.DefaultFile = globalConf.DefaultDefaultFile
	}
	if vdc.Resources = vs.Spec.Deployment.VarnishResources; vdc.Resources == nil {
		vdc.Resources = &globalConf.DefaultVarnishResources
	}
	if vdc.ExporterResources = vs.Spec.Deployment.VarnishExporterResources; vdc.ExporterResources == nil {
		vdc.ExporterResources = &globalConf.DefaultVarnishExporterResources
	}
	if vdc.VarnishRestartPolicy = vs.Spec.Deployment.VarnishRestartPolicy; vdc.VarnishRestartPolicy == nil {
		vdc.VarnishRestartPolicy = &globalConf.DefaultVarnishRestartPolicy
	}
	if vdc.VolumeMountName = vs.Spec.Deployment.SharedVolume.Name; vdc.VolumeMountName == "" {
		vdc.VolumeMountName = globalConf.DefaultVolumeMountName
	}
	if vdc.VolumeMountPath = vs.Spec.Deployment.SharedVolume.Path; vdc.VolumeMountPath == "" {
		vdc.VolumeMountPath = globalConf.DefaultVolumeMountPath
	}
	if vdc.LivenessProbe = vs.Spec.Deployment.LivenessProbe; vdc.LivenessProbe == nil {
		vdc.LivenessProbe = globalConf.DefaultLivenessProbe
	}
	if vdc.ReadinessProbe = vs.Spec.Deployment.ReadinessProbe; vdc.ReadinessProbe == nil {
		vdc.ReadinessProbe = &globalConf.DefaultReadinessProbe
	}

	return &vdc, nil
}

func newVarnishDeployment(globalConf *config.Config, deploymentConf *varnishDeploymentConfig) (*appsv1.Deployment, error) {

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentConf.Name,
			Labels:    globalConf.VarnishCommonLabels,
			Namespace: deploymentConf.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deploymentConf.VarnishReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: deploymentConf.Labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deploymentConf.Labels,
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: deploymentConf.VolumeMountName,
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:  "varnish",
							Image: globalConf.VarnishImageFullPath,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: globalConf.VarnishPort,
								},
							},
							Env: []v1.EnvVar{
								{Name: "APP_SELECTOR_STRING", Value: labels.SelectorFromSet(deploymentConf.AppSelector).String()},
								{Name: "BACKENDS_FILE", Value: deploymentConf.BackendsFile},
								{Name: "DEFAULT_FILE", Value: deploymentConf.DefaultFile},
								{Name: "NAMESPACE", Value: deploymentConf.Namespace},
								{Name: "TARGET_PORT", Value: strconv.FormatInt(int64(deploymentConf.Port.Port), 10)},
								{Name: "VARNISH_PORT", Value: strconv.FormatInt(int64(globalConf.VarnishPort), 10)},
								{Name: "VARNISH_MEMORY", Value: deploymentConf.VarnishMemory},
								{Name: "VCL_DIR", Value: globalConf.VCLDir},
							},
							Resources: *deploymentConf.Resources,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      deploymentConf.VolumeMountName,
									MountPath: deploymentConf.VolumeMountPath,
								},
							},
							LivenessProbe:   deploymentConf.LivenessProbe,
							ReadinessProbe:  deploymentConf.ReadinessProbe,
							ImagePullPolicy: globalConf.VarnishImagePullPolicy,
						},
						{
							Name:  "varnish-exporter",
							Image: globalConf.VarnishImageFullPath,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: globalConf.VarnishExporterPort,
								},
							},
							Resources: *deploymentConf.ExporterResources,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      deploymentConf.VolumeMountName,
									MountPath: deploymentConf.VolumeMountPath,
								},
							},
						},
					},
					RestartPolicy:      *deploymentConf.VarnishRestartPolicy,
					ServiceAccountName: deploymentConf.ServiceAccountName,
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: deploymentConf.ImagePullSecretName,
						},
					},
					Affinity:    deploymentConf.Affinity,
					Tolerations: deploymentConf.Tolerations,
				},
			},
		},
	}

	return &deployment, nil
}
