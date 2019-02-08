package controller

import (
	"context"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	vslabels "icm-varnish-k8s-operator/pkg/labels"
	"sort"

	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileVarnish) getVarnishNodes(vs *v1alpha1.VarnishService) ([]VarnishNode, error) {
	found := &v1.EndpointsList{}
	varnishEndpointsSelector := labels.SelectorFromSet(vslabels.CombinedComponentLabels(vs, v1alpha1.VarnishComponentCachedService))

	err := r.List(context.TODO(), &client.ListOptions{LabelSelector: varnishEndpointsSelector, Namespace: vs.Namespace}, found)
	if err != nil {
		return nil, errors.Annotatef(err, "could not retrieve varnish nodes from namespace %s for varnishservice %s", vs.Namespace, vs.Name)
	}

	if len(found.Items) == 0 {
		return nil, errors.NotFoundf("no endpoints from namespace %s for varnishservice %s", vs.Namespace, vs.Name)
	}

	var varnishNodes []VarnishNode
	for _, endpoints := range found.Items {
		for _, endpoint := range endpoints.Subsets {
			for _, address := range append(endpoint.Addresses, endpoint.NotReadyAddresses...) {
				for _, port := range endpoint.Ports {
					if port.Name == vs.Spec.Service.VarnishPort.Name {
						varnishNode := VarnishNode{IP: address.IP, Port: port.Port, PodName: address.TargetRef.Name}
						varnishNodes = append(varnishNodes, varnishNode)
						break
					}
				}
			}
		}
	}

	// sort slices so they also look the same for the code using it
	// prevents cases when generated VCL files change only because of order changes in the slice
	sort.SliceStable(varnishNodes, func(i, j int) bool {
		return varnishNodes[i].IP < varnishNodes[j].IP
	})

	return varnishNodes, nil
}
