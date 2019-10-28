// The app watches kubernetes endpoints api for changes in deployment
// and then re-writes varnish vcl file with any new/removed backends.

package main

import (
	"flag"
	"icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/kwatcher/config"
	"icm-varnish-k8s-operator/pkg/kwatcher/controller"
	"icm-varnish-k8s-operator/pkg/logger"
	"log"

	"github.com/go-logr/zapr"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"go.uber.org/zap"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	Version = "undefined" //will be overwritten by the correct version during docker build
)

func main() {
	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.Parse()

	kwatcherConfig, err := config.Load()
	if err != nil {
		log.Fatalf("could not load kwatcher config: %v", err)
	}

	logr := logger.NewLogger(kwatcherConfig.LogFormat, kwatcherConfig.LogLevel)
	logr = logr.With(logger.FieldKwatcherVersion, Version)

	logr.Infof("Version: %s", Version)
	logr.Infof("Log level: %s", kwatcherConfig.LogLevel.String())

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
		Namespace: kwatcherConfig.Namespace,
		Scheme:    scheme,
	})

	if err != nil {
		logr.With(zap.Error(err)).Fatalf("could not initialize manager")
	}

	logr.Infow("Registering Components")

	// Setup controller
	if err = controller.SetupVarnishReconciler(mgr, kwatcherConfig, logr); err != nil {
		logr.With(zap.Error(err)).Fatalw("could not setup controller")
	}

	logr.Infow("Starting Varnish Watcher")
	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		logr.With(err).Fatalf("Failed to start manager")
	}
}
