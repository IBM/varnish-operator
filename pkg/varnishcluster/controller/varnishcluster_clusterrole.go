package controller

import (
	"context"

	vcapi "github.com/cin/varnish-operator/api/v1alpha1"
	"github.com/cin/varnish-operator/pkg/labels"
	"github.com/cin/varnish-operator/pkg/logger"
	"github.com/cin/varnish-operator/pkg/names"
	"github.com/cin/varnish-operator/pkg/varnishcluster/compare"

	"github.com/pkg/errors"

	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ReconcileVarnishCluster) reconcileClusterRole(ctx context.Context, instance *vcapi.VarnishCluster) error {
	role := &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   names.ClusterRole(instance.Name, instance.Namespace),
			Labels: labels.CombinedComponentLabels(instance, vcapi.VarnishComponentClusterRole),
			Annotations: map[string]string{
				annotationVarnishClusterNamespace: instance.Namespace,
				annotationVarnishClusterName:      instance.Name,
			},
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"list", "watch", "get", "update"},
			},
			{
				APIGroups: []string{"caching.ibm.com"},
				Resources: []string{"varnishclusters"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets", "configmaps"},
				Verbs:     []string{"list", "get", "watch"},
			},
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentClusterRole)
	logr = logr.With(logger.FieldComponentName, role.Name)

	found := &rbac.ClusterRole{}
	err := r.Get(context.TODO(), types.NamespacedName{Name: role.Name}, found)
	// If the role does not exist, create it
	// Else if there was a problem doing the GET, just return
	// Else if the role exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating ClusterRole", "new", role)
		if err = r.Create(ctx, role); err != nil {
			return errors.Wrap(err, "Unable to create ClusterRole")
		}
	} else if err != nil {
		return errors.Wrap(err, "Could not Get ClusterRole")
	} else if !compare.EqualClusterRole(found, role) {
		logr.Infoc("Updating ClusterRole", "diff", compare.DiffClusterRole(found, role))
		found.Rules = role.Rules
		found.Labels = role.Labels
		found.Annotations = role.Annotations
		if err = r.Update(ctx, found); err != nil {
			return errors.Wrap(err, "Could not Update ClusterRole")
		}
	} else {
		logr.Debugw("No updates for ClusterRole")
	}
	return nil
}
