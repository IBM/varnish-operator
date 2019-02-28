package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"

	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileServiceAccount(instance *icmapiv1alpha1.VarnishService) (string, error) {
	saName := instance.Name + "-varnish-serviceaccount"
	serviceAccount := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saName,
			Namespace: instance.Namespace,
			Labels:    labels.CombinedComponentLabels(instance, icmapiv1alpha1.VarnishComponentServiceAccount),
		},
		ImagePullSecrets: []v1.LocalObjectReference{{Name: instance.Spec.Deployment.Container.ImagePullSecret}},
	}

	logr := r.logger.With("name", serviceAccount.Name, "namespace", serviceAccount.Namespace)

	// Set controller reference for service object
	if err := controllerutil.SetControllerReference(instance, serviceAccount, r.scheme); err != nil {
		return "", logr.RErrorw(err, "Cannot set controller reference for Service account")
	}

	found := &v1.ServiceAccount{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: serviceAccount.Name, Namespace: serviceAccount.Namespace}, found)
	// If the service account does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the service exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Service sccount", "new", serviceAccount)
		if err = r.Create(context.TODO(), serviceAccount); err != nil {
			return "", logr.RErrorw(err, "Unable to create Service account")
		}
	} else if err != nil {
		return "", logr.RErrorw(err, "Could not get Service account")
	} else if !compare.EqualServiceAccount(found, serviceAccount) {
		logr.Infoc("Updating Service account", "diff", compare.DiffServiceAccount(found, serviceAccount))
		found.ImagePullSecrets = serviceAccount.ImagePullSecrets
		found.Labels = serviceAccount.Labels
		if err = r.Update(context.TODO(), found); err != nil {
			return "", logr.RErrorw(err, "Unable to update service account")
		}
	} else {
		logr.Debugw("No updates for Service account")
	}
	return saName, nil
}
