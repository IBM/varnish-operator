package controller

import (
	"context"
	"icm-varnish-k8s-operator/pkg/logger"
	"reflect"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

const (
	annotationConfigMapVersion    = "configMapVersion"
	annotationVCLVersion          = "VCLVersion"
	annotationActiveVCLConfigName = "activeVCLConfigName"
)

func (r *ReconcileVarnish) reconcilePod(ctx context.Context, filesChanged bool, pod *v1.Pod, cm *v1.ConfigMap) error {
	activeVCL, err := getActiveVCLConfig()
	if err != nil {
		return err
	}

	podCopy := &v1.Pod{}
	pod.DeepCopyInto(podCopy)

	if podCopy.Annotations == nil {
		podCopy.Annotations = make(map[string]string)
	}

	activeVCLConfigMapVersion := extractConfigMapVersion(activeVCL.Name)
	latestConfigMapInUse := cm.GetResourceVersion() == activeVCLConfigMapVersion

	// update version annotations only if the latest config map is in use or
	// if the config map changed but the files were not (e.g only labels were updated)
	if latestConfigMapInUse || (!latestConfigMapInUse && !filesChanged) {
		if cm.Annotations["VCLVersion"] == "" {
			//make sure we don't leave an outdated annotation if the new version of config map has no user version anymore
			delete(podCopy.Annotations, annotationVCLVersion)
		} else {
			podCopy.Annotations[annotationVCLVersion] = cm.Annotations["VCLVersion"]
		}

		podCopy.Annotations[annotationConfigMapVersion] = cm.GetResourceVersion()
	}

	podCopy.Annotations[annotationActiveVCLConfigName] = activeVCL.Name
	if !reflect.DeepEqual(pod.Annotations, podCopy.Annotations) {
		if err = r.Update(ctx, podCopy); err != nil {
			return errors.Wrapf(err, "failed to update pod")
		}

		logger.FromContext(ctx).Infow("Pod has been successfully updated")
	} else {
		logger.FromContext(ctx).Debugw("No updates for pod")
	}

	return nil
}
