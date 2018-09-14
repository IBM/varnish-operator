package varnishservice

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"reflect"

	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileRole(instance *icmapiv1alpha1.VarnishService) (string, error) {
	role := &rbacv1beta1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-role",
			Namespace: instance.Namespace,
		},
		Rules: []rbacv1beta1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"list", "watch"},
			},
		},
	}

	// Set controller reference for role
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return "", logger.RError(err, "Cannot set controller reference for service", "namespace", role.Namespace, "name", role.Name)
	}

	found := &rbacv1beta1.Role{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, found)
	// If the role does not exist, create it
	if err != nil && kerrors.IsNotFound(err) {
		// logger.Info("Creating role", "name", role.Name, "namespace", role.Namespace)
		logger.Info("Creating role", "config", role)
		if err = r.Create(context.TODO(), role); err != nil {
			return "", logger.RError(err, "Unable to create role")
		}
		// If there was a problem doing the GET, just return
	} else if err != nil {
		return "", logger.RError(err, "Could not Get role")
		// If the role exists, and it is different, update
	} else if !reflect.DeepEqual(found.Rules, role.Rules) {
		found.Rules = role.Rules
		logger.Info("Updating role", "config", found)
		if err = r.Update(context.TODO(), found); err != nil {
			return "", logger.RError(err, "Could not Update role")
		}
	}
	// If no changes, do nothing
	logger.Info("no updates for role", "name", role.Name, "namespace", role.Namespace)
	return role.Name, nil
}
