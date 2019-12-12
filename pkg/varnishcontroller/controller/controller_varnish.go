package controller

import (
	"context"
	"fmt"
	"icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/events"
	"strings"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

const (
	// For VCL version name we use config map resource version which is a number.
	// Varnish doesn't accept config name that have numbers in the beginning. Even if it is disguised as strings (e.g. "1243").
	// For that reasons we prepend this prefix.
	VCLVersionPrefix = "v"
)

func (r *ReconcileVarnish) reconcileVarnish(ctx context.Context, vc *v1alpha1.VarnishCluster, pod *v1.Pod, cm *v1.ConfigMap) error {
	logr := logger.FromContext(ctx)
	logr.Debugw("Starting varnish reload...")
	start := time.Now()
	out, err := r.varnish.Reload(createVCLConfigName(cm.GetResourceVersion()), *vc.Spec.VCL.EntrypointFileName)
	if err != nil {
		if strings.Contains(string(out), "VarnishClusterVCL compilation failed") {
			vcEventMsg := "VarnishClusterVCL compilation failed for pod " + pod.Name + ". See pod logs for details"
			podEventMsg := "VarnishClusterVCL compilation failed. See logs for details"
			r.eventHandler.Warning(pod, events.EventReasonVCLCompilationError, podEventMsg)
			r.eventHandler.Warning(vc, events.EventReasonVCLCompilationError, vcEventMsg)
			logr.Warnw(string(out))
			return nil
		}

		podEventMsg := "VarnishClusterVarnish reload failed for pod " + pod.Name + ". See pod logs for details"
		vcEventMsg := "VarnishClusterVarnish reload failed. See logs for details"
		r.eventHandler.Warning(pod, events.EventReasonReloadError, podEventMsg)
		r.eventHandler.Warning(vc, events.EventReasonReloadError, vcEventMsg)
		return errors.Wrap(err, string(out))
	}
	logr.Debugf("VarnishClusterVarnish successfully reloaded in %f seconds", time.Since(start).Seconds())
	return nil
}

// creates the VarnishClusterVCL config name from config map version
func createVCLConfigName(configMapVersion string) string {
	return fmt.Sprintf("%s-%s-%d", VCLVersionPrefix, configMapVersion, time.Now().Unix())
}
