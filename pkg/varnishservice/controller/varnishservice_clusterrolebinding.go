package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"
	"icm-varnish-k8s-operator/pkg/varnishservice/logger"

	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileClusterRoleBinding(instance *icmapiv1alpha1.VarnishService, roleName, serviceAccountName string) error {
	clusterRoleBinding := &rbacv1beta1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   instance.Name + "-varnish-clusterrolebinding-" + instance.Namespace,
			Labels: combinedLabels(instance, "clusterrolebinding"),
		},
		Subjects: []rbacv1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: instance.Namespace,
			},
		},
		RoleRef: rbacv1beta1.RoleRef{
			Kind:     "ClusterRole",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	logr := logger.With("name", clusterRoleBinding.Name, "namespace", clusterRoleBinding.Namespace)

	// Set controller reference for clusterRoleBinding
	if err := controllerutil.SetControllerReference(instance, clusterRoleBinding, r.scheme); err != nil {
		return logr.RErrorw(err, "Cannot set controller reference for ClusterRoleBinding")
	}

	found := &rbacv1beta1.ClusterRoleBinding{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: clusterRoleBinding.Name, Namespace: clusterRoleBinding.Namespace}, found)
	// If the role does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the clusterRoleBinding exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating ClusterRoleBinding", "new", clusterRoleBinding)
		if err = r.Create(context.TODO(), clusterRoleBinding); err != nil {
			return logr.RErrorw(err, "Unable to create ClusterRoleBinding")
		}
	} else if err != nil {
		return logr.RErrorw(err, "Could not Get ClusterRoleBinding")
	} else if !compare.EqualClusterRoleBinding(found, clusterRoleBinding) {
		logr.Infoc("Updating ClusterRoleBinding", "diff", compare.DiffClusterRoleBinding(found, clusterRoleBinding))
		found.Subjects = clusterRoleBinding.Subjects
		found.RoleRef = clusterRoleBinding.RoleRef
		found.Labels = clusterRoleBinding.Labels
		if err = r.Update(context.TODO(), found); err != nil {
			return logr.RErrorw(err, "Could not Update ClusterRoleBinding")
		}
	} else {
		logr.Debugw("No updates for ClusterRoleBinding")
	}
	return nil
}
