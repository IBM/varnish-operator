package predicates

import (
	"icm-varnish-k8s-operator/pkg/logger"
	"sort"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/google/go-cmp/cmp"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &endpointSelectorPredicate{}

type endpointSelectorPredicate struct {
	selectors []labels.Selector
	logger    *logger.Logger
}

// NewEndpointsSelectors filters out endpoints that doesn't match any of the provided selectors.
func NewEndpointsSelectors(selector []labels.Selector, logr *logger.Logger) predicate.Predicate {
	return &endpointSelectorPredicate{
		selectors: selector,
		logger:    logr,
	}
}

func (ep *endpointSelectorPredicate) shared(obj runtime.Object, meta metav1.Object) bool {
	if _, ok := obj.(*v1.Endpoints); !ok {
		return true
	}

	for _, selector := range ep.selectors {
		if selector.Matches(labels.Set(meta.GetLabels())) {
			return true
		}
	}

	return false
}

func (ep *endpointSelectorPredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Object, e.Meta)
}

func (ep *endpointSelectorPredicate) Delete(e event.DeleteEvent) bool {
	return false // happens only when the VarnishCluster is deleted. Don't do anything in that case.
}

func (ep *endpointSelectorPredicate) Update(e event.UpdateEvent) bool {
	newEndpoints, ok := e.ObjectNew.(*v1.Endpoints)
	if !ok {
		return true
	}

	found := false
	for _, selector := range ep.selectors {
		if selector.Matches(labels.Set(newEndpoints.GetLabels())) {
			found = true
			break
		}
	}

	if !found {
		return false
	}

	oldEndpoints, ok := e.ObjectOld.(*v1.Endpoints)
	if !ok {
		ep.logger.Errorf("Wrong object type. Got %T Expected %T", e.ObjectNew, oldEndpoints)
		return false
	}

	if cmp.Equal(getIPs(oldEndpoints.Subsets), getIPs(newEndpoints.Subsets)) &&
		cmp.Equal(getPorts(oldEndpoints.Subsets), getPorts(newEndpoints.Subsets)) {
		return false
	}

	return true
}

func (ep *endpointSelectorPredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Object, e.Meta)
}

func getIPs(eps []v1.EndpointSubset) []string {
	ips := make([]string, 0)
	for _, ep := range eps {
		for _, addr := range append(ep.Addresses, ep.NotReadyAddresses...) {
			ips = append(ips, addr.IP)
		}

	}
	sort.Strings(ips)
	return ips
}

func getPorts(eps []v1.EndpointSubset) []int32 {
	ports := make([]int32, 0)
	for _, ep := range eps {
		for _, port := range ep.Ports {
			ports = append(ports, port.Port)
		}

	}
	sort.Slice(ports, func(i, j int) bool {
		return ports[i] < ports[j]
	})
	return ports
}
