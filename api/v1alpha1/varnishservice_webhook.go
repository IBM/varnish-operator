package v1alpha1

import (
	"fmt"
	"icm-varnish-k8s-operator/pkg/logger"
	"regexp"
	"strings"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var webhookLogger = &logger.Logger{SugaredLogger: zap.NewNop().Sugar()}

func SetWebhookLogger(l *logger.Logger) {
	webhookLogger = l
}

var (
	varnishArgsKeyRegexp  = regexp.MustCompile(`^-\w$`)
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

func (in *VarnishService) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-icm-ibm-com-v1alpha1-varnishservice,mutating=true,failurePolicy=fail,groups=icm.ibm.com,resources=varnishservices,verbs=create;update,versions=v1alpha1,name=mvarnishservice.kb.io

var _ webhook.Defaulter = &VarnishService{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (in *VarnishService) Default() {
	// defaulting logic goes here.
}

// note: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-icm-ibm-com-v1alpha1-varnishservice,mutating=false,failurePolicy=fail,groups=icm.ibm.com,resources=varnishservices,versions=v1alpha1,name=vvarnishservice.kb.io

var _ webhook.Validator = &VarnishService{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (in *VarnishService) ValidateCreate() error {
	logr := webhookLogger.With(logger.FieldComponent, VarnishComponentValidatingWebhook)
	logr = logr.With(logger.FieldNamespace, in.Namespace)
	logr = logr.With(logger.FieldVarnishService, in.Name)

	logr.Debug("Validating webhook has been called on create request")
	if err := validVarnishArgs(in.Spec.StatefulSet.Container.VarnishArgs); err != nil {
		return err
	}
	if err := validPorts(in.Spec.Service); err != nil {
		return err
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (in *VarnishService) ValidateUpdate(old runtime.Object) error {
	logr := webhookLogger.With(logger.FieldComponent, VarnishComponentValidatingWebhook)
	logr = logr.With(logger.FieldNamespace, in.Namespace)
	logr = logr.With(logger.FieldVarnishService, in.Name)

	logr.Debug("Validating webhook has been called on update request")
	if err := validVarnishArgs(in.Spec.StatefulSet.Container.VarnishArgs); err != nil {
		return err
	}
	if err := validPorts(in.Spec.Service); err != nil {
		return err
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (in *VarnishService) ValidateDelete() error {
	logr := webhookLogger.With(logger.FieldComponent, VarnishComponentValidatingWebhook)
	logr = logr.With(logger.FieldNamespace, in.Namespace)
	logr = logr.With(logger.FieldVarnishService, in.Name)

	logr.Debug("Validating webhook has been called on delete request")
	return nil
}

func validVarnishArgs(args []string) error {
	for i := 0; i < len(args); {
		if !varnishArgsKeyRegexp.MatchString(args[i]) {
			return errors.Errorf(
				`varnish args must follow pattern: ["key"[, "value"][,"key"[, "value"]]...] where key follows regexp "%s" and value is optional. eg ["-s", "malloc,1024M", "-p", "default_ttl=3600", "-T", "127.0.0.1:6082"]`,
				varnishArgsKeyRegexp.String(),
			)
		}
		if _, found := disallowedVarnishArgs[args[i]]; found {
			return errors.Errorf("cannot include args %s", disallowedVarnishArgsAsString)
		}
		i++
		if i < len(args) && !varnishArgsKeyRegexp.MatchString(args[i]) {
			i++
		}
	}
	return nil
}

func validPorts(service VarnishServiceService) error {
	varnishPortName, varnishExporterPortName := "varnish", "varnishexporter"
	if service.VarnishPort.Name != "" {
		varnishPortName = service.VarnishPort.Name
	}

	if service.VarnishExporterPort.Name != "" {
		varnishExporterPortName = service.VarnishExporterPort.Name
	}

	for idx, port := range service.Ports {
		if port.Name == varnishPortName {
			return errors.Errorf("cannot name port %s in .spec.service.ports[%d] (duplicate of varnishPort)", varnishPortName, idx)
		}
		if port.Name == varnishExporterPortName {
			return errors.Errorf("cannot name port %s in .spec.service.ports[%d] (duplicate of varnishExporterPort)", varnishExporterPortName, idx)
		}
	}
	return nil
}
