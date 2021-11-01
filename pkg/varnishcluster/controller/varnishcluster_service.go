package controller

import (
	"context"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishCluster) reconcileServiceNoCache(ctx context.Context, instance, instanceStatus *vcapi.VarnishCluster) (map[string]string, error) {
	selector := make(map[string]string, len(instance.Spec.Backend.Selector))
	for k, v := range instance.Spec.Backend.Selector {
		selector[k] = v
	}
	selectorLabels := vclabels.ComponentLabels(instance, vcapi.VarnishComponentNoCacheService)
	inheritedLabels := vclabels.InheritLabels(instance)
	svcLabels := make(map[string]string, len(selectorLabels)+len(inheritedLabels))
	for k, v := range inheritedLabels {
		svcLabels[k] = v
	}
	for k, v := range selectorLabels {
		svcLabels[k] = v
	}

	serviceNoCache := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      names.NoCacheService(instance.Name),
			Namespace: instance.Namespace,
			Labels:    svcLabels,
		},
		Spec: v1.ServiceSpec{
			Selector: selector,
			Ports: []v1.ServicePort{
				{
					Name:       "backend",
					Protocol:   v1.ProtocolTCP,
					Port:       *instance.Spec.Service.Port,
					TargetPort: *instance.Spec.Backend.Port,
				},
			},
			SessionAffinity: v1.ServiceAffinityNone,
			Type:            v1.ServiceTypeClusterIP,
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentNoCacheService)
	logr = logr.With(logger.FieldComponent, serviceNoCache.Name)
	ctx = logger.ToContext(ctx, logr)

	if err := r.reconcileServiceGeneric(ctx, instance, serviceNoCache); err != nil {
		return selectorLabels, err
	}
	return selectorLabels, nil
}

func (r *ReconcileVarnishCluster) reconcileService(ctx context.Context, instance, instanceStatus *vcapi.VarnishCluster, varnishSelector map[string]string) error {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        instance.Name,
			Namespace:   instance.Namespace,
			Labels:      vclabels.CombinedComponentLabels(instance, vcapi.VarnishComponentCacheService),
			Annotations: instance.Spec.Service.Annotations,
		},
	}

	logr := logger.FromContext(ctx).With(logger.FieldComponent, vcapi.VarnishComponentCacheService)
	logr = logr.With(logger.FieldComponent, service.Name)
	ctx = logger.ToContext(ctx, logr)

	service.Spec = v1.ServiceSpec{
		Selector: vclabels.CombinedComponentLabels(instance, vcapi.VarnishComponentVarnish),
		Ports: []v1.ServicePort{
			{
				Name:       vcapi.VarnishPortName,
				Protocol:   v1.ProtocolTCP,
				Port:       *instance.Spec.Service.Port,
				TargetPort: intstr.FromString(vcapi.VarnishPortName),
			},
			{
				Name:       vcapi.VarnishMetricsPortName,
				Protocol:   v1.ProtocolTCP,
				Port:       *instance.Spec.Service.MetricsPort,
				TargetPort: intstr.FromString(vcapi.VarnishMetricsPortName),
			},
			{
				Name:       vcapi.VarnishControllerMetricsPortName,
				Protocol:   v1.ProtocolTCP,
				Port:       vcapi.VarnishControllerMetricsPort,
				TargetPort: intstr.FromString(vcapi.VarnishControllerMetricsPortName),
			},
		},
		SessionAffinity: v1.ServiceAffinityNone,
		Type:            instance.Spec.Service.Type,
	}

	if err := r.reconcileServiceGeneric(ctx, instance, service); err != nil {
		return err
	}

	instanceStatus.Status.VarnishPodsSelector = labels.FormatLabels(service.Spec.Selector)
	return nil
}

func (r *ReconcileVarnishCluster) reconcileServiceGeneric(ctx context.Context, instance *vcapi.VarnishCluster, desired *v1.Service) error {
	logr := logger.FromContext(ctx)

	// Set controller reference for desired object
	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return errors.Wrap(err, "Cannot set controller reference for desired")
	}
	r.scheme.Default(desired)

	found := &v1.Service{}

	err := r.Get(ctx, types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, found)
	// If the desired does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the desired exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Service", "new", desired)
		if err = r.Create(ctx, desired); err != nil {
			return errors.Wrap(err, "Unable to create service")
		}
	} else if err != nil {
		return errors.Wrap(err, "Could not Get desired")
	} else {
		// ClusterIP is immutable once created, so always enforce the same as existing
		desired.Spec.ClusterIP = found.Spec.ClusterIP
		desired.Spec.ClusterIPs = found.Spec.ClusterIPs
		desired.Spec.IPFamilies = found.Spec.IPFamilies
		desired.Spec.IPFamilyPolicy = found.Spec.IPFamilyPolicy
		desired.Spec.InternalTrafficPolicy = found.Spec.InternalTrafficPolicy
		// use nodePort from the spec or the one allocated by Kubernetes
		if desired.Spec.Type == v1.ServiceTypeLoadBalancer || desired.Spec.Type == v1.ServiceTypeNodePort {
			inheritNodePorts(desired.Spec.Ports, found.Spec.Ports)
		}

		if !compare.EqualService(found, desired) {
			logr.Infoc("Updating Service", "diff", compare.DiffService(found, desired))
			found.Spec = desired.Spec
			found.Labels = desired.Labels
			found.Annotations = desired.Annotations
			if err = r.Update(ctx, found); err != nil {
				return errors.Wrap(err, "Unable to update desired")
			}
		} else {
			logr.Debugw("No updates for Service")
		}
	}

	return nil
}

func inheritNodePorts(to []v1.ServicePort, from []v1.ServicePort) {
	for toIndex, toPort := range to {
		if toPort.NodePort != 0 { //the node port is set by the user
			continue
		}

		for fromIndex, fromPort := range from {
			// set the node port allocated by kubernetes
			if fromPort.Port == toPort.Port {
				to[toIndex].NodePort = from[fromIndex].NodePort
			}
		}
	}
}
