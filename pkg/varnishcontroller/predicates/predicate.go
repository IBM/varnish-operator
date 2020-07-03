package predicates

import (
	"github.com/google/go-cmp/cmp"
	"github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sort"
)

var _ predicate.Predicate = &varnishControllerPredicate{}

type varnishControllerPredicate struct {
	clusterUID         types.UID
	logger             *logger.Logger
	endpointsSelectors []labels.Selector
}

func NewVarnishControllerPredicate(clusterUID types.UID, endpointsSelectors []labels.Selector, logr *logger.Logger) predicate.Predicate {
	if logr == nil {
		logr = logger.NewNopLogger()
	}
	return &varnishControllerPredicate{
		clusterUID:         clusterUID,
		logger:             logr,
		endpointsSelectors: endpointsSelectors,
	}
}

func (p *varnishControllerPredicate) Create(e event.CreateEvent) bool {
	switch v := e.Object.(type) {
	case *v1alpha1.VarnishCluster:
		if e.Meta.GetUID() != p.clusterUID {
			return false
		}
	case *v1.Endpoints:
		if !p.endpointMatchesSelector(v) {
			return false
		}
	}

	p.logger.Debugf("Create event for resource %T: %s/%s", e.Object, e.Meta.GetNamespace(), e.Meta.GetName())
	return true
}

func (p *varnishControllerPredicate) Delete(e event.DeleteEvent) bool {
	switch e.Object.(type) {
	case *v1alpha1.VarnishCluster:
		return false
	case *v1.Endpoints:
		return false
	}

	p.logger.Debugf("Delete event for resource %T: %s/%s", e.Object, e.Meta.GetNamespace(), e.Meta.GetName())
	return true
}

func (p *varnishControllerPredicate) Update(e event.UpdateEvent) bool {
	if reflect.TypeOf(e.ObjectNew) != reflect.TypeOf(e.ObjectOld) {
		p.logger.Errorf("New and Old object kinds are different. New: %T, Old: %T", e.ObjectNew, e.ObjectOld)
		return false
	}

	switch e.ObjectNew.(type) {
	case *v1alpha1.VarnishCluster:
		if !p.allowVarnishClusterUpdateEvent(e) {
			return false
		}
	case *v1.Endpoints:
		if !p.allowEndpointsUpdateEvent(e) {
			return false
		}
	}

	p.logger.Debugf("Update event for resource %T: %s/%s", e.ObjectNew, e.MetaNew.GetNamespace(), e.MetaNew.GetName())
	return true
}

func (p *varnishControllerPredicate) Generic(e event.GenericEvent) bool {
	switch endpoint := e.Object.(type) {
	case *v1alpha1.VarnishCluster:
		if e.Meta.GetUID() != p.clusterUID {
			return false
		}
	case *v1.Endpoints:
		if !p.endpointMatchesSelector(endpoint) {
			return false
		}
	}

	p.logger.Debugf("Generic event for resource %T: %s/%s", e.Object, e.Meta.GetNamespace(), e.Meta.GetName())
	return true
}

func (p *varnishControllerPredicate) allowVarnishClusterUpdateEvent(e event.UpdateEvent) bool {
	newCluster := e.ObjectNew.(*v1alpha1.VarnishCluster)
	oldCluster := e.ObjectOld.(*v1alpha1.VarnishCluster)
	if e.MetaNew.GetUID() != p.clusterUID || e.MetaOld.GetUID() != p.clusterUID {
		return false
	}

	if newCluster.Status.VCL.ConfigMapVersion != oldCluster.Status.VCL.ConfigMapVersion {
		return true
	}

	if !cmp.Equal(newCluster.Spec.Backend, oldCluster.Spec.Backend) {
		return true
	}
	return false
}

func (p *varnishControllerPredicate) allowEndpointsUpdateEvent(e event.UpdateEvent) bool {
	newEndpoints := e.ObjectNew.(*v1.Endpoints)
	oldEndpoints := e.ObjectOld.(*v1.Endpoints)
	if !p.endpointMatchesSelector(newEndpoints) {
		return false
	}

	if cmp.Equal(getIPs(oldEndpoints.Subsets), getIPs(newEndpoints.Subsets)) &&
		cmp.Equal(getPorts(oldEndpoints.Subsets), getPorts(newEndpoints.Subsets)) {
		return false
	}
	return true
}

func (p *varnishControllerPredicate) endpointMatchesSelector(endpoint *v1.Endpoints) bool {
	for _, selector := range p.endpointsSelectors {
		if selector.Matches(labels.Set(endpoint.GetLabels())) {
			return true
		}
	}

	return false
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
