package v1alpha1

import (
	"time"

	"github.com/gogo/protobuf/proto"

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

	if in.Monitoring == nil {
		in.Monitoring = &VarnishClusterMonitoring{}
	}

	if in.Monitoring.PrometheusServiceMonitor == nil {
		in.Monitoring.PrometheusServiceMonitor = &VarnishClusterMonitoringPrometheusServiceMonitor{}
	}

	// set default if empty or invalid value get here
	if _, err := time.ParseDuration(in.Monitoring.PrometheusServiceMonitor.ScrapeInterval); err != nil {
		in.Monitoring.PrometheusServiceMonitor.ScrapeInterval = "30s"
	}

	if in.Service.MetricsPort == nil {
		in.Service.MetricsPort = proto.Int32(VarnishPrometheusExporterPort)
	}

	if in.Service.Type == "" {
		in.Service.Type = v1.ServiceTypeClusterIP
	}

	if in.Backend.ZoneBalancing == nil {
		in.Backend.ZoneBalancing = &VarnishClusterBackendZoneBalancing{}
	}

	if in.HaproxySidecar == nil {
		in.HaproxySidecar = &HaproxySidecar{
			Enabled: false,
		}
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
