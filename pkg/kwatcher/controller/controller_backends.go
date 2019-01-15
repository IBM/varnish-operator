package controller

import (
	"context"
	"icm-varnish-k8s-operator/pkg/kwatcher/backends"

	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileVarnish) getBackends(namespace string, labels labels.Selector, targetPort int32) ([]backends.Backend, error) {
	found := &v1.EndpointsList{}

	opts := &client.ListOptions{
		LabelSelector: labels,
		Namespace:     namespace,
	}
	err := r.List(context.TODO(), opts, found)
	if err != nil {
		return nil, errors.Annotatef(err, "could not retrieve backends from namespace %s with labels %s", namespace, labels.String())
	}

	if len(found.Items) == 0 {
		return nil, errors.NotFoundf("no endpoints from namespace %s matching labels %s", namespace, labels.String())
	}

	var backendList []backends.Backend

	for _, endpoints := range found.Items {
		for _, endpoint := range endpoints.Subsets {
			for _, address := range endpoint.Addresses {
				for _, port := range endpoint.Ports {
					if port.Port == targetPort {
						nodeLabels, err := r.getNodeLabels(*address.NodeName)
						if err != nil {
							return nil, errors.Trace(err)
						}
						b := backends.Backend{IP: address.IP, NodeLabels: nodeLabels, PodName: address.TargetRef.Name}
						backendList = append(backendList, b)
						break
					}
				}
			}
		}
	}
	return backendList, nil
}
