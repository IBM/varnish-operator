package varnishservice

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/compare"
	"icm-varnish-k8s-operator/pkg/logger"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileRole(instance *icmapiv1alpha1.VarnishService) (string, error) {
	role := &rbacv1beta1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-varnish-role",
			Namespace: instance.Namespace,
			Labels:    combinedLabels(instance, "role"),
		},
		Rules: []rbacv1beta1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"list", "watch"},
			},
		},
	}

	logr := logger.WithValues("name", role.Name, "namespace", role.Namespace)

	// Set controller reference for role
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return "", logr.RError(err, "Cannot set controller reference for service")
	}

	found := &rbacv1beta1.Role{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, found)
	// If the role does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the role exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Info("Creating role", "new", role)
		if err = r.Create(context.TODO(), role); err != nil {
			return "", logr.RError(err, "Unable to create role")
		}
	} else if err != nil {
		return "", logr.RError(err, "Could not Get role")
	} else if !compare.EqualRole(found, role) {
		logr.Info("Updating role", "diff", compare.DiffRole(found, role))
		found.Rules = role.Rules
		found.Labels = role.Labels
		if err = r.Update(context.TODO(), found); err != nil {
			return "", logr.RError(err, "Could not Update role")
		}
	} else {
		logr.V(5).Info("no updates for role")

	}
	return role.Name, nil
}
