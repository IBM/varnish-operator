package controller

import (
	"context"
	"reflect"

	"github.com/pkg/errors"

	"k8s.io/api/core/v1"
)

const (
	annotationConfigMapVersion    = "configMapVersion"
	annotationVCLVersion          = "VCLVersion"
	annotationActiveVCLConfigName = "activeVCLConfigName"
)

func (r *ReconcileVarnish) reconcilePod(filesChanged bool, pod *v1.Pod, cm *v1.ConfigMap) error {
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
		if err = r.Update(context.Background(), podCopy); err != nil {
			return errors.Wrapf(err, "failed to update pod")
		}

		r.logger.Infow("Pod has been successfully updated")
	} else {
		r.logger.Debugw("No updates for pod")
	}

	return nil
}
