package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/names"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	finalizerClusterRole        = "clusterrole.finalizers.varnishcluster.icm.ibm.com"
	finalizerClusterRoleBinding = "clusterrolebinding.finalizers.varnishcluster.icm.ibm.com"
	finalizerServiceMonitor     = "prometheus-servicemonitor.finalizers.varnishcluster.icm.ibm.com"
	finalizerGrafanaDashboard   = "grafana-dashboard.finalizers.varnishcluster.icm.ibm.com"
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

	if instance.Spec.Monitoring.PrometheusServiceMonitor.Enabled && instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace != "" {
		controllerutil.AddFinalizer(instance, finalizerServiceMonitor)
	}

	grafanaDashboard := instance.Spec.Monitoring.GrafanaDashboard
	if grafanaDashboard != nil && grafanaDashboard.Enabled && grafanaDashboard.Namespace != "" {
		controllerutil.AddFinalizer(instance, finalizerServiceMonitor)
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

	// delete only if the the servicemonitor was installed in a different namespace
	// Otherwise it should be garbage collected by Kubernetes as the owner reference will be set
	if instance.Spec.Monitoring.PrometheusServiceMonitor.Enabled && instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace != "" {
		serviceMonitorName := types.NamespacedName{
			Namespace: instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace,
			Name:      names.ServiceMonitor(instance.Name),
		}
		if err := r.deleteServiceMonitor(ctx, serviceMonitorName); err != nil {
			return err
		}
	}

	if err := r.removeFinalizer(ctx, finalizerServiceMonitor, instance); err != nil {
		return nil
	}

	// delete only if the the servicemonitor was installed in a different namespace
	// Otherwise it should be garbage collected by Kubernetes as the owner reference will be set
	grafanaDashboard := instance.Spec.Monitoring.GrafanaDashboard
	if grafanaDashboard != nil && grafanaDashboard.Enabled && grafanaDashboard.Namespace != "" {
		serviceMonitorName := types.NamespacedName{
			Namespace: grafanaDashboard.Namespace,
			Name:      names.GrafanaDashboard(instance.Name),
		}
		if err := r.deleteGrafanaDashboard(ctx, serviceMonitorName); err != nil {
			return err
		}
	}

	if err := r.removeFinalizer(ctx, finalizerGrafanaDashboard, instance); err != nil {
		return nil
	}

	return nil
}

func (r *ReconcileVarnishCluster) deleteCR(ctx context.Context, ns types.NamespacedName) error {
	cr := &rbac.ClusterRole{}
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
	crb := &rbac.ClusterRoleBinding{}
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

func (r *ReconcileVarnishCluster) deleteServiceMonitor(ctx context.Context, ns types.NamespacedName) error {
	serviceMonitor := &unstructured.Unstructured{}
	serviceMonitor.SetGroupVersionKind(serviceMonitorGVK)
	err := r.Get(ctx, ns, serviceMonitor)
	if kerrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "could not get current state of servicemonitor")
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentPrometheusServiceMonitor)
	logr = logr.With(logger.FieldComponentName, serviceMonitor.GetName())
	logr.Infoc("Deleting existing servicemonitor")
	return r.Delete(ctx, serviceMonitor)
}

func (r *ReconcileVarnishCluster) deleteGrafanaDashboard(ctx context.Context, ns types.NamespacedName) error {
	grafanaDashboard := &v1.ConfigMap{}
	err := r.Get(ctx, ns, grafanaDashboard)
	if kerrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "could not get current state of Grafana dashboard ConfigMap")
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentGrafanaDashboard)
	logr = logr.With(logger.FieldComponentName, grafanaDashboard.GetName())
	logr.Infoc("Deleting existing Grafana dashboard ConfigMap")
	return r.Delete(ctx, grafanaDashboard)
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
