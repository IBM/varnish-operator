package controller

import (
	"context"
	"fmt"
	"sort"

	"github.com/cin/varnish-operator/pkg/varnishcontroller/podutil"

	"github.com/cin/varnish-operator/pkg/varnishcontroller/events"

	"github.com/cin/varnish-operator/api/v1alpha1"
	vclabels "github.com/cin/varnish-operator/pkg/labels"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileVarnish) getBackendEndpoints(ctx context.Context, vc *v1alpha1.VarnishCluster) ([]PodInfo, int32, float64, float64, error) {
	varnishNodeLabels, err := r.getNodeLabels(ctx, r.config.NodeName)
	if err != nil {
		return nil, 0, 0, 0, errors.WithStack(err)
	}

	// Check for deprecated topology labels
	zoneLabel := v1.LabelTopologyZone
	if _, ok := varnishNodeLabels[v1.LabelFailureDomainBetaZone]; ok {
		zoneLabel = v1.LabelFailureDomainBetaZone
	}

	currentZone := varnishNodeLabels[zoneLabel]

	actualLocalWeight := 1.0
	actualRemoteWeight := 1.0

	ns := []string{r.config.Namespace}
	if len(vc.Spec.Backend.Namespaces) > 0 {
		ns = vc.Spec.Backend.Namespaces
	}

	selector := labels.SelectorFromSet(vc.Spec.Backend.Selector)
	backendList, portNumber, err := r.getPodsInfo(ctx, vc, ns, selector, *vc.Spec.Backend.Port, vc.Spec.Backend.OnlyReady)
	if err != nil {
		return nil, 0, 0, 0, errors.WithStack(err)
	}

	if !checkMultizone(backendList, zoneLabel, currentZone) {
		return backendList, portNumber, actualLocalWeight, actualRemoteWeight, nil
	}

	backendRatio := calculateBackendRatio(backendList, currentZone, zoneLabel)

	switch vc.Spec.Backend.ZoneBalancing.Type {
	case v1alpha1.VarnishClusterBackendZoneBalancingTypeAuto:
		baseLocalWeight := 10
		for i, backend := range backendList {
			if backend.NodeLabels[zoneLabel] == currentZone {
				actualLocalWeight = float64(baseLocalWeight) * backendRatio
				backendList[i].Weight = actualLocalWeight
			} else {
				actualRemoteWeight = float64(1)
			}
		}

	case v1alpha1.VarnishClusterBackendZoneBalancingTypeThresholds:
		thresholds := vc.Spec.Backend.ZoneBalancing.Thresholds

		if len(thresholds) < 1 {
			break
		}

		sort.Slice(thresholds, func(i, j int) bool { return *thresholds[i].Threshold > *thresholds[j].Threshold })

		currentLocalWeight := *thresholds[0].Local
		currentRemoteWeight := *thresholds[0].Remote

		for _, thd := range thresholds {
			if int(backendRatio*100) <= *thd.Threshold {
				currentLocalWeight = *thd.Local
				currentRemoteWeight = *thd.Remote
			} else {
				break
			}
		}

		for i, backend := range backendList {
			if backend.NodeLabels[zoneLabel] == currentZone {
				actualLocalWeight = float64(currentLocalWeight)
				backendList[i].Weight = actualLocalWeight
			} else {
				actualRemoteWeight = float64(currentRemoteWeight)
				backendList[i].Weight = actualRemoteWeight
			}
		}
	default:
		// When Zone balancing is disabled we don't need to modify backend weight
		// since they already have equal weight
		break
	}

	return backendList, portNumber, actualLocalWeight, actualRemoteWeight, nil
}

