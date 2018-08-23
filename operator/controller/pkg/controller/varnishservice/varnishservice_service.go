package varnishservice

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/operator/controller/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/operator/controller/pkg/logger"
	"reflect"

	"github.com/imdario/mergo"
	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileBackendService(instance *icmapiv1alpha1.VarnishService, applicationPort *v1.ServicePort) error {
	backendService := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-backend",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"component": "backend",
			},
		},
		Spec: v1.ServiceSpec{
			Selector: instance.Spec.Service.Selector,
			Ports:    []v1.ServicePort{*applicationPort},
		},
	}
	for k, v := range instance.Spec.Service.Selector {
		backendService.Spec.Selector[k] = v
	}

	if err := r.reconcileServiceGeneric(instance, &backendService); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileVarnishService) reconcileFrontendService(instance *icmapiv1alpha1.VarnishService, applicationPort *v1.ServicePort, varnishSelector *map[string]string) error {
	frontendService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-frontend",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"component": "frontend",
			},
		},
	}
	instance.Spec.Service.DeepCopyInto(&frontendService.Spec)

	frontendService.Spec.Ports = []v1.ServicePort{
		{
			Name:       "http",
			Port:       applicationPort.Port,
			Protocol:   v1.ProtocolTCP,
			TargetPort: intstr.FromInt(r.globalConf.VarnishTargetPort),
		},
	}

	frontendService.Spec.Selector = *varnishSelector

	if err := r.reconcileServiceGeneric(instance, frontendService); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileVarnishService) reconcileServiceGeneric(instance *icmapiv1alpha1.VarnishService, service *v1.Service) error {
	// Set controller reference for service object
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return logger.RError(err, "Cannot set controller reference for service", "namespace", service.Namespace, "name", service.Name)
	}

	found := &v1.Service{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	// If the service does not exist, create it
	if err != nil && kerrors.IsNotFound(err) {
		logger.Info("Creating service", "config", service)
		// logger.Info("Creating service", "namespace", service.Namespace, "name", service.Name)
		if err = r.Create(context.TODO(), service); err != nil {
			return logger.RError(err, "Unable to create service")
		}
		// If there was a problem doing the GET, just return
	} else if err != nil {
		return logger.RError(err, "Could not Get service")
		// If the service exists, and it is different, update
	} else if !reflect.DeepEqual(service.Spec, found.Spec) {
		mergo.Merge(found, service, mergo.WithOverride)
		logger.Info("Updating service", "config", found)
		// logger.Info("Updating service", "namespace", service.Namespace, "name", service.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return logger.RError(err, "Unable to update service")
		}
	}
	// If no changes, do nothing
	logger.Info("No updates for Service", "name", service.Name)
	return nil
}
