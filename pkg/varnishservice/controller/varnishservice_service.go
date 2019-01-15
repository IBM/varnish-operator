package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"
	"icm-varnish-k8s-operator/pkg/varnishservice/config"
	"icm-varnish-k8s-operator/pkg/varnishservice/logger"
	"strconv"

	"k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcileNoCachedService(instance, instanceStatus *icmapiv1alpha1.VarnishService, applicationPort *v1.ServicePort) (map[string]string, error) {
	selector := make(map[string]string, len(instance.Spec.Service.Selector))
	for k, v := range instance.Spec.Service.Selector {
		selector[k] = v
	}
	selectorLabels := generateLabels(instance, "nocached-service")
	inheritedLabels := inheritLabels(instance)
	labels := make(map[string]string, len(selectorLabels)+len(inheritedLabels))
	for k, v := range inheritedLabels {
		labels[k] = v
	}
	for k, v := range selectorLabels {
		labels[k] = v
	}
	noCachedService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-varnish-nocached",
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Selector: selector,
			Ports:    []v1.ServicePort{*applicationPort},
		},
	}

	if err := r.reconcileServiceGeneric(instance, &instanceStatus.Status.Service.NoCached, noCachedService); err != nil {
		return selectorLabels, err
	}
	return selectorLabels, nil
}

func (r *ReconcileVarnishService) reconcileCachedService(instance, instanceStatus *icmapiv1alpha1.VarnishService, applicationPort *v1.ServicePort, varnishSelector map[string]string) error {
	cachedService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-varnish-cached",
			Namespace: instance.Namespace,
			Labels:    combinedLabels(instance, "cached-service"),
		},
	}

	prometheusAnnotations := map[string]string{
		"prometheus.io/scrape": "true",
		"prometheus.io/port":   strconv.FormatInt(int64(config.GlobalConf.VarnishExporterPort), 10),
	}

	if config.GlobalConf.PrometheusAnnotations {
		cachedService.Annotations = prometheusAnnotations
	}

	instance.Spec.Service.DeepCopyInto(&cachedService.Spec)

	cachedService.Spec.Ports = []v1.ServicePort{
		{
			Name:       "http",
			Port:       applicationPort.Port,
			Protocol:   v1.ProtocolTCP,
			TargetPort: intstr.FromInt(config.GlobalConf.VarnishTargetPort),
		},
		{
			Name:     "exporter",
			Port:     config.GlobalConf.VarnishExporterPort,
			Protocol: v1.ProtocolTCP,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: config.GlobalConf.VarnishExporterTargetPort,
			},
		},
	}

	cachedService.Spec.Selector = varnishSelector

	if err := r.reconcileServiceGeneric(instance, &instanceStatus.Status.Service.Cached, cachedService); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileVarnishService) reconcileServiceGeneric(instance *icmapiv1alpha1.VarnishService, instanceServiceStatus *icmapiv1alpha1.VarnishServiceSingleServiceStatus, desired *v1.Service) error {
	logr := logger.With("name", desired.Name, "namespace", desired.Namespace)

	// Set controller reference for desired object
	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return logr.RErrorw(err, "Cannot set controller reference for desired")
	}
	r.scheme.Default(desired)

	found := &v1.Service{}

	err := r.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, found)
	// If the desired does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the desired exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Service", "new", desired)
		if err = r.Create(context.TODO(), desired); err != nil {
			return logr.RErrorw(err, "Unable to create service")
		}
	} else if err != nil {
		return logr.RErrorw(err, "Could not Get desired")
	} else {
		// ClusterIP is immutable once created, so always enforce the same as existing
		desired.Spec.ClusterIP = found.Spec.ClusterIP
		if !compare.EqualService(found, desired) {
			logr.Infoc("Updating Service", "diff", compare.DiffService(found, desired))
			found.Spec = desired.Spec
			found.Labels = desired.Labels
			if err = r.Update(context.TODO(), found); err != nil {
				return logr.RErrorw(err, "Unable to update desired")
			}
		} else {
			logr.Debugw("No updates for Service")
		}
	}

	*instanceServiceStatus = icmapiv1alpha1.VarnishServiceSingleServiceStatus{
		ServiceStatus: found.Status,
		Name:          found.Name,
		IP:            found.Spec.ClusterIP,
	}
	return nil
}
