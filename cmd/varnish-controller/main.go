// The app watches kubernetes endpoints api for changes in deployment
// and then re-writes varnish vcl file with any new/removed backends.

package main

import (
	"flag"
	"fmt"
	"icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/config"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/controller"
	"icm-varnish-k8s-operator/pkg/varnishcontroller/varnishadm"
	"log"

	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/go-logr/zapr"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// Version defines varnish-controller version. Will be overwritten by the correct version during docker build
var (
	Version = "undefined"
)

func main() {
	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.Parse()

	varnishControllerConfig, err := config.Load()
	if err != nil {
		log.Fatalf("could not load varnish-controller config: %v", err)
	}

	logr := logger.NewLogger(varnishControllerConfig.LogFormat, varnishControllerConfig.LogLevel)
	logr = logr.With(logger.FieldVarnishControllerVersion, Version)

	logr.Infof("Version: %s", Version)
	logr.Infof("Log level: %s", varnishControllerConfig.LogLevel.String())

	ctrl.SetLogger(zapr.NewLogger(logr.Desugar())) //set logger for controller-runtime to see internal library logs

	scheme := runtime.NewScheme()
	err = clientgoscheme.AddToScheme(scheme)
	if err != nil {
		logr.With(zap.Error(err)).Fatalf("unable to set up standard schemes config")
	}

	err = v1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
	if err != nil {
		logr.With(zap.Error(err)).Fatalf("unable to set up varnish operator schemes config")
	}

	clientConfig, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalf("could not load rest client config. Error: %s", err)
	}

	mgr, err := ctrl.NewManager(clientConfig, ctrl.Options{
		Namespace:              varnishControllerConfig.Namespace,
		Scheme:                 scheme,
		HealthProbeBindAddress: fmt.Sprintf(":%d", v1alpha1.HealthCheckPort),
	})

	if err != nil {
		logr.With(zap.Error(err)).Fatalf("could not initialize manager")
	}

	err = mgr.AddReadyzCheck("ping", healthz.Ping)
	if err != nil {
		logr.With(zap.Error(err)).Fatal("unable to set readiness check")
	}

	logr.Infow("Registering Components")

	// Setup controller
	varnishAdm := varnishadm.NewVarnishAdministartor(varnishControllerConfig.VarnishPingTimeout,
		varnishControllerConfig.VarnishPingDelay,
		config.VCLConfigDir,
		varnishControllerConfig.VarnishAdmArgs)

	if err = controller.SetupVarnishReconciler(mgr, varnishControllerConfig, varnishAdm, logr); err != nil {
		logr.With(zap.Error(err)).Fatalw("could not setup controller")
	}
	logr.Infow("Looking up for a Varnish service")
	if err = varnishAdm.Ping(); err != nil {
		logr.With(err).Fatalf("Varnish is unreachable")
	}
	logr.Infow("Starting Varnish Controller")
	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		logr.With(err).Fatalf("Failed to start manager")
	}
}
