package controller

import (
	"context"

	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileVarnish) getNodeLabels(nodeName string) (map[string]string, error) {
	found := &v1.Node{}
	err := r.Get(context.TODO(), types.NamespacedName{Namespace: v1.NamespaceAll, Name: nodeName}, found)
	if err != nil && kerrors.IsNotFound(err) {
		return nil, errors.Annotatef(err, "could not find node with name %s", nodeName)
	} else if err != nil {
		return nil, errors.Annotatef(err, "problem calling Get on Node %s", nodeName)
	} else {
		return found.Labels, nil
	}
}
