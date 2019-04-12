package controller

import (
	"context"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileVarnish) getNodeLabels(ctx context.Context, nodeName string) (map[string]string, error) {
	found := &v1.Node{}
	err := r.Get(ctx, types.NamespacedName{Namespace: v1.NamespaceAll, Name: nodeName}, found)
	if err != nil && kerrors.IsNotFound(err) {
		return nil, errors.Wrapf(err, "could not find node with name %s", nodeName)
	} else if err != nil {
		return nil, errors.Wrapf(err, "problem calling Get on Node %s", nodeName)
	} else {
		return found.Labels, nil
	}
}
