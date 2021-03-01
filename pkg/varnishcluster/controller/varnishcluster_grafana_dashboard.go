package controller

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileVarnishCluster) reconcileGrafanaDashboard(ctx context.Context, instance *vcapi.VarnishCluster) error {
	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentGrafanaDashboard)
	logr = logr.With(logger.FieldComponentName, names.GrafanaDashboard(instance.Name))
	ctx = logger.ToContext(ctx, logr)

	err := r.garbageCollectGrafanaDashboards(ctx, instance)
	if err != nil {
		return errors.WithStack(err)
	}

	if instance.Spec.Monitoring.GrafanaDashboard == nil ||
		(instance.Spec.Monitoring.GrafanaDashboard != nil && !instance.Spec.Monitoring.GrafanaDashboard.Enabled) {
		return nil
	}

	if instance.Spec.Monitoring.GrafanaDashboard.Title == "" {
		instance.Spec.Monitoring.GrafanaDashboard.Title = fmt.Sprintf("Varnish (%s)", instance.Name)
	}

	if instance.Spec.Monitoring.GrafanaDashboard.Namespace != "" {
		installationNamespace := instance.Spec.Monitoring.GrafanaDashboard.Namespace
		err := r.Get(ctx, types.NamespacedName{Name: installationNamespace}, &v1.Namespace{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				errMsg := fmt.Sprintf("Can't install Grafana dashboard. Namespace %q doesn't exist", installationNamespace)
				logr.Warn(errMsg)
				r.events.Warning(instance, EventReasonNamespaceNotFound, errMsg)
				return nil
			}
			return errors.WithStack(err)
		}
	}

	grafanaDashboard, err := r.createGrafanaDashboardConfigMap(instance)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := r.applyGrafanaDashboardConfigMap(ctx, instance, grafanaDashboard); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ReconcileVarnishCluster) applyGrafanaDashboardConfigMap(ctx context.Context, instance *vcapi.VarnishCluster, grafanaDashboard *v1.ConfigMap) error {
	logr := logger.FromContext(ctx)
	found := &v1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Namespace: grafanaDashboard.GetNamespace(), Name: names.GrafanaDashboard(instance.Name)}, found)
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Grafana dashboard ConfigMap")
		if err = r.Create(ctx, grafanaDashboard); err != nil {
			return errors.Wrap(err, "Unable to create Grafana dashboard ConfigMap")
		}
	} else if err != nil {
		return errors.Wrap(err, "Could not get Grafana dashboard ConfigMap")
	} else if !compare.EqualConfigMap(found, grafanaDashboard) {
		logr.Infoc("Updating Grafana dashboard ConfigMap", "diff", compare.DiffConfigMap(found, grafanaDashboard))
		found.Data = grafanaDashboard.Data
		found.SetLabels(grafanaDashboard.GetLabels())
		if err = r.Update(ctx, found); err != nil {
			return errors.Wrap(err, "Unable to update Grafana dashboard ConfigMap")
		}
	} else {
		logr.Debugw("No updates for Grafana dashboard ConfigMap")
	}

	return nil
}

