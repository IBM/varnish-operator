package controller

import (
	"context"
	"fmt"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	annotationVCLVersion           = "VCLVersion"
	annotationHaproxyConfigVersion = "HaproxyConfigVersion"
)

func (r *ReconcileVarnishCluster) reconcileConfigMap(ctx context.Context, podsSelector map[string]string, instance, instanceStatus *vcapi.VarnishCluster) error {
	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentVCLFileConfigMap)
	logr = logr.With(logger.FieldComponentName, instance.Spec.VCL.ConfigMapName)

	cm := &v1.ConfigMap{}
	cmLabels := vclabels.CombinedComponentLabels(instance, vcapi.VarnishComponentVCLFileConfigMap)
	err := r.Get(ctx, types.NamespacedName{Name: *instance.Spec.VCL.ConfigMapName, Namespace: instance.Namespace}, cm)
	// if the ConfigMap does not exist, create it and set it with the default VCL files
	// Else if there was a problem doing the Get, just return an error
	// Else fill in missing values -- "OwnerReference" or Labels
	// Else do nothing
	if err != nil && kerrors.IsNotFound(err) {
		cm = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      *instance.Spec.VCL.ConfigMapName,
				Labels:    cmLabels,
				Namespace: instance.Namespace,
			},
			Data: map[string]string{
				*instance.Spec.VCL.EntrypointFileName: entrypointVCLFileContent,
				"backends.vcl.tmpl":                   backendsVCLTmplFileContent,
			},
		}
		if err := controllerutil.SetControllerReference(instance, cm, r.scheme); err != nil {
			return errors.Wrap(err, "could not initialize default ConfigMap")
		}

		logr.Infoc("Creating ConfigMap with default VCL files", "new", cm)
		if err = r.Create(ctx, cm); err != nil {
			return errors.Wrap(err, "could not create ConfigMap")
		}
	} else if err != nil {
		return errors.Wrap(err, "could not get current state of ConfigMap")
	} else {
		cmCopy := cm.DeepCopy() //create a copy to check later if the config map changed and needs to be updated
		// TODO: there may be a problem if the configmap is already owned by something else. That will prevent the `Watch` fn (in varnishcluster_controller.go#run) from detecting updates to the ConfigMap. It will also cause this code to throw an unhandled error that we may want to handle
		if err = controllerutil.SetControllerReference(instance, cm, r.scheme); err != nil {
			return errors.Wrap(err, "could not set controller as the OwnerReference for existing ConfigMap")
		}
		// don't trample on any labels created by user
		if cm.Labels == nil {
			cm.Labels = make(map[string]string, len(cmLabels))
		}
		for l, v := range cmLabels {
			cm.Labels[l] = v
		}

		if !compare.EqualConfigMap(cm, cmCopy) {
			logr.Infow("Updating ConfigMap with defaults", "diff", compare.DiffConfigMap(cm, cmCopy))
			if err = r.Update(ctx, cm); err != nil {
				return errors.Wrap(err, "could not update configmap")
			}
		} else {
			logr.Debugw("No updates for ConfigMap")
		}
	}

	instanceStatus.Status.VCL.ConfigMapVersion = cm.GetResourceVersion()
	if cm.Annotations != nil && cm.Annotations[annotationVCLVersion] != "" {
		v := cm.Annotations[annotationVCLVersion]
		instanceStatus.Status.VCL.Version = &v
	} else {
		instanceStatus.Status.VCL.Version = nil //ensure the status field is empty if the annotation is
	}

	availabilityString, err := r.availabilityString(podsSelector, "configMapVersion", instance.Status.VCL.ConfigMapVersion, logr)
	if err != nil {
		return err
	}
	instanceStatus.Status.VCL.Availability = availabilityString
	return nil
}

func (r *ReconcileVarnishCluster) availabilityString(podsSelector map[string]string, annotationKey string, cmVersion string, logr *logger.Logger) (string, error) {
	pods := &v1.PodList{}
	selector := labels.SelectorFromSet(podsSelector)
	if err := r.List(context.Background(), pods, &client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return "", errors.Wrap(err, "can't get list of pods")
	}
	latest, outdated := 0, 0
	for _, item := range pods.Items {
		//do not count pods that are not updated with VCL version. Those are pods that are just created and not fully functional
		if item.Annotations[annotationKey] == "" {
			logr.Debugf("ConfigMapVersion annotation (%s) is not present. Skipping the pod - %s.", annotationKey, item.Name)
		} else if item.Annotations[annotationKey] == cmVersion {
			latest++
		} else {
			outdated++
		}
	}
	return fmt.Sprintf("%d latest / %d outdated", latest, outdated), nil
}
