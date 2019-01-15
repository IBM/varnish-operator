// The app watches kubernetes endpoints api for changes in deployment
// and then re-writes varnish vcl file with any new/removed backends.

package main

import (
	"flag"
	"icm-varnish-k8s-operator/pkg/kwatcher/config"
	"icm-varnish-k8s-operator/pkg/kwatcher/controller"
	"icm-varnish-k8s-operator/pkg/kwatcher/logger"
	"log"

	"go.uber.org/zap"

	kconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func init() {
	flag.Parse()
}

func main() {
	cfg, err := kconfig.GetConfig()
	if err != nil {
		logger.Panicw("could not load config", zap.Error(err))
	}

	mgr, err := manager.New(cfg, manager.Options{Namespace: config.GlobalConf.Namespace})
	if err != nil {
		logger.Panicw("could not initialize manager", zap.Error(err))
	}

	logger.Infow("Registering Components")

	// Setup controller
	if err = controller.Add(mgr); err != nil {
		logger.Panicw("could not setup controller", zap.Error(err))
	}

	logger.Infow("Starting Varnish Watcher")

	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
