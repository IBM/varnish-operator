package controller

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ibm/varnish-operator/pkg/logger"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

const (
	annotationConfigMapVersion     = "configMapVersion"
	annotationVCLVersion           = "VCLVersion"
	annotationActiveVCLConfigName  = "activeVCLConfigName"
	annotationLocalBackendsWeight  = "localBackendsWeight"
	annotationRemoteBackendsWeight = "remoteBackendsWeight"
)

func (r *ReconcileVarnish) reconcilePod(ctx context.Context, filesChanged bool, pod *v1.Pod, cm *v1.ConfigMap, localWeight float64, remoteWeight float64) error {
	activeVCLName, err := r.varnish.GetActiveConfigurationName()
	if err != nil {
		return err
	}

	podCopy := &v1.Pod{}
	pod.DeepCopyInto(podCopy)

	if podCopy.Annotations == nil {
		podCopy.Annotations = make(map[string]string)
	}

	activeVCLConfigMapVersion := extractConfigMapVersion(activeVCLName)
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

	podCopy.Annotations[annotationActiveVCLConfigName] = activeVCLName
	podCopy.Annotations[annotationLocalBackendsWeight] = fmt.Sprintf("%f", localWeight)
	podCopy.Annotations[annotationRemoteBackendsWeight] = fmt.Sprintf("%f", remoteWeight)
	logger.FromContext(ctx).Debugf("Local backends weight: %f", localWeight)
	logger.FromContext(ctx).Debugf("Remote backends weight: %f", remoteWeight)
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

// returns the config name the was used for VarnishClusterVCL config name creation
func extractConfigMapVersion(VCLConfigName string) string {
	parts := strings.Split(VCLConfigName, "-")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2]
}
