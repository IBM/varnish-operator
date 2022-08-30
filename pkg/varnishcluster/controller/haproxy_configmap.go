package controller

import (
	"bytes"
	"context"
	"text/template"

	"github.com/ibm/varnish-operator/pkg/names"

	"github.com/ibm/varnish-operator/api/v1alpha1"
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishCluster) reconcileHaproxyConfigMap(ctx context.Context, podsSelector map[string]string, instance *vcapi.VarnishCluster, instanceStatus *vcapi.VarnishCluster) error {
	logr := logger.FromContext(ctx).With(logger.FieldComponent, names.HaproxyConfigMap(instance.Name))
	logr = logr.With(logger.FieldComponentName, names.HaproxyConfigMap(instance.Name))

	cm := &v1.ConfigMap{}
	cmLabels := vclabels.CombinedComponentLabels(instance, names.HaproxyConfigMap(instance.Name))
	err := r.Get(ctx, types.NamespacedName{Name: names.HaproxyConfigMap(instance.Name), Namespace: instance.Namespace}, cm)
	if err != nil && kerrors.IsNotFound(err) {
		if err := r.updateHaproxyConfigMap(instance, podsSelector, cm, cmLabels, instanceStatus, logr); err != nil {
			return err
		}
		logr.Infoc("Creating HAProxy ConfigMap", "new", cm)
		if err := r.Create(ctx, cm); err != nil {
			return errors.Wrap(err, "could not create ConfigMap")
		}
	} else if err != nil {
		return errors.Wrap(err, "could not get current state of HAProxy ConfigMap")
	} else {
		cmCopy := cm.DeepCopy() //create a copy to check later if the config map changed and needs to be updated
		if err := r.updateHaproxyConfigMap(instance, podsSelector, cm, cmLabels, instanceStatus, logr); err != nil {
			return err
		}
		if !compare.EqualConfigMap(cm, cmCopy) {
			logr.Infow("Updating HAProxy ConfigMap with defaults", "diff", compare.DiffConfigMap(cm, cmCopy))
			if err = r.Update(ctx, cm); err != nil {
				return errors.Wrap(err, "could not update configmap")
			}
		} else {
			logr.Debugw("No updates for ConfigMap")
		}
	}
	return nil
}

func (r *ReconcileVarnishCluster) templatizeHaproxyConfig(instance *vcapi.VarnishCluster, tmpl string) (string, error) {
	t, err := template.New("haproxy-config").Parse(tmpl)
	if err != nil {
		return "", errors.WithStack(err)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, instance.Spec.HaproxySidecar); err != nil {
		return "", errors.WithStack(err)
	}
	return b.String(), nil
}

func (r *ReconcileVarnishCluster) updateHaproxyConfigMap(instance *vcapi.VarnishCluster, podsSelector map[string]string, cm *v1.ConfigMap, cmLabels map[string]string, instanceStatus *vcapi.VarnishCluster, logr *logger.Logger) error {
	haproxyConfig, err := r.templatizeHaproxyConfig(instance, haproxyConfigTemplate)
	if err != nil {
		return err
	}
	cm.Data = map[string]string{v1alpha1.HaproxyConfigFileName: haproxyConfig}

	// don't trample on any labels created by user
	if cm.Labels == nil {
		cm.Labels = make(map[string]string, len(cmLabels))
	}
	for l, v := range cmLabels {
		cm.Labels[l] = v
	}

	instanceStatus.Status.HAProxy.ConfigMapVersion = cm.GetResourceVersion()
	if cm.Annotations != nil && cm.Annotations[annotationHaproxyConfigVersion] != "" {
		v := cm.Annotations[annotationHaproxyConfigVersion]
		instanceStatus.Status.HAProxy.Version = &v
	} else {
		instanceStatus.Status.HAProxy.Version = nil //ensure the status field is empty if the annotation is
	}

	availabilityString, err := r.availabilityString(podsSelector, annotationHaproxyConfigVersion, instance.Status.HAProxy.ConfigMapVersion, logr)
	if err != nil {
		return err
	}
	instanceStatus.Status.HAProxy.Availability = availabilityString

	cm.ObjectMeta.Name = names.HaproxyConfigMap(instance.Name)
	cm.ObjectMeta.Namespace = instance.Namespace
	if err := controllerutil.SetControllerReference(instance, cm, r.scheme); err != nil {
		return errors.Wrap(err, "could not initialize haproxy ConfigMap")
	}
	return nil
}
