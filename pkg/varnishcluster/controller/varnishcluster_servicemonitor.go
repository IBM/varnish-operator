package controller

import (
	"context"
	"fmt"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var serviceMonitorGVK = schema.GroupVersionKind{
	Group:   "monitoring.coreos.com",
	Version: "v1",
	Kind:    "ServiceMonitor",
}

var serviceMonitorListGVK = schema.GroupVersionKind{
	Group:   "monitoring.coreos.com",
	Version: "v1",
	Kind:    "ServiceMonitorList",
}

func (r *ReconcileVarnishCluster) reconcileServiceMonitor(ctx context.Context, instance *vcapi.VarnishCluster) error {
	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentPrometheusServiceMonitor)
	logr = logr.With(logger.FieldComponentName, names.ServiceMonitor(instance.Name))
	ctx = logger.ToContext(ctx, logr)

	err := r.cleanupNotNeededServiceMonitors(ctx, instance)
	if err != nil {
		return errors.WithStack(err)
	}

	if !instance.Spec.Monitoring.PrometheusServiceMonitor.Enabled {
		return nil
	}

	if instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace != "" {
		installationNamespace := instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace
		err := r.Get(ctx, types.NamespacedName{Name: installationNamespace}, &v1.Namespace{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				errMsg := fmt.Sprintf("Can't install ServiceMonitor. Namespace %q doesn't exist", installationNamespace)
				logger.FromContext(ctx).Warn(errMsg)
				r.events.Warning(instance, EventReasonNamespaceNotFound, errMsg)
				return nil
			}
			return errors.WithStack(err)
		}
	}

	serviceMonitor, err := r.createServiceMonitorObject(instance)
	if err != nil {
		return err
	}

	if err := r.applyServiceMonitor(ctx, instance, serviceMonitor); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ReconcileVarnishCluster) applyServiceMonitor(ctx context.Context, instance *vcapi.VarnishCluster, serviceMonitor *unstructured.Unstructured) error {
	logr := logger.FromContext(ctx)
	found := &unstructured.Unstructured{}
	found.SetGroupVersionKind(serviceMonitorGVK)
	err := r.Get(ctx, types.NamespacedName{Namespace: serviceMonitor.GetNamespace(), Name: names.ServiceMonitor(instance.Name)}, found)
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating ServiceMonitor", "new", serviceMonitor)
		if err = r.Create(ctx, serviceMonitor); err != nil {
			return errors.Wrap(err, "Unable to create ServiceMonitor")
		}
	} else if _, ok := errors.Cause(err).(*meta.NoKindMatchError); ok {
		r.events.Warning(instance, EventReasonServiceMonitorKindNotFound, "ServiceMonitor can't be installed. Prometheus operator needs to be installed first")
		logr.Warn("ServiceMonitor can't be installed. Prometheus operator needs to be installed first")
		return nil
	} else if err != nil {
		return errors.Wrap(err, "Could not get ServiceMonitor")
	} else if !compare.EqualServiceMonitor(found, serviceMonitor) {
		logr.Infoc("Updating ServiceMonitor", "diff", compare.DiffServiceMonitor(found, serviceMonitor))
		found.Object["spec"] = serviceMonitor.Object["spec"]
		found.SetLabels(serviceMonitor.GetLabels())
		if err = r.Update(ctx, found); err != nil {
			return errors.Wrap(err, "Unable to update ServiceMonitor")
		}
	} else {
		logr.Debugw("No updates for ServiceMonitor")
	}
	return nil
}