// Deletes Grafana dashboards that are no longer needed. For example if the namespace is changed, the ConfigMap from the previous namespace should be deleted.
// Also if the dashboard installation is disabled through the spec we also need to delete the ConfigMap.
func (r *ReconcileVarnishCluster) garbageCollectGrafanaDashboards(ctx context.Context, instance *vcapi.VarnishCluster) error {
	grafanaDashboardSpec := instance.Spec.Monitoring.GrafanaDashboard

	//delete all of them if Grafana dashboard installation is disabled
	if (grafanaDashboardSpec != nil && !grafanaDashboardSpec.Enabled) || grafanaDashboardSpec == nil {
		if err := r.removeAllGrafanaDashboards(ctx, instance); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	grafanaDashboardToBeLeftNamespace := instance.Namespace
	if grafanaDashboardSpec.Namespace != "" {
		grafanaDashboardToBeLeftNamespace = grafanaDashboardSpec.Namespace
	}

	if err := r.removePreviouslyCreatedDashboards(ctx, instance, grafanaDashboardToBeLeftNamespace); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// removePreviousDashboards removes ConfigMaps that are left as a result of config change.
// For example, if the user deployed the dashboards to namespace foo and then decided to move it back to namespace bar, we need to delete the one created in foo.
func (r *ReconcileVarnishCluster) removePreviouslyCreatedDashboards(ctx context.Context, instance *vcapi.VarnishCluster, grafanaDashboardToBeLeftNamespace string) error {
	dashboardsList := &v1.ConfigMapList{}
	err := r.List(ctx, dashboardsList, client.MatchingLabels(labels.CombinedComponentLabels(instance, vcapi.VarnishComponentGrafanaDashboard)))
	if err != nil {
		return errors.WithStack(err)
	}
	for _, item := range dashboardsList.Items {
		if item.GetNamespace() != grafanaDashboardToBeLeftNamespace {
			logger.FromContext(ctx).Infof("Deleting Grafana dashboard ConfigMap %s/%s", item.GetNamespace(), item.GetName())
			err := r.Delete(ctx, &item)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

func (r *ReconcileVarnishCluster) removeAllGrafanaDashboards(ctx context.Context, instance *vcapi.VarnishCluster) error {
	dashboardsList := &v1.ConfigMapList{}
	err := r.List(ctx, dashboardsList, client.MatchingLabels(labels.CombinedComponentLabels(instance, vcapi.VarnishComponentGrafanaDashboard)))
	if err != nil {
		return errors.WithStack(err)
	}
	for _, item := range dashboardsList.Items {
		logger.FromContext(ctx).Infof("Deleting Grafana dashboard ConfigMap %s/%s", item.GetNamespace(), item.GetName())
		err := r.Delete(ctx, &item)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *ReconcileVarnishCluster) createGrafanaDashboardConfigMap(instance *vcapi.VarnishCluster) (*v1.ConfigMap, error) {
	grafanaDashboard := &v1.ConfigMap{}
	grafanaDashboard.SetName(names.GrafanaDashboard(instance.Name))
	grafanaDashboard.SetNamespace(instance.Namespace)
	if instance.Spec.Monitoring.GrafanaDashboard.Namespace != "" {
		grafanaDashboard.SetNamespace(instance.Spec.Monitoring.GrafanaDashboard.Namespace)

	}

	cmLabels := labels.CombinedComponentLabels(instance, vcapi.VarnishComponentGrafanaDashboard)
	additionalLabels := instance.Spec.Monitoring.GrafanaDashboard.Labels
	if len(additionalLabels) > 0 {
		for key, value := range additionalLabels {
			cmLabels[key] = value
		}
	} else {
		cmLabels["grafana_dashboard"] = "1" // Default setting. Set the same value as the Grafana chart uses
	}
	grafanaDashboard.SetLabels(cmLabels)

	// Set reference only if the ConfigMap is installed in the same namespace. It can't be set cross namespace.
	if instance.Spec.Monitoring.GrafanaDashboard.Namespace == "" {
		err := controllerutil.SetControllerReference(instance, grafanaDashboard, r.scheme)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	dashboardData, err := generateGrafanaDashboardData(instance)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	grafanaDashboard.Data = dashboardData
	return grafanaDashboard, nil
}

func generateGrafanaDashboardData(instance *vcapi.VarnishCluster) (map[string]string, error) {
	data := map[string]interface{}{
		"DatasourceName": *instance.Spec.Monitoring.GrafanaDashboard.DatasourceName,
		"Title": instance.Spec.Monitoring.GrafanaDashboard.Title,
	}

	t, err := template.New("GrafanaDashboard").Parse(generateGrafanaDashboardTemplate())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dashboardData := map[string]string{names.GrafanaDashboardFile(instance.Name): b.String()}

	return dashboardData, nil
}
