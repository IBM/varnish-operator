package controller

import (
	"context"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/labels"
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

func (r *ReconcileVarnishCluster) reconcileRole(ctx context.Context, instance *vcapi.VarnishCluster) error {
	role := &rbac.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.Role(instance.Name),
			Namespace: instance.Namespace,
			Labels:    labels.CombinedComponentLabels(instance, vcapi.VarnishComponentRole),
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints", "configmaps"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{"caching.ibm.com"},
				Resources: []string{"varnishclusters"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "patch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "watch"},
			},
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentRole)
	logr = logr.With(logger.FieldComponentName, role.Name)

	// Set controller reference for role
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return errors.Wrap(err, "Cannot set controller reference for service")
	}

	found := &rbac.Role{}

	err := r.Get(ctx, types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, found)
	// If the role does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the role exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Role", "new", role)
		if err = r.Create(ctx, role); err != nil {
			return errors.Wrap(err, "Unable to create role")
		}
	} else if err != nil {
		return errors.Wrap(err, "Could not Get role")
	} else if !compare.EqualRole(found, role) {
		logr.Infoc("Updating Role", "diff", compare.DiffRole(found, role))
		found.Rules = role.Rules
		found.Labels = role.Labels
		if err = r.Update(ctx, found); err != nil {
			return errors.Wrap(err, "Could not Update role")
		}
	} else {
		logr.Debugw("No updates for Role")

	}
	return nil
}
