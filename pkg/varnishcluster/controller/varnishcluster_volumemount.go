package controller

import (
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func haproxyConfigVolumeMount(readOnly bool) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.HaproxyConfigVolume,
		MountPath: vcapi.HaproxyConfigDir,
		ReadOnly:  readOnly,
	}
}

func haproxyScriptsVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.HaproxyScriptsVolume,
		MountPath: "/haproxy-scripts",
		ReadOnly:  true,
	}
}

func varnishSecretVolumeMount() v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.VarnishSecretVolume,
		MountPath: "/etc/varnish-secret",
		ReadOnly:  true,
	}
}

func varnishSettingsVolumeMount(readOnly bool) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.VarnishSettingsVolume,
		MountPath: "/etc/varnish",
		ReadOnly:  readOnly,
	}
}

func varnishSharedVolumeMount(readOnly bool) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      vcapi.VarnishSharedVolume,
		MountPath: "/var/lib/varnish",
		ReadOnly:  readOnly,
	}
}
