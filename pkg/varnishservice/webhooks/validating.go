package webhooks

import (
	"context"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/varnishservice/logger"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	atypes "sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type validationWebhook struct {
	client  client.Client
	decoder atypes.Decoder
}

// podValidator implements inject.Client.
// A client will be automatically injected by Kubebuilder internals.
var _ inject.Client = &validationWebhook{}

func (v *validationWebhook) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

// podValidator implements inject.Decoder.
// A decoder will be automatically injected by Kubebuilder internals.
var _ inject.Decoder = &validationWebhook{}

func (v *validationWebhook) InjectDecoder(d atypes.Decoder) error {
	v.decoder = d
	return nil
}

// Handle implements admission webhook interface
func (v *validationWebhook) Handle(ctx context.Context, req atypes.Request) atypes.Response {
	vs := &v1alpha1.VarnishService{}
	logger.Debugw("Validating webhook called.")

	err := v.decoder.Decode(req, vs)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	// TODO validation logic goes here
	// example: check if number of replicas are less than 100
	//if vs.Spec.Deployment.Replicas != nil && *vs.Spec.Deployment.Replicas > 100 {
	//	return admission.ValidationResponse(false, "Should be less than 100 replicas")
	//}

	return admission.ValidationResponse(true, "")
}
