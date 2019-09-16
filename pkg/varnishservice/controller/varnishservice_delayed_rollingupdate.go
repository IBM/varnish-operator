package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	klabels "k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileVarnishService) reconcileDelayedRollingUpdate(ctx context.Context, instance, instanceStatus *icmapiv1alpha1.VarnishService, sts *appsv1.StatefulSet) error {
	logr := logger.FromContext(ctx)

	if instance.Spec.StatefulSet.UpdateStrategy.Type != icmapiv1alpha1.VarnishUpdateStrategyDelayedRollingUpdate {
		// make sure there're no hanging timers.
		// could happen if the update strategy changed from DelayedRollingUpdate to something else
		r.reconcileTriggerer.Stop(icmapiv1alpha1.VarnishUpdateStrategyDelayedRollingUpdate, instance)
		return nil
	}

	if sts.Status.UpdatedReplicas == sts.Status.Replicas {
		logr.Debugf("All replicas are up to date")
		return nil
	}

	stsPods := &v1.PodList{}
	varnishPodLabels := klabels.SelectorFromSet(labels.CombinedComponentLabels(instance, icmapiv1alpha1.VarnishComponentVarnishes))
	if err := r.List(ctx, &client.ListOptions{LabelSelector: varnishPodLabels}, stsPods); err != nil {
		return errors.WithStack(err)
	}

	if len(stsPods.Items) == 0 {
		logr.Infof("No pods found to perform DelayedRollingUpdate")
		return nil
	}

	// an already existing timer means that we have already reloaded a pod and should wait for the timer to trigger the reconcile loop
	// the obvious way of checking the statefulset status proved to be error prone as there is a delay between pod deletion and statefulset status update.
	// During that time the operator thinks that we didn't reload a pod and will delete an another one by mistake.
	if r.reconcileTriggerer.TimerExists(icmapiv1alpha1.VarnishUpdateStrategyDelayedRollingUpdate, instance) {
		logr.Debugf("Timer exists. Waiting for a reconcile loop to be triggered by the timer.")
		return nil
	}

	// Don't reload a new pod if an existing is still not ready, even if it's time to reload an another one.
	// After the pod becomes ready, an another reconcile loop will be triggered and a new pod will be reloaded.
	if sts.Status.ReadyReplicas < sts.Status.Replicas {
		logr.Debugf("One or more pods are not ready yet. Not triggering an another pod update yet.")
		return nil
	}

	currentRevision := sts.Status.UpdateRevision

	var newestUpdatedPod, updateCandidate *v1.Pod
	for i, stsPod := range stsPods.Items { //find the newest pod
		if stsPod.Labels["controller-revision-hash"] != currentRevision {
			updateCandidate = &stsPods.Items[i]
			continue
		}
		if newestUpdatedPod == nil {
			newestUpdatedPod = &stsPods.Items[i]
			continue
		}
		if newestUpdatedPod.Status.StartTime.UnixNano() < stsPod.Status.StartTime.UnixNano() {
			newestUpdatedPod = &stsPods.Items[i]
			continue
		}
	}

	// could happen during scaling operations
	if updateCandidate == nil {
		logr.Debugf("Couldn't find a pod to update.")
		return nil
	}

	rollingUpdateDelay := time.Duration(instance.Spec.StatefulSet.UpdateStrategy.DelayedRollingUpdate.DelaySeconds) * time.Second

	// no pods updated yet. Update the first one without delay.
	if newestUpdatedPod == nil {
		logr.Debugf("Updating the first pod: %s", stsPods.Items[0].Name)
		if err := r.Delete(ctx, &stsPods.Items[0]); err != nil {
			return errors.WithStack(err)
		}
		// The pod's creation time will be later than the reconcile time
		// This will be fixed during next reconcile
		r.reconcileTriggerer.TriggerAfter(icmapiv1alpha1.VarnishUpdateStrategyDelayedRollingUpdate, rollingUpdateDelay, instance)
		return nil
	}

	nextUpdateTime := newestUpdatedPod.Status.StartTime.Add(rollingUpdateDelay)
	if (time.Now().Unix() - nextUpdateTime.Unix()) >= 0 { //if time to update an another pod
		logr.Infof("Updating pod %s according to DelayedRollingUpdate strategy", updateCandidate.Name)
		if err := r.Delete(ctx, updateCandidate); err != nil {
			return errors.WithStack(err)
		}

		logr.Debugf("Setting timer to trigger reconcile loop in %s for ", rollingUpdateDelay)
		// this time will not be correct as Kubernetes will have a delay before it creates a new pod
		// As we can't know when the pod will be created, set the timer delay based on the current time.
		// Later, when the pod will be created, we will set the correct time
		r.reconcileTriggerer.TriggerAfter(icmapiv1alpha1.VarnishUpdateStrategyDelayedRollingUpdate, rollingUpdateDelay, instance)
		return nil
	}

	// make sure we have the correct time
	r.reconcileTriggerer.TriggerAfter(icmapiv1alpha1.VarnishUpdateStrategyDelayedRollingUpdate, time.Until(nextUpdateTime), instance)
	logr.Debugf("DelayedRollingUpdate is in progress. Next pod update is in %s", time.Until(nextUpdateTime))
	return nil
}
