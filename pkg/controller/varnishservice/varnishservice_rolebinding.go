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

func (r *ReconcileVarnishService) reconcileRoleBinding(instance *icmapiv1alpha1.VarnishService, roleName, serviceAccountName string) error {
	roleBinding := &rbacv1beta1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-varnish-rolebinding",
			Namespace: instance.Namespace,
			Labels:    combinedLabels(instance, "rolebinding"),
		},
		Subjects: []rbacv1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: instance.Namespace,
			},
		},
		RoleRef: rbacv1beta1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	logr := logger.With("name", roleBinding.Name, "namespace", roleBinding.Namespace)

	// Set controller reference for roleBinding
	if err := controllerutil.SetControllerReference(instance, roleBinding, r.scheme); err != nil {
		return logr.RErrorw(err, "Cannot set controller reference for service")
	}

	found := &rbacv1beta1.RoleBinding{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}, found)
	// If the role does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the roleBinding exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating RoleBinding", "new", roleBinding)
		if err = r.Create(context.TODO(), roleBinding); err != nil {
			return logr.RErrorw(err, "Unable to create roleBinding")
		}
	} else if err != nil {
		return logr.RErrorw(err, "Could not Get roleBinding")
	} else if !compare.EqualRoleBinding(found, roleBinding) {
		logr.Debugw("Updating RoleBinding", "diff", compare.DiffRoleBinding(found, roleBinding))
		found.Subjects = roleBinding.Subjects
		found.RoleRef = roleBinding.RoleRef
		found.Labels = roleBinding.Labels
		if err = r.Update(context.TODO(), found); err != nil {
			return logr.RErrorw(err, "Could not Update roleBinding")
		}
	} else {
		logr.Debugw("No updates for Rolebinding")
	}
	return nil
}
