package varnishservice

import (
	"context"
	"errors"
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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileDeployment(instance, instanceStatus *icmapiv1alpha1.VarnishService, serviceAccountName string, applicationPort *v1.ServicePort, endpointSelector map[string]string) (map[string]string, error) {
	deployConfig, err := newVarnishDeploymentConfig(r.globalConf, instance, serviceAccountName, applicationPort, endpointSelector)
	if err != nil {
		return nil, logger.RError(err, "could not generate deployment config")
	}
	desired, err := newVarnishDeployment(r.globalConf, deployConfig)
	if err != nil {
		return nil, logger.RError(err, "could not generate deployment")
	}
	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return nil, logger.RError(err, "could not set controller as the OwnerReference for deployment", "name", desired.Name, "namespace", desired.Namespace)
	}

	found := &appsv1.Deployment{}

	err = r.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, found)
	// If the deployment does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the deployment exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logger.Info("Creating Deployment", "config", desired)
		err = r.Create(context.TODO(), desired)
		if err != nil {
			return nil, logger.RError(err, "could not create deployment", "name", desired.Name, "namespace", desired.Namespace)
		}
	} else if err != nil {
		return nil, logger.RError(err, "could not get current state of deployment", "name", desired.Name, "namespace", desired.Namespace)
	} else if !compare.EqualDeployment(desired, found) {
		logger.Info("Updating Deployment", "diff", compare.DiffDeployment(desired, found))
		found.Spec = desired.Spec
		err = r.Update(context.TODO(), found)
		if err != nil {
			return nil, logger.RError(err, "could not update deployment", "name", desired.Name, "namespace", desired.Namespace)
		}
	} else {
		logger.V5Info("No updates for Deployment")
	}

	instanceStatus.Status.Deployment = found.Status

	return deployConfig.Labels, nil
}

type varnishDeploymentConfig struct {
	Name                 string
	Namespace            string
	Labels               map[string]string
	EndpointSelector     map[string]string
	VarnishRestartPolicy *v1.RestartPolicy
	VarnishReplicas      int32
	VarnishMemory        string
	BackendsFile         string
	DefaultFile          string
	Resources            *v1.ResourceRequirements
	ExporterResources    *v1.ResourceRequirements
	SharedVolumeName     string
	SharedVolumePath     string
	LivenessProbe        *v1.Probe
	ReadinessProbe       *v1.Probe
	ServiceAccountName   string
	Port                 v1.ServicePort
	ImagePullSecretName  string
	Affinity             *v1.Affinity
	Tolerations          []v1.Toleration
}

func newVarnishDeploymentConfig(globalConf *config.Config, vs *icmapiv1alpha1.VarnishService, serviceAccountName string, applicationPort *v1.ServicePort, endpointSelector map[string]string) (*varnishDeploymentConfig, error) {
	vdc := varnishDeploymentConfig{
		EndpointSelector:   endpointSelector,
		ServiceAccountName: serviceAccountName,
		Port:               *applicationPort,
		Affinity:           vs.Spec.Deployment.Affinity,
		Tolerations:        vs.Spec.Deployment.Tolerations,
	}
	// required fields
	if vdc.Name = vs.Name + "-deployment"; vdc.Name == "-deployment" {
		return &vdc, errors.New("name not defined")
	}

	vdc.Labels = map[string]string{"varnish-component": vs.Name + "-varnish"}

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
	if vdc.SharedVolumeName = vs.Spec.Deployment.SharedVolume.Name; vdc.SharedVolumeName == "" {
		vdc.SharedVolumeName = globalConf.DefaultSharedVolumeName
	}
	if vdc.SharedVolumePath = vs.Spec.Deployment.SharedVolume.Path; vdc.SharedVolumePath == "" {
		vdc.SharedVolumePath = globalConf.DefaultSharedVolumePath
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
							Name: deploymentConf.SharedVolumeName,
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
								{Name: "ENDPOINT_SELECTOR_STRING", Value: labels.SelectorFromSet(deploymentConf.EndpointSelector).String()},
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
									Name:      deploymentConf.SharedVolumeName,
									MountPath: deploymentConf.SharedVolumePath,
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
									Name:      deploymentConf.SharedVolumeName,
									MountPath: deploymentConf.SharedVolumePath,
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
