package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&VarnishCluster{}, func(obj interface{}) { SetVarnishClusterDefaults(obj.(*VarnishCluster)) })
	scheme.AddTypeDefaultingFunc(&VarnishClusterList{}, func(obj interface{}) { SetVarnishClusterListDefaults(obj.(*VarnishClusterList)) })
	return nil
}

func SetVarnishClusterDefaults(in *VarnishCluster) {
	defaultVarnishCluster(in)
}

func SetVarnishClusterListDefaults(in *VarnishClusterList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetVarnishClusterDefaults(a)
	}
}

func defaultVarnishCluster(in *VarnishCluster) {
	defaultVarnishClusterSpec(&in.Spec)

	if in.Spec.Varnish == nil {
		in.Spec.Varnish = &VarnishClusterVarnish{}
	}
	defaultVarnish(in.Spec.Varnish)
}

func defaultVarnishClusterSpec(in *VarnishClusterSpec) {
	var defaultReplicasNumber int32 = 1
	if in.Replicas == nil {
		in.Replicas = &defaultReplicasNumber
	}

	if in.LogLevel == "" {
		in.LogLevel = "info"
	}
	if in.LogFormat == "" {
		in.LogFormat = "json"
	}

	if in.UpdateStrategy != nil {
		if in.UpdateStrategy.Type == "" {
			in.UpdateStrategy.Type = OnDeleteVarnishClusterStrategyType
		}

		if in.UpdateStrategy.Type == VarnishUpdateStrategyDelayedRollingUpdate {
			if in.UpdateStrategy.DelayedRollingUpdate == nil {
				in.UpdateStrategy.DelayedRollingUpdate = &UpdateStrategyDelayedRollingUpdate{
					DelaySeconds: 60,
				}
			}
		}
	} else {
		in.UpdateStrategy = &VarnishClusterUpdateStrategy{
			Type: OnDeleteVarnishClusterStrategyType,
		}
	}

	if in.Service.MetricsPort == 0 {
		in.Service.MetricsPort = VarnishPrometheusExporterPort
	}

	if in.Service.Type == "" {
		in.Service.Type = v1.ServiceTypeClusterIP
	}

	if in.Backend.ZoneBalancing == nil {
		in.Backend.ZoneBalancing = &VarnishClusterBackendZoneBalancing{}
	}
	defaultVarnishZoneBalancingType(in.Backend.ZoneBalancing)
}

func defaultVarnish(in *VarnishClusterVarnish) {
	if in.ImagePullPolicy == "" {
		in.ImagePullPolicy = v1.PullAlways
	}

	if in.Resources == nil {
		in.Resources = &v1.ResourceRequirements{}
	}

	if in.Controller == nil {
		in.Controller = &VarnishClusterVarnishController{}
	}
	defaultVarnishController(in.Controller)

	if in.MetricsExporter == nil {
		in.MetricsExporter = &VarnishClusterVarnishMetricsExporter{}
	}
	defaultVarnishMetricsExporter(in.MetricsExporter)
}

func defaultVarnishController(in *VarnishClusterVarnishController) {
	if in.ImagePullPolicy == "" {
		in.ImagePullPolicy = v1.PullAlways
	}
}

func defaultVarnishMetricsExporter(in *VarnishClusterVarnishMetricsExporter) {
	if in.ImagePullPolicy == "" {
		in.ImagePullPolicy = v1.PullAlways
	}
}

func defaultVarnishZoneBalancingType(in *VarnishClusterBackendZoneBalancing) {
	if in.Type == "" {
		in.Type = VarnishClusterBackendZoneBalancingTypeDisabled
	}
}
