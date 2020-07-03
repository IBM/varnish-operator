package controller

import (
	"context"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"github.com/pkg/errors"

	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishCluster) reconcileServiceAccount(ctx context.Context, instance *vcapi.VarnishCluster) error {
	saName := names.ServiceAccount(instance.Name)
	serviceAccount := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saName,
			Namespace: instance.Namespace,
			Labels:    labels.CombinedComponentLabels(instance, vcapi.VarnishComponentServiceAccount),
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentServiceAccount)
	logr = logr.With(logger.FieldComponentName, serviceAccount.Name)

	// Set controller reference for service object
	if err := controllerutil.SetControllerReference(instance, serviceAccount, r.scheme); err != nil {
		return errors.Wrap(err, "Cannot set controller reference for Service account")
	}

	found := &v1.ServiceAccount{}

	err := r.Get(ctx, types.NamespacedName{Name: serviceAccount.Name, Namespace: serviceAccount.Namespace}, found)
	// If the service account does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the service exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Service account", "new", serviceAccount)
		if err = r.Create(ctx, serviceAccount); err != nil {
			return errors.Wrap(err, "Unable to create Service account")
		}
	} else if err != nil {
		return errors.Wrap(err, "Could not get Service account")
	} else if !compare.EqualServiceAccount(found, serviceAccount) {
		logr.Infoc("Updating Service account", "diff", compare.DiffServiceAccount(found, serviceAccount))
		found.ImagePullSecrets = serviceAccount.ImagePullSecrets
		found.Labels = serviceAccount.Labels
		if err = r.Update(ctx, found); err != nil {
			return errors.Wrap(err, "Unable to update service account")
		}
	} else {
		logr.Debugw("No updates for Service account")
	}
	return nil
}
