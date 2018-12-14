package webhooks

import (
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/config"
	"icm-varnish-k8s-operator/pkg/logger"

	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"
)

func InstallWebhooks(mgr manager.Manager) {
	validatingWebhook, err := builder.NewWebhookBuilder().
		Name("validating-webhook.varnish-operator.icm.ibm.com").
		Validating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&v1alpha1.VarnishService{}).
		Handlers(&validationWebhook{}).
		FailurePolicy(admissionregistrationv1beta1.Ignore). //change to Fail for debugging
		Build()

	if err != nil {
		logger.RErrorw(err, "Can't create validating webhook")
		return
	}

	err = validatingWebhook.Validate()
	if err != nil {
		logger.RErrorw(err, "Invalid validating webhook")
		return
	}

	mutatingWebhook, err := builder.NewWebhookBuilder().
		Name("mutating-webhook.varnish-operator.icm.ibm.com").
		Mutating().
		Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
		WithManager(mgr).
		ForType(&v1alpha1.VarnishService{}).
		Handlers(&mutatingWebhook{scheme: mgr.GetScheme()}).
		FailurePolicy(admissionregistrationv1beta1.Ignore). //change to Fail for debugging
		Build()

	if err != nil {
		logger.RErrorw(err, "Can't create mutating webhook")
		return
	}

	err = mutatingWebhook.Validate()
	if err != nil {
		logger.RErrorw(err, "Invalid mutating webhook")
		return
	}

	srv, err := webhook.NewServer(config.GlobalConf.VarnishName+"-webhook-server", mgr, webhook.ServerOptions{
		Port:    9244,
		CertDir: "/tmp/varnish-operator/webhook/certs",
		BootstrapOptions: &webhook.BootstrapOptions{
			ValidatingWebhookConfigName: config.GlobalConf.VarnishName + "-validating-webhook-config",
			MutatingWebhookConfigName:   config.GlobalConf.VarnishName + "-mutating-webhook-config",
			Service: &webhook.Service{
				Namespace: config.GlobalConf.Namespace,
				Name:      config.GlobalConf.VarnishName + "-webhook-service",
				// Selectors should select the pods that runs this webhook server.
				Selectors: map[string]string{
					"admission-controller": "varnish-service-admission-controller",
				},
			},
		},
	})

	if err != nil {
		logger.RErrorw(err, "Can't create validating webhook server")
		return
	}

	_ = srv.Port //make Go not complain about unused variable. Should be removed when enabling webhooks
	// The webhooks are disabled due to a bug in kubernetes 1.11.
	// It leaded to errors like this (shortened):
	// Internal error occurred: jsonpatch replace operation does not apply: doc is missing key: /spec/service/ports/0/targetPort
	// It was caused by the mutating webhook that was setting default values for the service.
	// For now the defaults setting is happening in the reconcile loop until we decide to drop Kubernetes 1.11 support.
	// You can use mutating webhooks for different logic and also safely use the validating webhook functions if you need.
	// To do so, just uncomment the webhooks registering below and make sure you run the server not in Dryrun mode.
	//err = srv.Register(validatingWebhook, mutatingWebhook)
	//if err != nil {
	//	logger.RErrorw(err, "Can't register validating webhook in the admission server")
	//	return
	//}

	//logger.Infow("Admission controller is successfully registered")
}
