package controller

import (
	"sync"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type VarnishClusterContainers struct {
	vols *VarnishClusterVolumes
}

var varnishClusterContainersLock = &sync.Mutex{}
var varnishClusterContainersSingleton *VarnishClusterContainers

func getVarnishClusterContainersInstance() *VarnishClusterContainers {
	if varnishClusterContainersSingleton == nil {
		varnishClusterContainersLock.Lock()
		defer varnishClusterContainersLock.Unlock()
		if varnishClusterContainersSingleton == nil {
			varnishClusterContainersSingleton = &VarnishClusterContainers{
				vols: getVarnishClusterVolumeMountsInstance(),
			}
		}
	}
	return varnishClusterContainersSingleton
}

func (r *VarnishClusterContainers) createVarnishContainer(instance *vcapi.VarnishCluster, varnishdArgs []string, varnishImage string) v1.Container {
	//Varnish container
	return v1.Container{
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
			r.vols.createVarnishSharedVolumeMount(false),
			r.vols.createVarnishSettingsVolumeMount(true),
			r.vols.createVarnishSecretVolumeMount(),
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
	}
}

func (r *VarnishClusterContainers) createVarnishMetricsContainer(instance *vcapi.VarnishCluster, varnishImage string) v1.Container {
	varnishMetricsImage := imageNameGenerate(instance.Spec.Varnish.MetricsExporter.Image, varnishImage, vcapi.VarnishMetricsExporterImage)

	//Varnish metrics
	return v1.Container{
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
			r.vols.createVarnishSharedVolumeMount(true),
			r.vols.createVarnishSettingsVolumeMount(true),
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
	}
}

func (r *VarnishClusterContainers) createVarnishControllerContainer(instance *vcapi.VarnishCluster, varnishImage string, endpointSelector map[string]string) v1.Container {
	gvk := instance.GroupVersionKind()
	varnishControllerImage := imageNameGenerate(instance.Spec.Varnish.Controller.Image, varnishImage, vcapi.VarnishControllerImage)

	//Varnish controller
	return v1.Container{
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
			r.vols.createVarnishSharedVolumeMount(false),
			r.vols.createVarnishSettingsVolumeMount(false),
			r.vols.createVarnishSecretVolumeMount(),
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
	}
}

func (r *VarnishClusterContainers) createHaproxySidecarContainer(instance *vcapi.VarnishCluster) v1.Container {
	//haproxy sidecar
	return v1.Container{
		Name:            vcapi.HaproxyContainerName,
		Image:           instance.Spec.HaproxySidecar.Image,
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
				Name:      vcapi.HaproxyConfigVolume,
				MountPath: vcapi.HaproxyConfigMountPath,
				ReadOnly:  true,
			},
		},
	}
}

func (r *VarnishClusterContainers) createContainers(instance *vcapi.VarnishCluster, varnishdArgs []string, varnishImage string, endpointSelector map[string]string) []v1.Container {
	containers := []v1.Container{
		r.createVarnishContainer(instance, varnishdArgs, varnishImage),
		r.createVarnishMetricsContainer(instance, varnishImage),
		r.createVarnishControllerContainer(instance, varnishImage, endpointSelector),
	}
	if instance.Spec.HaproxySidecar != nil && instance.Spec.HaproxySidecar.Enabled {
		containers = append(containers, r.createHaproxySidecarContainer(instance))
	}
	return containers
}
