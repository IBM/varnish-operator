package controller

import (
	"context"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"github.com/pkg/errors"

	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishCluster) reconcileRoleBinding(ctx context.Context, instance *vcapi.VarnishCluster) error {
	roleBinding := &rbac.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.RoleBinding(instance.Name),
			Namespace: instance.Namespace,
			Labels:    vclabels.CombinedComponentLabels(instance, vcapi.VarnishComponentRoleBinding),
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      names.ServiceAccount(instance.Name),
				Namespace: instance.Namespace,
			},
		},
		RoleRef: rbac.RoleRef{
			Kind:     "Role",
			Name:     names.Role(instance.Name),
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentRoleBinding)
	logr = logr.With(logger.FieldComponentName, roleBinding.Name)

	// Set controller reference for roleBinding
	if err := controllerutil.SetControllerReference(instance, roleBinding, r.scheme); err != nil {
		return errors.Wrap(err, "Cannot set controller reference for service")
	}

	found := &rbac.RoleBinding{}

	err := r.Get(ctx, types.NamespacedName{Name: roleBinding.Name, Namespace: roleBinding.Namespace}, found)
	// If the role does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the roleBinding exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating RoleBinding", "new", roleBinding)
		if err = r.Create(ctx, roleBinding); err != nil {
			return errors.Wrap(err, "Unable to create roleBinding")
		}
	} else if err != nil {
		return errors.Wrap(err, "Could not Get roleBinding")
	} else if !compare.EqualRoleBinding(found, roleBinding) {
		logr.Infoc("Updating RoleBinding", "diff", compare.DiffRoleBinding(found, roleBinding))
		found.Subjects = roleBinding.Subjects
		found.RoleRef = roleBinding.RoleRef
		found.Labels = roleBinding.Labels
		if err = r.Update(ctx, found); err != nil {
			return errors.Wrap(err, "Could not Update roleBinding")
		}
	} else {
		logr.Debugw("No updates for Rolebinding")
	}
	return nil
}
