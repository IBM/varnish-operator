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

	DefaultHaproxySidecar(in.HaproxySidecar)

	defaultVarnishZoneBalancingType(in.Backend.ZoneBalancing)
}

func DefaultHaproxySidecar(haproxySidecar *HaproxySidecar) {
	if haproxySidecar == nil {
		haproxySidecar = &HaproxySidecar{
			Enabled: false,
		}
	} else if haproxySidecar.Enabled {
		if haproxySidecar.MaxConnections == nil {
			haproxySidecar.MaxConnections = proto.Int32(64)
		}
		if haproxySidecar.ConnectTimeout == nil {
			haproxySidecar.ConnectTimeout = proto.Int32(5000)
		}
		if haproxySidecar.ClientTimeout == nil {
			haproxySidecar.ClientTimeout = proto.Int32(50000)
		}
		if haproxySidecar.ServerTimeout == nil {
			haproxySidecar.ServerTimeout = proto.Int32(50000)
		}
		if haproxySidecar.StatRefreshRate == nil {
			haproxySidecar.StatRefreshRate = proto.Int32(10)
		}
		if haproxySidecar.BackendAdditionalFlags == "" {
			haproxySidecar.BackendAdditionalFlags = "none"
		}
		if haproxySidecar.BackendServerMaxAgeHeader == nil {
			haproxySidecar.BackendServerMaxAgeHeader = proto.Int32(31536000)
		}
		if haproxySidecar.BackendServerPort == nil {
			haproxySidecar.BackendServerPort = proto.Int32(443)
		}
	}
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

	if len(in.ExtraInitContainers) > 0 {
		for i, container := range in.ExtraInitContainers {
			if container.ImagePullPolicy == "" {
				in.ExtraInitContainers[i].ImagePullPolicy = v1.PullIfNotPresent
			}
			if len(container.TerminationMessagePolicy) == 0 {
				in.ExtraInitContainers[i].TerminationMessagePolicy = v1.TerminationMessageReadFile
			}
			if len(container.TerminationMessagePath) == 0 {
				in.ExtraInitContainers[i].TerminationMessagePath = "/dev/termination-log"
			}
		}
	}

	if len(in.ExtraVolumeClaimTemplates) > 0 {
		for i, template := range in.ExtraVolumeClaimTemplates {
			if template.Spec.VolumeMode == nil {
				volumeMode := v1.PersistentVolumeFilesystem
				in.ExtraVolumeClaimTemplates[i].Spec.VolumeMode = &volumeMode
			}
		}
	}
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
