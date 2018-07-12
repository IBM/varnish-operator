package varnishservicehandler

import (
	"fmt"
	icmapiv1alpha1 "icm-varnish-k8s-operator/operator/controller/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/operator/controller/pkg/config"
	"reflect"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const (
	varnishServiceName = "varnish-exporter"
)

// VarnishServiceHandler describes the functions that handle events coming in for the VarnishService CRD
type VarnishServiceHandler struct {
	Conf *config.Config
	Client kubernetes.Interface
}

// ObjectAdded is called when a new instance of a VarnishService is detected, and so creates all of the necessary resources
func (h *VarnishServiceHandler) ObjectAdded(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}

	if err := applyHeadlessService(h.Client, h.Conf, vs); err != nil {
		return errors.Trace(err)
	}

	log.Infof("adding %+v", vs)
	return nil
}

// ObjectUpdated is called when a change to an existing VarnishService is detected, and applies any relevant changes
func (h *VarnishServiceHandler) ObjectUpdated(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}

	if err := applyHeadlessService(h.Client, h.Conf, vs); err != nil {
		return errors.Trace(err)
	}

	log.Infof("adding %+v", vs)
	return nil
}

// ObjectDeleted is called when an existing instance of VarnishService is deleted, so it cleans up all dependent resources
func (h *VarnishServiceHandler) ObjectDeleted(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}

	if err := deleteHeadlessService(h.Client, vs); err != nil {
		return errors.Trace(err)
	}
	
	log.Infof("adding %+v", vs)
	return nil
}

type varnishServiceConf struct {
	AppName      string
	AppLabels    map[string]string
	AppSelectors map[string]string
	Type         v1.ServiceType
}

func newService(globalConf *config.Config, serviceConf *varnishServiceConf) *v1.Service {
	varnishServiceName := fmt.Sprintf("%s-%s", serviceConf.AppName, varnishServiceName)
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   varnishServiceName,
			Labels: serviceConf.AppLabels,
			Annotations: map[string]string{
				"icm.ibm.com/owner":    globalConf.VarnishName,
				"prometheus.io/scrape": "true",
				"prometheus.io/port":   string(globalConf.VarnishExporterTargetPort),
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       globalConf.VarnishName,
					Port:       globalConf.VarnishPort,
					TargetPort: intstr.IntOrString{IntVal: globalConf.VarnishTargetPort},
					Protocol:   v1.ProtocolTCP,
				},
				{
					Name:       globalConf.VarnishExporterName,
					Port:       globalConf.VarnishExporterPort,
					TargetPort: intstr.IntOrString{IntVal: globalConf.VarnishExporterTargetPort},
					Protocol:   v1.ProtocolTCP,
				},
			},
			Selector: serviceConf.AppSelectors,
			Type:     serviceConf.Type,
		},
	}
}

type varnishDeploymentConfig struct {
	AppName             string
	AppSelectors        map[string]string
	AppLabels           map[string]string
	VarnishReplicas     int32
	VarnishMemory       int32
	BackendsFile        string
	DefaultFile         string
	Namespace           string
	Resources           *v1.ResourceRequirements
	ExporterResources   *v1.ResourceRequirements
	LimitResourceCPU    string
	LimitResourceMem    string
	RequestsResourceCPU string
	RequestsResourceMem string
	VolumeMountName     string
	VolumeMountPath     string
	LivenessProbe       *v1.Probe
	ReadinessProbe      *v1.Probe
	ImagePullPolicy     v1.PullPolicy
	ServiceAccountName  string
}

func newVarnishDeployment(globalConf *config.Config, deploymentConf *varnishDeploymentConfig) (*appsv1.Deployment, error) {
	varnishDeploymentName := fmt.Sprintf("%s-%s", deploymentConf.AppName, varnishServiceName)

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   varnishDeploymentName,
			Labels: deploymentConf.AppLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deploymentConf.VarnishReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: deploymentConf.AppSelectors,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deploymentConf.AppSelectors,
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
							Name:  varnishDeploymentName,
							Image: globalConf.VarnishImageFullPath,
							Ports: []v1.ContainerPort{
								{
									Name:          globalConf.VarnishName,
									HostPort:      globalConf.VarnishPort,
									ContainerPort: globalConf.VarnishPort,
								},
							},
							Env: []v1.EnvVar{
								{Name: "APP_NAME", Value: deploymentConf.AppName},
								{Name: "BACKENDS_FILE", Value: deploymentConf.BackendsFile},
								{Name: "DEFAULT_FILE", Value: deploymentConf.DefaultFile},
								{Name: "NAMESPACE", Value: deploymentConf.Namespace},
								{Name: "PORT_WATCH_NAME", Value: varnishDeploymentName},
								{Name: "VARNISH_PORT", Value: string(globalConf.VarnishPort)},
								{Name: "VARNISH_MEMORY", Value: fmt.Sprintf("%dM", deploymentConf.VarnishMemory)},
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
							ImagePullPolicy: deploymentConf.ImagePullPolicy,
						},
						{
							Name:  fmt.Sprintf("%s-exporter", varnishDeploymentName),
							Image: globalConf.VarnishImageFullPath,
							Ports: []v1.ContainerPort{
								{
									Name:          fmt.Sprintf("%s-exporter", varnishDeploymentName),
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
					RestartPolicy:      globalConf.RestartPolicy,
					ServiceAccountName: deploymentConf.ServiceAccountName,
					ImagePullSecrets: []v1.LocalObjectReference{
						{
							Name: globalConf.ImagePullSecret,
						},
					},
				},
			},
		},
	}

	return &deployment, nil
}
