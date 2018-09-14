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

func (r *ReconcileVarnishService) reconcileRoleBinding(instance *icmapiv1alpha1.VarnishService, roleName, serviceAccountName string) error {
	roleBinding := &rbacv1beta1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-rolebinding",
			Namespace: instance.Namespace,
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

	// Set controller reference for roleBinding
	if err := controllerutil.SetControllerReference(instance, roleBinding, r.scheme); err != nil {
		return logger.RError(err, "Cannot set controller reference for service", "namespace", roleBinding.Namespace, "name", roleBinding.Name)
	}

	found := &rbacv1beta1.RoleBinding{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}, found)
	// If the role does not exist, create it
	if err != nil && kerrors.IsNotFound(err) {
		logger.Info("Creating roleBinding", "name", roleBinding.Name, "namespace", roleBinding.Namespace)
		if err = r.Create(context.TODO(), roleBinding); err != nil {
			return logger.RError(err, "Unable to create roleBinding")
		}
		// If there was a problem doing the GET, just return
	} else if err != nil {
		return logger.RError(err, "Could not Get roleBinding")
		// If the roleBinding exists, and it is different, update
	} else if !reflect.DeepEqual(found.Subjects, roleBinding.Subjects) || !reflect.DeepEqual(found.RoleRef, roleBinding.RoleRef) {
		found.Subjects = roleBinding.Subjects
		found.RoleRef = roleBinding.RoleRef
		logger.Info("Updating roleBinding", "name", roleBinding.Name, "namespace", roleBinding.Namespace)
		if err = r.Update(context.TODO(), found); err != nil {
			return logger.RError(err, "Could not Update roleBinding")
		}
	}
	// If no changes, do nothing
	logger.Info("no updates for rolebinding", "name", roleBinding.Name, "namespace", roleBinding.Namespace)
	return nil
}
