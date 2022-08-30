package controller

import (
	"github.com/gogo/protobuf/proto"
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func varnishClusterVolumes(instance *vcapi.VarnishCluster) []v1.Volume {
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
	if instance.Spec.HaproxySidecar != nil && instance.Spec.HaproxySidecar.Enabled {
		haproxyVolumes := []v1.Volume{
			{
				Name: vcapi.HaproxyConfigVolume,
				VolumeSource: v1.VolumeSource{
					EmptyDir: &v1.EmptyDirVolumeSource{},
				},
			},
			{
				Name: vcapi.HaproxyScriptsVolume,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: vcapi.HaproxyScriptsVolume + "-configmap",
						},
						DefaultMode: proto.Int32(0777),
					},
				},
			},
		}
		volumes = append(volumes, haproxyVolumes...)
	}
	return append(volumes, instance.Spec.Varnish.ExtraVolumes...)
}
