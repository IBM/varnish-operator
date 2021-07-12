package controller

import (
	"context"
	"strings"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"github.com/gogo/protobuf/proto"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishCluster) reconcileStatefulSet(ctx context.Context, instance, instanceStatus *vcapi.VarnishCluster, endpointSelector map[string]string) (*appsv1.StatefulSet, map[string]string, error) {
	varnishdArgs := getSanitizedVarnishArgs(&instance.Spec)
	varnishLabels := vclabels.CombinedComponentLabels(instance, vcapi.VarnishComponentVarnish)
	var varnishImage string
	if instance.Spec.Varnish.Image == "" {
		varnishImage = r.config.CoupledVarnishImage
	} else {
		varnishImage = instance.Spec.Varnish.Image
	}

	var updateStrategy appsv1.StatefulSetUpdateStrategy
	switch instance.Spec.UpdateStrategy.Type {
	case vcapi.OnDeleteVarnishClusterStrategyType,
		vcapi.DelayedRollingUpdateVarnishClusterStrategyType:
		updateStrategy.Type = appsv1.OnDeleteStatefulSetStrategyType
	case vcapi.RollingUpdateVarnishClusterStrategyType:
		updateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType
		updateStrategy.RollingUpdate = instance.Spec.UpdateStrategy.RollingUpdate
	}

	desired := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.StatefulSet(instance.Name),
			Labels:    varnishLabels,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: names.HeadlessService(instance.Name),
			Replicas:    instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: varnishLabels,
			},
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			UpdateStrategy:       updateStrategy,
			RevisionHistoryLimit: proto.Int(10),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: varnishLabels,
				},
				Spec: v1.PodSpec{
					// Share a single process namespace between all of the containers in a pod.
					// When this is set containers will be able to view and signal processes from other containers
					// in the same pod.It is required for the pod to provide reliable way to collect metrics.
					// Otherwise metrics collection container may only collect general varnish process metrics.
					ShareProcessNamespace:         proto.Bool(true),
					Volumes:                       getVarnishClusterVolumeMountsInstance().createVolumes(instance),
					Containers:                    getVarnishClusterContainersInstance().createContainers(instance, varnishdArgs, varnishImage, endpointSelector),
					TerminationGracePeriodSeconds: proto.Int64(30),
					DNSPolicy:                     v1.DNSClusterFirst,
					SecurityContext:               &v1.PodSecurityContext{},
					ServiceAccountName:            names.ServiceAccount(instance.Name),
					Affinity:                      instance.Spec.Affinity,
					Tolerations:                   instance.Spec.Tolerations,
					RestartPolicy:                 v1.RestartPolicyAlways,
				},
			},
		},
	}

	if instance.Spec.Varnish.ImagePullSecret != "" {
		desired.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			{
				Name: instance.Spec.Varnish.ImagePullSecret,
			},
		}
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentVarnish)
	logr = logr.With(logger.FieldComponentName, desired.Name)

	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return nil, nil, errors.Wrap(err, "could not set controller as the OwnerReference for statefulset")
	}
	r.scheme.Default(desired)

	found := &appsv1.StatefulSet{}

	err := r.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, found)
	// If the statefulset does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the statefulset exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating StatefulSet", "new", desired)
		err = r.Create(ctx, desired)
		if err != nil {
			return nil, nil, errors.Wrap(err, "could not create statefulset")
		}
	} else if err != nil {
		return nil, nil, errors.Wrap(err, "could not get current state of statefulset")
	} else {
		// the pod selector is immutable once set, so always enforce the same as existing
		desired.Spec.Selector = found.Spec.Selector
		desired.Spec.Template.Labels = found.Spec.Template.Labels
		desired.Spec.Template.Annotations = found.Spec.Template.Annotations
		if !compare.EqualStatefulSet(found, desired) {
			logr.Infoc("Updating StatefulSet", "diff", compare.DiffStatefulSet(found, desired))
			found.Spec = desired.Spec
			found.Labels = desired.Labels
			if err = r.Update(ctx, found); err != nil {
				return nil, nil, errors.Wrap(err, "could not update statefulset")
			}
		} else {
			logr.Debugw("No updates for StatefulSet")
		}
	}

	instanceStatus.Status.VarnishArgs = strings.Join(varnishdArgs, " ")
	instanceStatus.Status.Replicas = found.Status.Replicas

	return found, varnishLabels, nil
}

func imageNameGenerate(specified, base, suffix string) string {
	if specified != "" {
		return specified
	}
	baseName, tag := splitNameAndTag(base)
	return baseName + suffix + tag
}

func splitNameAndTag(fullName string) (image, tag string) {
	parts := strings.Split(fullName, ":")
	image = parts[0]
	if len(parts) > 1 {
		tag = ":" + parts[1]
	}
	return
}
