package varnishservice

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/compare"
	"icm-varnish-k8s-operator/pkg/logger"
	"io/ioutil"

	"go.uber.org/zap"

	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var defaultVCL, backendsVCLTmpl string

func init() {
	readVCL := func(file string) string {
		bs, err := ioutil.ReadFile("config/vcl/" + file)
		if err != nil {
			// TODO: use new logger implementation in other branch for this
			logger.Panicw("could not find file for ConfigMap", "filename", file, zap.Error(err))
		}
		return string(bs)
	}
	defaultVCL = readVCL("default.vcl")
	backendsVCLTmpl = readVCL("backends.vcl.tmpl")
}

func (r *ReconcileVarnishService) reconcileConfigMap(instance *icmapiv1alpha1.VarnishService) (map[string]string, error) {
	logr := logger.With("name", instance.Spec.VCLConfigMap.Name, "namespace", instance.Namespace)

	found := &v1.ConfigMap{}

	selectorLabels := generateLabels(instance, "default-file-configmap")
	inheritedLabels := inheritLabels(instance)
	labels := make(map[string]string, len(selectorLabels)+len(inheritedLabels))
	for k, v := range inheritedLabels {
		labels[k] = v
	}
	for k, v := range selectorLabels {
		labels[k] = v
	}

	err := r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.VCLConfigMap.Name, Namespace: instance.Namespace}, found)
	// if the ConfigMap does not exist, create it and set it with the default VCL files
	// Else if there was a problem doing the Get, just return an error
	// Else fill in missing values -- "OwnerReference" or Labels
	// Else do nothing
	if err != nil && kerrors.IsNotFound(err) {
		desired := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.Spec.VCLConfigMap.Name,
				Labels:    labels,
				Namespace: instance.Namespace,
			},
			Data: map[string]string{
				instance.Spec.VCLConfigMap.DefaultFile:      defaultVCL,
				instance.Spec.VCLConfigMap.BackendsTmplFile: backendsVCLTmpl,
			},
		}
		if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
			return selectorLabels, logr.RErrorw(err, "could not initialize default ConfigMap")
		}

		logr.Info("Creating ConfigMap with default VCL files", "new", desired)
		if err = r.Create(context.TODO(), desired); err != nil {
			return selectorLabels, logr.RErrorw(err, "could not create ConfigMap")
		}
	} else if err != nil {
		return selectorLabels, logr.RErrorw(err, "could not get current state of ConfigMap")
	} else {
		foundCopy := found.DeepCopy()
		// TODO: there may be a problem if the configmap is already owned by something else. That will prevent the `Watch` fn (in varnishservice_controller.go#run) from detecting updates to the ConfigMap. It will also cause this code to throw an unhandled error that we may want to handle
		if err = controllerutil.SetControllerReference(instance, foundCopy, r.scheme); err != nil {
			return selectorLabels, logr.RErrorw(err, "could not set controller as the OwnerReference for existing ConfigMap")
		}
		// don't trample on any labels created by user
		for l, v := range labels {
			foundCopy.Labels[l] = v
		}

		if !compare.EqualConfigMap(found, foundCopy) {
			logger.Infow("Updating ConfigMap with defaults", "diff", compare.DiffConfigMap(found, foundCopy))
			if err = r.Update(context.TODO(), foundCopy); err != nil {
				return selectorLabels, logger.RErrorw(err, "could not update deployment")
			}
		} else {
			logr.Debugw("No updates for ConfigMap")
		}
	}
	return selectorLabels, nil
}