// Deletes ServiceMonitors that are no longer needed. For example if the namespace is changed, the resource from the previous namespace should be deleted.
// Also if the ServiceMonitor installation is disabled through the spec we also need to delete the ServiceMonitor.
func (r *ReconcileVarnishCluster) cleanupNotNeededServiceMonitors(ctx context.Context, instance *vcapi.VarnishCluster) error {
	//delete all of them if ServiceMonitor installation is disabled
	if !instance.Spec.Monitoring.PrometheusServiceMonitor.Enabled {
		if err := r.removeAllServiceMonitors(ctx, instance); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	serviceMonitorToBeLeftNamespace := instance.Namespace
	if instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace != "" {
		serviceMonitorToBeLeftNamespace = instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace
	}

	if err := r.removePreviouslyCreatedServiceMonitors(ctx, instance, serviceMonitorToBeLeftNamespace); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ReconcileVarnishCluster) removePreviouslyCreatedServiceMonitors(ctx context.Context, instance *vcapi.VarnishCluster, serviceMonitorToBeLeftNamespace string) error {
	serviceMonitorList := &unstructured.UnstructuredList{}
	serviceMonitorList.SetGroupVersionKind(serviceMonitorListGVK)
	err := r.List(ctx, serviceMonitorList, client.MatchingLabels(labels.CombinedComponentLabels(instance, vcapi.VarnishComponentPrometheusServiceMonitor)))
	if err != nil {
		if _, ok := errors.Cause(err).(*meta.NoKindMatchError); ok { //if no such kind then no such resources as well
			return nil
		}
		return errors.WithStack(err)
	}
	for _, item := range serviceMonitorList.Items {
		if item.GetNamespace() != serviceMonitorToBeLeftNamespace {
			logger.FromContext(ctx).Infof("Deleting ServiceMonitor %s/%s", item.GetNamespace(), item.GetName())
			err := r.Delete(ctx, &item)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

func (r *ReconcileVarnishCluster) removeAllServiceMonitors(ctx context.Context, instance *vcapi.VarnishCluster) error {
	serviceMonitorList := &unstructured.UnstructuredList{}
	serviceMonitorList.SetGroupVersionKind(serviceMonitorListGVK)
	err := r.List(ctx, serviceMonitorList, client.MatchingLabels(labels.CombinedComponentLabels(instance, vcapi.VarnishComponentPrometheusServiceMonitor)))
	if err != nil {
		if _, ok := errors.Cause(err).(*meta.NoKindMatchError); ok { //if no such kind then no such resources as well
			return nil
		}
		return errors.WithStack(err)
	}

	for _, item := range serviceMonitorList.Items {
		logger.FromContext(ctx).Infow("Deleting ServiceMonitor %s/%s", item.GetNamespace(), item.GetName())
		err := r.Delete(ctx, &item)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *ReconcileVarnishCluster) createServiceMonitorObject(instance *vcapi.VarnishCluster) (*unstructured.Unstructured, error) {
	serviceMonitor := &unstructured.Unstructured{}
	serviceMonitor.SetGroupVersionKind(serviceMonitorGVK)
	serviceMonitor.SetName(names.ServiceMonitor(instance.Name))
	serviceMonitor.SetLabels(labels.CombinedComponentLabels(instance, vcapi.VarnishComponentPrometheusServiceMonitor))

	serviceMonitor.SetNamespace(instance.Namespace)
	if instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace != "" {
		serviceMonitor.SetNamespace(instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace)
	}

	additionalLabels := instance.Spec.Monitoring.PrometheusServiceMonitor.Labels
	if len(additionalLabels) > 0 {
		newLabels := serviceMonitor.GetLabels() //inherit existing labels
		for key, value := range additionalLabels {
			newLabels[key] = value
		}
		serviceMonitor.SetLabels(newLabels)
	}

	matchLabels := map[string]interface{}{}
	componentLabels := labels.ComponentLabels(instance, vcapi.VarnishComponentCacheService)
	for key, value := range componentLabels {
		matchLabels[key] = value
	}
	scrapeInterval := instance.Spec.Monitoring.PrometheusServiceMonitor.ScrapeInterval
	serviceMonitor.Object["spec"] = map[string]interface{}{
		"selector": map[string]interface{}{
			"matchLabels": matchLabels,
		},
		"endpoints": []interface{}{
			map[string]interface{}{
				"port":     "metrics",
				"interval": scrapeInterval,
			},
			map[string]interface{}{
				"port":     "ctrl-metrics",
				"interval": scrapeInterval,
			},
		},
		"namespaceSelector": map[string]interface{}{
			"matchNames": []interface{}{instance.Namespace},
		},
	}

	// Set reference only if the servicemonitor is installed in the same namespace. It can't be set cross namespace.
	if instance.Spec.Monitoring.PrometheusServiceMonitor.Namespace == "" {
		err := controllerutil.SetControllerReference(instance, serviceMonitor, r.scheme)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return serviceMonitor, nil
}
