package controller

import (
	"sync"

	"github.com/gogo/protobuf/proto"
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

type VarnishClusterVolumes struct {
}

var varnishClusterVolumesLock = &sync.Mutex{}
var varnishClusterVolumesSingleton *VarnishClusterVolumes

func getVarnishClusterVolumeMountsInstance() *VarnishClusterVolumes {
	if varnishClusterVolumesSingleton == nil {
		varnishClusterVolumesLock.Lock()
		defer varnishClusterVolumesLock.Unlock()
		if varnishClusterVolumesSingleton == nil {
			varnishClusterVolumesSingleton = &VarnishClusterVolumes{}
		}
	}
	return varnishClusterVolumesSingleton
}

func (r *VarnishClusterVolumes) createHaproxyConfigVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.HaproxyConfigVolume,
		MountPath: "/usr/local/etc/haproxy",
		ReadOnly:  true,
	}
}

func (r *VarnishClusterVolumes) createVarnishSharedVolumeMount(readOnly bool) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.VarnishSharedVolume,
		MountPath: "/var/lib/varnish",
		ReadOnly:  readOnly,
	}
}

func (r *VarnishClusterVolumes) createVarnishSettingsVolumeMount(readOnly bool) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.VarnishSettingsVolume,
		MountPath: "/etc/varnish",
		ReadOnly:  readOnly,
	}
}

func (r *VarnishClusterVolumes) createVarnishSecretVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.VarnishSecretVolume,
		MountPath: "/etc/varnish-secret",
		ReadOnly:  true,
	}
}

func (r *VarnishClusterVolumes) createHaproxyScriptsVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.HaproxyScriptsVolume,
		MountPath: "/haproxy-scripts",
		ReadOnly:  true,
	}
}

func (r *VarnishClusterVolumes) createVolumes(instance *vcapi.VarnishCluster) []v1.Volume {
	varnishSecretName, varnishSecretKeyName := namesForInstanceSecret(instance)
	volumes := []v1.Volume{
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
	}
	if instance.Spec.HaproxySidecar.Enabled {
		haproxyVolume := v1.Volume{
			Name: vcapi.HaproxyConfigVolume,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: instance.Spec.HaproxySidecar.ConfigMapName,
					},
				},
			},
		}
		volumes = append(volumes, haproxyVolume)

		haproxyScriptsVolume := v1.Volume{
			Name: vcapi.HaproxyScriptsVolume,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: vcapi.HaproxyScriptsVolume + "-configmap",
					},
					DefaultMode: proto.Int32(0777),
				},
			},
		}
		volumes = append(volumes, haproxyScriptsVolume)
	}
	return volumes
}
