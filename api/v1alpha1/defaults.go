package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	v1 "k8s.io/api/core/v1"
)

func RegisterDefaults(scheme *runtime.Scheme) error {
	scheme.AddTypeDefaultingFunc(&VarnishService{}, func(obj interface{}) { SetVarnishServiceDefaults(obj.(*VarnishService)) })
	scheme.AddTypeDefaultingFunc(&VarnishServiceList{}, func(obj interface{}) { SetVarnishServiceListDefaults(obj.(*VarnishServiceList)) })
	return nil
}

func SetVarnishServiceDefaults(in *VarnishService) {
	defaultVarnishService(in)
}

func SetVarnishServiceListDefaults(in *VarnishServiceList) {
	for i := range in.Items {
		a := &in.Items[i]
		SetVarnishServiceDefaults(a)
	}
}

func defaultVarnishService(in *VarnishService) {
	defaultVarnishServiceSpec(&in.Spec)
	defaultVarnishContainer(&in.Spec.StatefulSet.Container)
	defaultVarnishServiceService(&in.Spec.Service)
}

func defaultVarnishServiceSpec(in *VarnishServiceSpec) {
	if in.LogLevel == "" {
		in.LogLevel = "info"
	}
	if in.LogFormat == "" {
		in.LogFormat = "json"
	}

	if in.StatefulSet.UpdateStrategy.Type == "" {
		in.StatefulSet.UpdateStrategy.Type = appsv1.OnDeleteStatefulSetStrategyType
	}

	if in.StatefulSet.UpdateStrategy.Type == VarnishUpdateStrategyDelayedRollingUpdate {
		if in.StatefulSet.UpdateStrategy.DelayedRollingUpdate == nil {
			in.StatefulSet.UpdateStrategy.DelayedRollingUpdate = &UpdateStrategyDelayedRollingUpdate{
				DelaySeconds: 60,
			}
		}
	}
}

func defaultVarnishContainer(in *VarnishContainer) {
	if in.ImagePullPolicy == "" {
		in.ImagePullPolicy = v1.PullAlways
	}
	if in.RestartPolicy == "" {
		in.RestartPolicy = v1.RestartPolicyAlways
	}
}

func defaultVarnishServiceService(in *VarnishServiceService) {
	s := &v1.Service{Spec: in.ServiceSpec}
	s.Spec.DeepCopyInto(&in.ServiceSpec)

	if in.VarnishPort.Name == "" {
		in.VarnishPort.Name = "varnish"
	}
	if in.VarnishPort.TargetPort == (intstr.IntOrString{}) {
		in.VarnishPort.TargetPort = intstr.IntOrString{
			Type:   intstr.Int,
			IntVal: in.VarnishPort.Port,
		}
	}

	if len(in.VarnishPort.Protocol) == 0 {
		in.VarnishPort.Protocol = v1.ProtocolTCP
	}

	if in.VarnishExporterPort.Name == "" {
		in.VarnishExporterPort.Name = "varnishexporter"
	}

	if in.VarnishExporterPort.Port == 0 {
		in.VarnishExporterPort.Port = VarnishPrometheusExporterPort
	}

	if len(in.VarnishExporterPort.Protocol) == 0 {
		in.VarnishExporterPort.Protocol = v1.ProtocolTCP
	}

	if len(in.SessionAffinity) == 0 {
		in.SessionAffinity = v1.ServiceAffinityNone
	}

	if len(in.Type) == 0 {
		in.Type = v1.ServiceTypeClusterIP
	}

	//we don't support running exporter on custom port yet so ignore the value set by user for now
	in.VarnishExporterPort.TargetPort = intstr.FromInt(VarnishPrometheusExporterPort)
}
