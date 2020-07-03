package controller

import (
	"context"
	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/labels"
	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/names"
	"github.com/ibm/varnish-operator/pkg/varnishcluster/compare"

	"k8s.io/apimachinery/pkg/util/intstr"

	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishCluster) reconcileHeadlessService(ctx context.Context, instance *vcapi.VarnishCluster) error {
	namespacedName := types.NamespacedName{Namespace: instance.Namespace, Name: names.HeadlessService(instance.Name)}
	logr := logger.FromContext(ctx)

	logr = logr.With(logger.FieldComponent, vcapi.VarnishComponentHeadlessService)
	logr = logr.With(logger.FieldComponentName, namespacedName.Name)

	desired := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Labels:    labels.CombinedComponentLabels(instance, vcapi.VarnishComponentHeadlessService),
			Namespace: namespacedName.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "varnish",
					Protocol:   v1.ProtocolTCP,
					Port:       vcapi.VarnishPort,
					TargetPort: intstr.FromInt(vcapi.VarnishPort),
					NodePort:   0,
				},
			},
			ClusterIP:       v1.ClusterIPNone,
			Type:            v1.ServiceTypeClusterIP,
			SessionAffinity: v1.ServiceAffinityNone,
			Selector:        vclabels.CombinedComponentLabels(instance, vcapi.VarnishComponentVarnish),
		},
	}

	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return errors.Wrap(err, "could not set controller as the OwnerReference for headless service")
	}

	found := &v1.Service{}

	err := r.Get(ctx, namespacedName, found)
	// If the headless service does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the headless service exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating Headless Service", "new", desired)
		if err = r.Create(ctx, desired); err != nil {
			return errors.Wrap(err, "could not create Headless Service")
		}
	} else if err != nil {
		return errors.Wrap(err, "could not get current state of Headless Service")
	} else if !compare.EqualService(found, desired) {
		logr.Infoc("Updating Headless Service", "diff", compare.DiffService(found, desired))
		if err = r.Update(ctx, found); err != nil {
			return errors.Wrap(err, "could not update Headless Service")
		}
	} else {
		logr.Debugw("No updates for Headless Service")
	}

	return nil
}
