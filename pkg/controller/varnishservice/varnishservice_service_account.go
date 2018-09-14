package varnishservice

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"reflect"

	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileServiceAccount(instance *icmapiv1alpha1.VarnishService) (string, error) {
	saName := instance.Name + "-serviceaccount"
	serviceAccount := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      saName,
			Namespace: instance.Namespace,
		},
		ImagePullSecrets: []v1.LocalObjectReference{{Name: instance.Spec.Deployment.ImagePullSecretName}},
	}
	// Set controller reference for service object
	if err := controllerutil.SetControllerReference(instance, serviceAccount, r.scheme); err != nil {
		return "", logger.RError(err, "Cannot set controller reference for service account", "namespace", serviceAccount.Namespace, "name", serviceAccount.Name)
	}

	found := &v1.ServiceAccount{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: serviceAccount.Name, Namespace: serviceAccount.Namespace}, found)
	// If the service account does not exist, create it
	if err != nil && kerrors.IsNotFound(err) {
		logger.Info("Creating service", "config", serviceAccount)
		// logger.Info("Creating service", "namespace", service.Namespace, "name", service.Name)
		if err = r.Create(context.TODO(), serviceAccount); err != nil {
			return "", logger.RError(err, "Unable to create service account")
		}
		// If there was a problem doing the GET, just return
	} else if err != nil {
		return "", logger.RError(err, "Could not Get service account")
		// If the service exists, and it is different, update
	} else if !reflect.DeepEqual(serviceAccount.ImagePullSecrets, found.ImagePullSecrets) {
		found.ImagePullSecrets = serviceAccount.ImagePullSecrets
		logger.Info("Updating service account", "config", found)
		// logger.Info("Updating service", "namespace", service.Namespace, "name", service.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return "", logger.RError(err, "Unable to update service")
		}
	}
	// If no changes, do nothing
	logger.Info("No updates for service account", "name", serviceAccount.Name, "namespace", serviceAccount.Namespace)
	return saName, nil
}
