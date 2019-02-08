package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"

	policyv1beta1 "k8s.io/api/policy/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcilePodDisruptionBudget(instance *icmapiv1alpha1.VarnishService, podSelector map[string]string) error {
	namespacedName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-varnish-pdb"}
	logr := r.logger.With("name", namespacedName.Name, "namespace", namespacedName.Namespace)

	// if not specified, do not create one
	if instance.Spec.PodDisruptionBudget == nil {
		logr.Debugw("No poddisruptionbudget specified")
		return nil
	}

	var pdbs policyv1beta1.PodDisruptionBudgetSpec
	instance.Spec.PodDisruptionBudget.DeepCopyInto(&pdbs)
	pdbs.Selector = &metav1.LabelSelector{
		MatchLabels: podSelector,
	}

	desired := &policyv1beta1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Labels:    labels.CombinedComponentLabels(instance, icmapiv1alpha1.VarnishComponentPodDisruptionBudget),
			Namespace: namespacedName.Namespace,
		},
		Spec: pdbs,
	}

	if err := controllerutil.SetControllerReference(instance, desired, r.scheme); err != nil {
		return logr.RErrorw(err, "could not set controller as the OwnerReference for poddisruptionbudget")
	}

	found := &policyv1beta1.PodDisruptionBudget{}

	err := r.Get(context.TODO(), namespacedName, found)
	// If the poddisruptionbudget does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the poddisruptionbudget exists, and it is different, update
	// Else no changes, do nothing
	if err != nil && kerrors.IsNotFound(err) {
		logr.Infoc("Creating PodDisruptionBudget", "new", desired)
		if err = r.Create(context.TODO(), desired); err != nil {
			return logr.RErrorw(err, "could not create poddisruptionbudget")
		}
	} else if err != nil {
		return logr.RErrorw(err, "could not get current state of poddisruptionbudget")
	} else if !compare.EqualPodDisruptionBudget(found, desired) {
		logr.Infoc("Updating PodDisruptionBudget", "diff", compare.DiffPodDisruptionBudget(found, desired))
		if err = r.Update(context.TODO(), found); err != nil {
			return logr.RErrorw(err, "could not update poddisruptionbudget")
		}
	} else {
		logr.Debugw("No updates for poddisruptionbudget")
	}
	return nil
}
