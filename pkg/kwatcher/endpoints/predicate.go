package endpoints

import (
	"icm-varnish-k8s-operator/pkg/logger"
	"sort"

	"github.com/google/go-cmp/cmp"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &EndpointPredicate{}

type EndpointPredicate struct {
	labels map[string]string
	logger *logger.Logger
}

func NewPredicate(selectorString string, logr *logger.Logger) (predicate.Predicate, error) {
	ls, err := labels.ConvertSelectorToLabelsMap(selectorString)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse selector string")
	}
	ep := &EndpointPredicate{
		labels: ls,
		logger: logr,
	}
	return ep, nil
}

func (ep *EndpointPredicate) shared(meta metav1.Object) bool {
	inc := meta.GetLabels()
	for label, v := range ep.labels {
		if incV, found := inc[label]; !found || v != incV {
			return false
		}
	}
	return true
}

func (ep *EndpointPredicate) Create(e event.CreateEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *EndpointPredicate) Delete(e event.DeleteEvent) bool {
	return ep.shared(e.Meta)
}

func (ep *EndpointPredicate) Update(e event.UpdateEvent) bool {
	if !ep.shared(e.MetaNew) {
		return false
	}

	newEndpoints, ok := e.ObjectNew.(*v1.Endpoints)
	if !ok {
		ep.logger.Errorf("Wrong object type. Got %T Expected %T", e.ObjectNew, newEndpoints)
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

func (ep *EndpointPredicate) Generic(e event.GenericEvent) bool {
	return ep.shared(e.Meta)
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
