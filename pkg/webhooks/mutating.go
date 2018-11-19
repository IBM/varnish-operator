package webhooks

import (
	"context"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	atypes "sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type mutatingWebhook struct {
	scheme  *runtime.Scheme
	client  client.Client
	decoder atypes.Decoder
}

// podValidator implements inject.Client.
// A client will be automatically injected by Kubebuilder internals.
var _ inject.Client = &mutatingWebhook{}

func (v *mutatingWebhook) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

// podValidator implements inject.Decoder.
// A decoder will be automatically injected by Kubebuilder internals.
var _ inject.Decoder = &mutatingWebhook{}

func (v *mutatingWebhook) InjectDecoder(d atypes.Decoder) error {
	v.decoder = d
	return nil
}

// Handle implements admission webhook interface
func (v *mutatingWebhook) Handle(ctx context.Context, req atypes.Request) atypes.Response {
	original := &v1alpha1.VarnishService{}
	logger.Debugw("Mutating webhook called.")

	err := v.decoder.Decode(req, original)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	mutated := original.DeepCopy()
	v.scheme.Default(mutated)

	return admission.PatchResponse(original, mutated)
}
