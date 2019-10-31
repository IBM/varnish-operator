package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"

	"github.com/pkg/errors"

	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileClusterRole(ctx context.Context, instance *icmapiv1alpha1.VarnishService) (string, error) {
	role := &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   instance.Name + "-varnish-clusterrole-" + instance.Namespace,
			Labels: labels.CombinedComponentLabels(instance, icmapiv1alpha1.VarnishComponentClusterRole),
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"list", "watch"},
			},
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentClusterRole)
	logr = logr.With(logger.FieldComponentName, role.Name)

	// Set controller reference for role
	if err := controllerutil.SetControllerReference(instance, role, r.scheme); err != nil {
		return "", errors.Wrap(err, "Cannot set controller reference for ClusterRole")
	}

	found := &rbac.ClusterRole{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: role.Name, Namespace: role.Namespace}, found)
	// If the role does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the role exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating ClusterRole", "new", role)
		if err = r.Create(ctx, role); err != nil {
			return "", errors.Wrap(err, "Unable to create ClusterRole")
		}
	} else if err != nil {
		return "", errors.Wrap(err, "Could not Get ClusterRole")
	} else if !compare.EqualClusterRole(found, role) {
		logr.Infoc("Updating ClusterRole", "diff", compare.DiffClusterRole(found, role))
		found.Rules = role.Rules
		found.Labels = role.Labels
		if err = r.Update(ctx, found); err != nil {
			return "", errors.Wrap(err, "Could not Update ClusterRole")
		}
	} else {
		logr.Debugw("No updates for ClusterRole")

	}
	return role.Name, nil
}
