package compare

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	deployIgnoreFields = cmpopts.IgnoreFields(appsv1.Deployment{}, "Spec.Template.Spec.DeprecatedServiceAccount")
	compareQuantity    = cmp.Comparer(func(x, y resource.Quantity) bool { return x.Cmp(y) == 0 })
	deployOpts         = []cmp.Option{cmpopts.IgnoreFields(appsv1.Deployment{}, sharedIgnoreFields...), deployIgnoreFields, compareQuantity}
)

var (
	deployContainerPortDefaults = &v1.ContainerPort{
		Protocol: v1.ProtocolTCP,
	}

	readinessProbeDefaults = v1.Probe{
		TimeoutSeconds:   int32(1),
		PeriodSeconds:    int32(10),
		SuccessThreshold: int32(1),
		FailureThreshold: int32(3),
	}

	deployContainerDefaults = &v1.Container{
		TerminationMessagePath:   v1.TerminationMessagePathDefault,
		TerminationMessagePolicy: v1.TerminationMessageReadFile,
		ImagePullPolicy:          v1.PullIfNotPresent,
	}

	thirty         = int64(30)
	ten            = int32(10)
	sixHundred     = int32(600)
	deployDefaults = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &thirty,
					DNSPolicy:                     v1.DNSClusterFirst,
					SecurityContext:               &v1.PodSecurityContext{},
					SchedulerName:                 v1.DefaultSchedulerName,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			RevisionHistoryLimit:    &ten,
			ProgressDeadlineSeconds: &sixHundred,
		},
	}
	twentyFivePercent = intstr.FromString("25%")
)

func withDeploymentDefaults(d *appsv1.Deployment) *appsv1.Deployment {
	var dd appsv1.Deployment
	d.DeepCopyInto(&dd)

	mergo.Merge(&dd, deployDefaults)
	for c := range dd.Spec.Template.Spec.Containers {
		container := &dd.Spec.Template.Spec.Containers[c]
		mergo.Merge(container, deployContainerDefaults)
		if container.ReadinessProbe != nil {
			mergo.Merge(container.ReadinessProbe, readinessProbeDefaults)
		}
		for p := range dd.Spec.Template.Spec.Containers[c].Ports {
			mergo.Merge(&dd.Spec.Template.Spec.Containers[c].Ports[p], deployContainerPortDefaults)
		}
	}

	if dd.Spec.Strategy.RollingUpdate == nil && dd.Spec.Strategy.Type == appsv1.RollingUpdateDeploymentStrategyType {
		dd.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &twentyFivePercent,
			MaxSurge:       &twentyFivePercent,
		}
	}
	return &dd
}

func withDeploymentInheritance(desired, found *appsv1.Deployment) {
	if desired.Annotations == nil {
		desired.Annotations = make(map[string]string)
	}
	if desired.Annotations["deployment.kubernetes.io/revision"] == "" {
		desired.Annotations["deployment.kubernetes.io/revision"] = found.Annotations["deployment.kubernetes.io/revision"]
	}
}

// EqualDeployment compares 2 deployments for equality
func EqualDeployment(desired, found *appsv1.Deployment) bool {
	desiredCopy := withDeploymentDefaults(desired)
	withDeploymentInheritance(desiredCopy, found)
	return cmp.Equal(desiredCopy, found, deployOpts...)
}

// DiffDeployment generates a patch diff between 2 deployments
func DiffDeployment(desired, found *appsv1.Deployment) string {
	desiredCopy := withDeploymentDefaults(desired)
	withDeploymentInheritance(desiredCopy, found)
	return cmp.Diff(desiredCopy, found, deployOpts...)
}
