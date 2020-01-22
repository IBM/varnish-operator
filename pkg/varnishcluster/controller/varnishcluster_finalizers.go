package controller

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/names"
	v1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	finalizerClusterRole        = "clusterrole.finalizers.varnishcluster.icm.ibm.com"
	finalizerClusterRoleBinding = "clusterrolebinding.finalizers.varnishcluster.icm.ibm.com"
)

func (r *ReconcileVarnishCluster) reconcileFinalizers(ctx context.Context, instance *icmapiv1alpha1.VarnishCluster) error {
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() { //object is being deleted, don't set finalizers
		return nil
	}

	logr := logger.FromContext(ctx)

	existingFinalizers := instance.Finalizers
	desiredFinalizers := []string{finalizerClusterRole, finalizerClusterRoleBinding}
	for _, desiredFinalizer := range desiredFinalizers {
		controllerutil.AddFinalizer(instance, desiredFinalizer)
	}

	if !cmp.Equal(instance.Finalizers, existingFinalizers) {
		logr.Infow("Updating finalizers", "diff", cmp.Diff(existingFinalizers, instance.Finalizers))
		err := r.Update(ctx, instance)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *ReconcileVarnishCluster) finalizerCleanUp(ctx context.Context, instance *icmapiv1alpha1.VarnishCluster) error {
	if err := r.deleteCR(ctx, types.NamespacedName{Name: names.ClusterRole(instance.Name, instance.Namespace)}); err != nil {
		return errors.WithStack(err)
	}

	if err := r.removeFinalizer(ctx, finalizerClusterRole, instance); err != nil {
		return err
	}

	if err := r.deleteCRB(ctx, types.NamespacedName{Name: names.ClusterRoleBinding(instance.Name, instance.Namespace)}); err != nil {
		return err
	}

	if err := r.removeFinalizer(ctx, finalizerClusterRoleBinding, instance); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileVarnishCluster) deleteCR(ctx context.Context, ns types.NamespacedName) error {
	cr := &v1.ClusterRole{}
	err := r.Get(ctx, ns, cr)
	if kerrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "could not get current state of clusterrole")
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentClusterRole)
	logr = logr.With(logger.FieldComponentName, cr.Name)
	logr.Infoc("Deleting existing clusterrole")
	return r.Delete(ctx, cr)
}

func (r *ReconcileVarnishCluster) deleteCRB(ctx context.Context, ns types.NamespacedName) error {
	crb := &v1.ClusterRoleBinding{}
	err := r.Get(ctx, ns, crb)
	if kerrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "could not get current state of clusterrolebinding")
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentClusterRoleBinding)
	logr = logr.With(logger.FieldComponentName, crb.Name)
	logr.Infoc("Deleting existing clusterrolebinding")
	return r.Delete(ctx, crb)
}

func (r *ReconcileVarnishCluster) removeFinalizer(ctx context.Context, finalizerName string, instance *icmapiv1alpha1.VarnishCluster) error {
	if containsString(instance.Finalizers, finalizerName) {
		controllerutil.RemoveFinalizer(instance, finalizerName)
		if err := r.Update(ctx, instance); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