func (r *ReconcileVarnish) getVarnishEndpoints(ctx context.Context, vc *v1alpha1.VarnishCluster) ([]PodInfo, error) {
	varnishLables := labels.SelectorFromSet(vclabels.CombinedComponentLabels(vc, v1alpha1.VarnishComponentVarnish))
	varnishPort := intstr.FromString(v1alpha1.VarnishPortName)

	varnishEndpoints, _, err := r.getPodsInfo(ctx, vc, []string{r.config.Namespace}, varnishLables, varnishPort, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return varnishEndpoints, nil
}

func (r *ReconcileVarnish) getPodsInfo(ctx context.Context, vc *v1alpha1.VarnishCluster, namespaces []string, labels labels.Selector, validPort intstr.IntOrString, onlyReady bool) ([]PodInfo, int32, error) {
	var pods []v1.Pod
	for _, namespace := range namespaces {
		listOptions := []client.ListOption{
			client.MatchingLabelsSelector{Selector: labels},
			client.InNamespace(namespace),
		}
		found := &v1.PodList{}
		err := r.List(ctx, found, listOptions...)
		if err != nil {
			return nil, 0, errors.Wrapf(err, "could not retrieve endpoints from namespace %v with labels %s", namespaces, labels.String())
		}

		pods = append(pods, found.Items...)
	}

	var portNumber int32
	var podInfoList []PodInfo

	if len(pods) == 0 {
		r.logger.Infof("No pods found by labels %v in namespace(s) %v", labels.String(), namespaces)
		return podInfoList, 0, nil
	}

	for _, pod := range pods {
		if len(pod.Status.PodIP) == 0 || len(pod.Spec.NodeName) == 0 {
			continue
		}

		if onlyReady && !podutil.PodReady(pod) {
			continue
		}

		portFound := false
		for _, container := range pod.Spec.Containers {
			for _, containerPort := range container.Ports {
				if containerPort.ContainerPort == validPort.IntVal || containerPort.Name == validPort.StrVal {
					portFound = true
					portNumber = containerPort.ContainerPort
					var backendWeight = 1.0
					nodeLabels, err := r.getNodeLabels(ctx, pod.Spec.NodeName)
					if err != nil {
						return nil, 0, errors.WithStack(err)
					}
					b := PodInfo{IP: pod.Status.PodIP, NodeLabels: nodeLabels, PodName: pod.Name, Weight: backendWeight}
					podInfoList = append(podInfoList, b)
					break
				}
			}
		}

		if !portFound {
			errMsg := fmt.Sprintf("Backend pod %s/%s ignored since none of its containers have port %q defined", pod.Namespace, pod.Name, validPort.String())
			r.eventHandler.Warning(vc, events.EventReasonBackendIgnored, errMsg)
			r.logger.Warnf(errMsg)
		}
	}

	// sort slices so they also look the same for the code using it
	// prevents cases when generated VCL files change only because
	// of order changes in the slice
	sort.SliceStable(podInfoList, func(i, j int) bool {
		return podInfoList[i].IP < podInfoList[j].IP
	})

	return podInfoList, portNumber, nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func calculateBackendRatio(backends []PodInfo, currentZone string, zoneLabel string) float64 {
	var zones []string
	var remoteCount, localCount int
	for _, b := range backends {
		if _, ok := b.NodeLabels[zoneLabel]; ok {
			if !containsString(zones, b.NodeLabels[zoneLabel]) {
				zones = append(zones, b.NodeLabels[zoneLabel])
			}
			if b.NodeLabels[zoneLabel] == currentZone {
				localCount++
			} else {
				remoteCount++
			}
		}
	}

	zoneCount := len(zones)

	// current zone 1 pods / (( remote zone 2 pods + remote zone 3 pods ) / num of remote zones )
	backendRatio := float64(localCount) / (float64(remoteCount) / float64(zoneCount-1))

	return backendRatio
}

func checkMultizone(endpoints []PodInfo, zoneLabel string, currentZone string) bool {
	for _, b := range endpoints {
		if _, ok := b.NodeLabels[zoneLabel]; ok {
			if b.NodeLabels[zoneLabel] != currentZone {
				return true
			}
		}
	}
	return false
}
