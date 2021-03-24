package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/events"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/varnishadm"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

const (
	// For VCL version name we use config map resource version which is a number.
	// Varnish doesn't accept config name that have numbers in the beginning. Even if it is disguised as strings (e.g. "1243").
	// For that reasons we prepend this prefix.
	VCLVersionPrefix = "v-"
)

func (r *ReconcileVarnish) reconcileVarnish(ctx context.Context, vc *v1alpha1.VarnishCluster, pod *v1.Pod, cm *v1.ConfigMap) error {
	logr := logger.FromContext(ctx)
	logr.Debugw("Starting varnish reload...")
	start := time.Now()
	out, err := r.varnish.Reload(createVCLConfigName(cm.GetResourceVersion()), *vc.Spec.VCL.EntrypointFileName)
	if err != nil {
		if strings.Contains(string(out), "VCL compilation failed") {
			r.metrics.VCLCompilationError.Set(1)
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

	r.metrics.VCLCompilationError.Set(0)
	logr.Debugf("VarnishClusterVarnish successfully reloaded in %f seconds", time.Since(start).Seconds())

	logr.Debugf("Cleaning up old VCL configs...")
	cleanedUpVCLs := 0
	configsList, err := r.varnish.List()
	if err != nil {
		return errors.WithStack(err)
	}

	// cleanup unused VCLs. It cleans up only VCLs created by varnish controller (those that start with our prefix)
	for _, vclConfig := range configsList {
		if vclConfig.Status == varnishadm.VCLStatusAvailable && strings.HasPrefix(vclConfig.Name, VCLVersionPrefix) {
			err := r.varnish.Discard(vclConfig.Name)
			if err != nil {
				return errors.Wrapf(err, "Can't delete VCL config %q", vclConfig.Name)
			}
			cleanedUpVCLs++
		}
	}

	logr.Debugf("Cleaned up %d VCL config(s)", cleanedUpVCLs)
	return nil
}

// creates the VarnishClusterVCL config name from config map version
func createVCLConfigName(configMapVersion string) string {
	return fmt.Sprintf("%s%s-%d", VCLVersionPrefix, configMapVersion, time.Now().Unix())
}
