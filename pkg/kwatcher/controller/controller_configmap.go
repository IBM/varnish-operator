package controller

import (
	"context"

	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileVarnish) getConfigMap(namespace, name string) (*v1.ConfigMap, error) {
	found := &v1.ConfigMap{}

	err := r.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, found)
	if err != nil && kerrors.IsNotFound(err) {
		return nil, errors.NewNotFound(err, "ConfigMap must exist to reconcile Varnish")
	} else if err != nil {
		return nil, errors.Annotatef(err, "could not Get ConfigMap")
	}

	return found, nil
}
