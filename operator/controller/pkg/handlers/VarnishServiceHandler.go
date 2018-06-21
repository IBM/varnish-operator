package handlers

import (
	"fmt"
	icmapiv1alpha1 "icm-varnish-k8s-operator/operator/controller/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/operator/controller/pkg/config"
	"reflect"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	varnishServiceName = "varnish-exporter"
)

// VarnishServiceHandler describes the functions that handle events coming in for the VarnishService CRD
type VarnishServiceHandler struct {
	Conf *config.Config
}

// ObjectAdded prints out the VarnishService
func (h *VarnishServiceHandler) ObjectAdded(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}
	log.Infof("adding %+v", vs)
	return nil
}

// ObjectUpdated prints out the VarnishService
func (h *VarnishServiceHandler) ObjectUpdated(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}
	log.Infof("adding %+v", vs)
	return nil
}

// ObjectDeleted prints out the VarnishService
func (h *VarnishServiceHandler) ObjectDeleted(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}
	log.Infof("adding %+v", vs)
	return nil
}

type headlessConfig struct {
	ServiceName       string
	AppLabels         map[string]string
	AppSelectors      map[string]string
	VarnishBackedPort apiv1.ServicePort
	OtherPorts        []apiv1.ServicePort
}

func newHeadlessService(globalConf config.Config, headlessConf headlessConfig) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   headlessConf.ServiceName,
			Labels: headlessConf.AppLabels,
			Annotations: map[string]string{
				"icm.ibm.com/owner":               varnishServiceName,
				"icm.ibm.com/varnish-backed-port": string(headlessConf.VarnishBackedPort.Port),
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports:     append(headlessConf.OtherPorts, headlessConf.VarnishBackedPort),
			Selector:  headlessConf.AppSelectors,
			ClusterIP: "None",
			Type:      apiv1.ServiceTypeClusterIP,
		},
	}
}

type varnishDeploymentConfig struct {
	AppName         string
	AppSelectors    map[string]string
	AppLabels       map[string]string
	VarnishReplicas int32
	VarnishMemory   int32
	BackendsFile    string
	DefaultFile     string
	Namespace       string
}

func newVarnishDeployment(globalConf config.Config, deploymentConf varnishDeploymentConfig) *appsv1.Deployment {
	replicas := deploymentConf.VarnishReplicas
	varnishDeploymentName := fmt.Sprintf("%s-%s", deploymentConf.AppName, varnishServiceName)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   varnishDeploymentName,
			Labels: deploymentConf.AppLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: deploymentConf.AppSelectors,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deploymentConf.AppSelectors,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            deploymentConf.AppName,
							Image:           globalConf.FullImagePath(),
							ImagePullPolicy: apiv1.PullAlways,
							Ports: []apiv1.ContainerPort{
								{
									Name:          globalConf.VarnishName,
									HostPort:      globalConf.VarnishPort,
									ContainerPort: globalConf.VarnishPort,
								},
							},
							Env: []apiv1.EnvVar{
								{Name: "APP_NAME", Value: deploymentConf.AppName},
								{Name: "BACKENDS_FILE", Value: deploymentConf.BackendsFile},
								{Name: "DEFAULT_FILE", Value: deploymentConf.DefaultFile},
								{Name: "NAMESPACE", Value: deploymentConf.Namespace},
								{Name: "PORT_WATCH_NAME", Value: varnishDeploymentName},
								{Name: "VARNISH_PORT", Value: string(globalConf.VarnishPort)},
								{Name: "VARNISH_MEMORY", Value: fmt.Sprintf("%dM", deploymentConf.VarnishMemory)},
								{Name: "VCL_DIR", Value: globalConf.VCLDir},
							},
							// TODO: continue converting yaml file to golang
							Resources: apiv1.ResourceRequirements{
								Limits:   apiv1.ResourceList{},
								Requests: apiv1.ResourceList{},
							},
						},
					},
				},
			},
		},
	}
}
