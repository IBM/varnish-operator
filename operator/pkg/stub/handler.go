package stub

import (
	icmv1alpha1 "icm-varnish-k8s-operator/operator/pkg/apis/icm/v1alpha1"

	juju "github.com/juju/errors"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/handler"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	"github.com/operator-framework/operator-sdk/pkg/sdk/types"
	"github.com/sirupsen/logrus"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHandler creates a new instance of Handler for reconciliation
func NewHandler() handler.Handler {
	return &Handler{}
}

// Handler has a Handle function defined for it that determines how reconciliation is done
type Handler struct {
	// Fill me
}

// Handle reconciles the updated custom resource with the current state
func (h *Handler) Handle(ctx types.Context, event types.Event) error {
	switch o := event.Object.(type) {
	case *icmv1alpha1.VarnishService:
		newState := o
		dep := deploymentForVarnishService(newState)

		if err := createIfNotExists(dep); err != nil {
			return juju.Trace(err)
		}

		if err := diffState(dep, newState); err != nil {
			return juju.Trace(err)
		}

	}
	return nil
}

// createIfNotExists creates the deployment if it previously did not exist
func createIfNotExists(dep *appsv1beta2.Deployment) error {
	err := action.Create(dep)
	if err != nil && !errors.IsAlreadyExists(err) {
		return juju.Annotate(err, "Failed to create deployment")
	}
	return nil
}

// TODO: this is a stub function that does nothing right now
func diffState(currDep *appsv1beta2.Deployment, newState *icmv1alpha1.VarnishService) error {
	if err := query.Get(currDep); err != nil {
		return juju.Annotate(err, "could not retrieve current state of VarnishService deployment")
	}
	currSize := *currDep.Spec.Replicas
	desiredSize := newState.Spec.Replicas
	if currSize != desiredSize {
		logrus.WithField("sizeDiff", desiredSize-currSize).Info("state change here")
	}
	return nil
}

func deploymentForVarnishService(cr *icmv1alpha1.VarnishService) *appsv1beta2.Deployment {
	labels := labelsForVarnishCache(cr.Name)
	dep := &appsv1beta2.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1beta2.SchemeGroupVersion.String(),
			Kind:       "deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: &cr.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Image: "registry.ng.bluemix.net/icm-varnish/varnish:1.0.1",
						Name:  "varnish",
						Ports: []v1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "varnish-backend",
						}},
					}},
				},
			},
		},
	}
	addOwnerRefToObject(dep, asOwner(cr))
	return dep
}

// labelsForVarnishCache generates the labels meant to apply to a varnish cache
func labelsForVarnishCache(name string) map[string]string {
	return map[string]string{
		"app":        "varnish-cache",
		"varnish_cr": name,
	}
}

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// asOwner returns an OwnerReference set as the VarnishService CR
func asOwner(cr *icmv1alpha1.VarnishService) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: cr.APIVersion,
		Kind:       cr.Kind,
		Name:       cr.Name,
		UID:        cr.UID,
		Controller: &trueVar,
	}
}
