package controller

import (
	"context"
	"fmt"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"
	"icm-varnish-k8s-operator/pkg/varnishservice/logger"
	"io/ioutil"

	"go.uber.org/zap"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	annotationVCLVersion = "VCLVersion"
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

func (r *ReconcileVarnishService) reconcileConfigMap(podsSelector map[string]string, instance, instanceStatus *icmapiv1alpha1.VarnishService) (*v1.ConfigMap, error) {
	logr := logger.With("name", instance.Spec.VCLConfigMap.Name, "namespace", instance.Namespace)

	cm := &v1.ConfigMap{}
	cmLabels := combinedLabels(instance, "vcl-file-configmap")
	err := r.Get(context.TODO(), types.NamespacedName{Name: instance.Spec.VCLConfigMap.Name, Namespace: instance.Namespace}, cm)
	// if the ConfigMap does not exist, create it and set it with the default VCL files
	// Else if there was a problem doing the Get, just return an error
	// Else fill in missing values -- "OwnerReference" or Labels
	// Else do nothing
	if err != nil && kerrors.IsNotFound(err) {
		cm = &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      instance.Spec.VCLConfigMap.Name,
				Labels:    cmLabels,
				Namespace: instance.Namespace,
			},
			Data: map[string]string{
				instance.Spec.VCLConfigMap.DefaultFile:      defaultVCL,
				instance.Spec.VCLConfigMap.BackendsTmplFile: backendsVCLTmpl,
			},
		}
		if err := controllerutil.SetControllerReference(instance, cm, r.scheme); err != nil {
			return nil, logr.RErrorw(err, "could not initialize default ConfigMap")
		}

		logr.Infoc("Creating ConfigMap with default VCL files", "new", cm)
		if err = r.Create(context.TODO(), cm); err != nil {
			return nil, logr.RErrorw(err, "could not create ConfigMap")
		}
	} else if err != nil {
		return nil, logr.RErrorw(err, "could not get current state of ConfigMap")
	} else {
		cmCopy := cm.DeepCopy() //create a copy to check later if the config map changed and needs to be updated
		// TODO: there may be a problem if the configmap is already owned by something else. That will prevent the `Watch` fn (in varnishservice_controller.go#run) from detecting updates to the ConfigMap. It will also cause this code to throw an unhandled error that we may want to handle
		if err = controllerutil.SetControllerReference(instance, cm, r.scheme); err != nil {
			return nil, logr.RErrorw(err, "could not set controller as the OwnerReference for existing ConfigMap")
		}
		// don't trample on any labels created by user
		if cm.Labels == nil {
			cm.Labels = make(map[string]string, len(cmLabels))
		}
		for l, v := range cmLabels {
			cm.Labels[l] = v
		}

		if !compare.EqualConfigMap(cm, cmCopy) {
			logger.Infow("Updating ConfigMap with defaults", "diff", compare.DiffConfigMap(cm, cmCopy))
			if err = r.Update(context.TODO(), cm); err != nil {
				return nil, logger.RErrorw(err, "could not update deployment")
			}
		} else {
			logr.Debugw("No updates for ConfigMap")
		}
	}

	instanceStatus.Status.VCL.ConfigMapVersion = cm.GetResourceVersion()
	if cm.Annotations != nil && cm.Annotations[annotationVCLVersion] != "" {
		v := cm.Annotations[annotationVCLVersion]
		instanceStatus.Status.VCL.Version = &v
	} else {
		instanceStatus.Status.VCL.Version = nil //ensure the status field is empty if the annotation is
	}

	pods := &v1.PodList{}
	selector := labels.SelectorFromSet(podsSelector)
	err = r.List(context.Background(), &client.ListOptions{LabelSelector: selector}, pods)
	if err != nil {
		return nil, logger.RErrorw(err, "can't get list of pods")
	}

	latest, outdated := 0, 0
	for _, item := range pods.Items {
		//do not count pods that are not updated with VCL version. Those are pods that are just created and not fully functional
		if item.Annotations["configMapVersion"] == "" {
			logger.Debugw(fmt.Sprintf("ConfigMapVersion annotation is not present. Skipping the pod."))
		} else if item.Annotations["configMapVersion"] == instance.Status.VCL.ConfigMapVersion {
			latest++
		} else {
			outdated++
		}
	}

	instanceStatus.Status.VCL.Availability = fmt.Sprintf("%d latest / %d outdated", latest, outdated)
	return cm, nil
}
