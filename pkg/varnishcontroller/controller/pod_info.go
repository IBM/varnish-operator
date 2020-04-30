package controller

import (
	"context"
	"sort"

	"icm-varnish-k8s-operator/api/v1alpha1"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileVarnish) getPodInfo(ctx context.Context, namespace string, labels labels.Selector, validPort intstr.IntOrString, vc *v1alpha1.VarnishCluster) ([]PodInfo, int32, float64, float64, error) {

	found := &v1.EndpointsList{}
	err := r.List(ctx, found, client.MatchingLabelsSelector{Selector: labels}, client.InNamespace(namespace))
	if err != nil {
		return nil, 0, 0, 0, errors.Wrapf(err, "could not retrieve backends from namespace %s with labels %s", namespace, labels.String())
	}

	if len(found.Items) == 0 {
		return nil, 0, 0, 0, errors.Errorf("no endpoints from namespace %s matching labels %s", namespace, labels.String())
	}

	varnishNodeLabels, err := r.getNodeLabels(ctx, r.config.NodeName)
	if err != nil {
		return nil, 0, 0, 0, errors.WithStack(err)
	}

	// Check for deprecated topology labels
	zoneLabel := "topology.kubernetes.io/zone"
	if _, ok := varnishNodeLabels["failure-domain.beta.kubernetes.io/zone"]; ok {
		zoneLabel = "failure-domain.beta.kubernetes.io/zone"
	}

	currentZone := varnishNodeLabels[zoneLabel]

	actualLocalWeight := 1.0
	actualRemoteWeight := 1.0

	var backendList []PodInfo
	var portNumber int32
	var multiZone bool

	for _, endpoints := range found.Items {
		for _, endpoint := range endpoints.Subsets {
			for _, address := range append(endpoint.Addresses, endpoint.NotReadyAddresses...) {
				for _, port := range endpoint.Ports {
					if port.Port == validPort.IntVal || port.Name == validPort.StrVal {
						var backendWeight float64 = 1.0
						portNumber = port.Port
						nodeLabels, err := r.getNodeLabels(ctx, *address.NodeName)
						if err != nil {
							return nil, 0, 0, 0, errors.WithStack(err)
						}
						if _, ok := nodeLabels[zoneLabel]; ok {
							if nodeLabels[zoneLabel] != currentZone {
								multiZone = true
							}
						}
						b := PodInfo{IP: address.IP, NodeLabels: nodeLabels, PodName: address.TargetRef.Name, Weight: backendWeight}
						backendList = append(backendList, b)
						break
					}
				}
			}
		}
	}

	// sort slices so they also look the same for the code using it
	// prevents cases when generated VCL files change only because
	// of order changes in the slice
	sort.SliceStable(backendList, func(i, j int) bool {
		return backendList[i].IP < backendList[j].IP
	})

	if !multiZone {
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
