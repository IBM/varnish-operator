// The app watches kubernetes endpoints api for changes in deployment
// and then re-writes varnish vcl file with any new/removed backends.

package main

import (
	"flag"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/kwatcher/config"
	"icm-varnish-k8s-operator/pkg/kwatcher/controller"
	"icm-varnish-k8s-operator/pkg/logger"
	"log"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	kconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func init() {
	flag.Parse()
}

func main() {
	clientConfig, err := kconfig.GetConfig()
	if err != nil {
		log.Fatalf("could not load rest client config. Error: %s", err)
	}

	kwatcherConfig, err := config.Load()
	if err != nil {
		log.Fatalf("could not load kwatcher config: %v", err)
	}

	logr := logger.NewLogger(kwatcherConfig.LogFormat, kwatcherConfig.LogLevel)

	mgr, err := manager.New(clientConfig, manager.Options{Namespace: kwatcherConfig.Namespace})
	if err != nil {
		logr.Fatalf("could not initialize manager", zap.Error(err))
	}

	// Setup Scheme for all resources
	AddToSchemes := runtime.SchemeBuilder{v1alpha1.SchemeBuilder.AddToScheme}
	if err := AddToSchemes.AddToScheme(mgr.GetScheme()); err != nil {
		logr.Fatal(err)
	}

	logr.Infow("Registering Components")

	// Setup controller
	if err = controller.Add(mgr, kwatcherConfig, logr); err != nil {
		logr.Fatalw("could not setup controller", zap.Error(err))
	}

	logr.Infow("Starting Varnish Watcher")

	logr.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
