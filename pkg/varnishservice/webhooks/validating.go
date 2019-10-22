package webhooks

import (
	"context"
	"fmt"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"net/http"
	"regexp"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	atypes "sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type validationWebhook struct {
	logger  *logger.Logger
	client  client.Client
	decoder atypes.Decoder
}

// podValidator implements inject.Client.
// A client will be automatically injected by Kubebuilder internals.
var _ inject.Client = &validationWebhook{}

var (
	varnishArgsKeyRegexp  = regexp.MustCompile(`-\w`)
	disallowedVarnishArgs = map[string]bool{
		"-a": true,
		"-f": true,
		"-F": true,
		"-n": true,
		"-S": true,
	}
	disallowedVarnishArgsAsString string
)

func init() {
	disallowedVarnishArgsAsArr := make([]string, len(disallowedVarnishArgs))
	i := 0
	for k := range disallowedVarnishArgs {
		disallowedVarnishArgsAsArr[i] = k
		i++
	}
	disallowedVarnishArgsAsString = fmt.Sprintf(`"%s"`, strings.Join(disallowedVarnishArgsAsArr, `", "`))
}

func (w *validationWebhook) InjectClient(c client.Client) error {
	w.client = c
	return nil
}

// podValidator implements inject.Decoder.
// A decoder will be automatically injected by Kubebuilder internals.
var _ inject.Decoder = &validationWebhook{}

func (w *validationWebhook) InjectDecoder(d atypes.Decoder) error {
	w.decoder = d
	return nil
}

// Handle implements admission webhook interface
func (w *validationWebhook) Handle(ctx context.Context, req atypes.Request) atypes.Response {
	vs := &v1alpha1.VarnishService{}
	w.logger.Debugw("Validating webhook called.")

	err := w.decoder.Decode(req, vs)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}

	if resp := validVarnishArgs(vs.Spec.StatefulSet.Container.VarnishArgs, w.logger); !resp.Response.Allowed {
		return resp
	}
	if resp := validPorts(vs.Spec.Service); !resp.Response.Allowed {
		return resp
	}

	return admission.ValidationResponse(true, "")
}

func validVarnishArgs(args []string, logr *logger.Logger) atypes.Response {
	for i := 0; i < len(args); {
		if !varnishArgsKeyRegexp.MatchString(args[i]) {
			return admission.ValidationResponse(false, `varnish args must follow pattern: ["key"[, "value"][,"key"[, "value"]]...] where key follows regexp "-\\w" and value is optional. eg ["-s", "malloc,1024M", "-p", "default_ttl=3600", "-T", "127.0.0.1:6082"]`)
		}
		if _, found := disallowedVarnishArgs[args[i]]; found {
			return admission.ValidationResponse(false, fmt.Sprintf("cannot include args %s", disallowedVarnishArgsAsString))
		}
		i++
		if i < len(args) && !varnishArgsKeyRegexp.MatchString(args[i]) {
			i++
		}
	}
	return admission.ValidationResponse(true, "")
}

func validPorts(service v1alpha1.VarnishServiceService) atypes.Response {
	varnishPortName, varnishExporterPortName := "varnish", "varnishexporter"
	if service.VarnishPort.Name != "" {
		varnishPortName = service.VarnishPort.Name
	}

	if service.VarnishExporterPort.Name != "" {
		varnishExporterPortName = service.VarnishExporterPort.Name
	}

	for idx, port := range service.Ports {
		if port.Name == varnishPortName {
			return admission.ValidationResponse(false, fmt.Sprintf("cannot name port %s in .spec.service.ports[%d] (duplicate of varnishPort)", varnishPortName, idx))
		}
		if port.Name == varnishExporterPortName {
			return admission.ValidationResponse(false, fmt.Sprintf("cannot name port %s in .spec.service.ports[%d] (duplicate of varnishExporterPort)", varnishExporterPortName, idx))
		}
	}
	return admission.ValidationResponse(true, "")
}
