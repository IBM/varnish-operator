package controller

import (
	"context"
	icmapiv1alpha1 "icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/labels"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishservice/compare"

	"github.com/pkg/errors"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileVarnishService) reconcilePodDisruptionBudget(ctx context.Context, instance *icmapiv1alpha1.VarnishService, podSelector map[string]string) error {
	namespacedName := types.NamespacedName{Namespace: instance.Namespace, Name: instance.Name + "-varnish-pdb"}
	logr := logger.FromContext(ctx)
	logr = logr.With(logger.FieldComponent, icmapiv1alpha1.VarnishComponentPodDisruptionBudget)
	logr = logr.With(logger.FieldComponentName, namespacedName.Name)
	ctx = logger.ToContext(ctx, logr)

	if instance.Spec.PodDisruptionBudget == nil {
		return r.deletePDBIfExists(ctx, namespacedName)
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
		return errors.Wrap(err, "could not set controller as the OwnerReference for poddisruptionbudget")
	}

	found, err := r.getPDB(ctx, namespacedName)
	// If the poddisruptionbudget does not exist, create it
	// Else if there was a problem doing the GET, just return an error
	// Else if the poddisruptionbudget exists, and it is different, update
	// Else no changes, do nothing
	switch {
	case kerrors.IsNotFound(err):
		logr.Infoc("Creating PodDisruptionBudget", "new", desired)
		return r.createPDB(ctx, desired)
	case err != nil:
		return errors.Wrap(err, "could not get current state of poddisruptionbudget")
	case !compare.EqualPodDisruptionBudget(found, desired):
		logr.Infoc("Updating PodDisruptionBudget", "diff", compare.DiffPodDisruptionBudget(found, desired))
		return r.updatePDB(ctx, found, desired)
	default:
		logr.Debugw("No updates for poddisruptionbudget")
		return nil
	}
}

func (r *ReconcileVarnishService) getPDB(ctx context.Context, ns types.NamespacedName) (*policyv1beta1.PodDisruptionBudget, error) {
	pdb := &policyv1beta1.PodDisruptionBudget{}
	err := r.Get(ctx, ns, pdb)
	return pdb, err
}

func (r *ReconcileVarnishService) createPDB(ctx context.Context, pdb *policyv1beta1.PodDisruptionBudget) error {
	err := r.Create(ctx, pdb)
	if err != nil {
		return errors.Wrap(err, "could not create poddisruptionbudget")
	}
	return nil
}

// In Kubernetes version <= 1.14 PodDisruptionBudget updates are forbidden:
// https://github.com/kubernetes/kubernetes/issues/45398
// That's why it needs to be recreated every time the spec changes. Until Kubernetes 1.14 will be deprecated.
func (r *ReconcileVarnishService) updatePDB(ctx context.Context, found, desired *policyv1beta1.PodDisruptionBudget) error {
	err := r.deletePDB(ctx, found)
	if err != nil {
		return err
	}

	return r.createPDB(ctx, desired)
}

func (r *ReconcileVarnishService) deletePDB(ctx context.Context, pdb *policyv1beta1.PodDisruptionBudget) error {
	err := r.Delete(ctx, pdb)
	if err != nil {
		return errors.Wrap(err, "could not delete poddisruptionbudget")
	}
	return nil
}

func (r *ReconcileVarnishService) deletePDBIfExists(ctx context.Context, ns types.NamespacedName) error {
	pdb, err := r.getPDB(ctx, ns)
	if kerrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "could not get current state of poddisruptionbudget")
	}

	l := logger.FromContext(ctx)
	l.Infoc("Deleting existing poddisruptionbudget")
	return r.deletePDB(ctx, pdb)
}
