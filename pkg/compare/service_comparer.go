package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/imdario/mergo"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	serviceIgnoreFields = cmpopts.IgnoreFields(v1.Service{}, "Spec.ClusterIP")
	serviceOpts         = []cmp.Option{cmpopts.IgnoreFields(v1.Service{}, sharedIgnoreMetadata...), cmpopts.IgnoreFields(v1.Service{}, sharedIgnoreStatus...), serviceIgnoreFields}
)

var (
	serviceDefaults = &v1.Service{
		Spec: v1.ServiceSpec{
			SessionAffinity: v1.ServiceAffinityNone,
			Type:            v1.ServiceTypeClusterIP,
		},
	}
)

func withServiceDefaults(s *v1.Service) *v1.Service {
	var ss v1.Service
	s.DeepCopyInto(&ss)

	mergo.Merge(&ss, serviceDefaults)

	for i := range ss.Spec.Ports {
		if ss.Spec.Ports[i].TargetPort == (intstr.IntOrString{}) {
			ss.Spec.Ports[i].TargetPort = intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: ss.Spec.Ports[i].Port,
			}
		}
	}

	return &ss
}

// EqualService compares 2 services for equality
func EqualService(found, desired *v1.Service) bool {
	desiredCopy := withServiceDefaults(desired)
	return cmp.Equal(found, desiredCopy, serviceOpts...)
}

// DiffService generates a patch diff between 2 services
func DiffService(found, desired *v1.Service) string {
	desiredCopy := withServiceDefaults(desired)
	return cmp.Diff(found, desiredCopy, serviceOpts...)
}
