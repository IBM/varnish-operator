package main

import (
	"flag"
	"fmt"
	"icm-varnish-k8s-operator/pkg/apis"
	"icm-varnish-k8s-operator/pkg/logger"
	vscfg "icm-varnish-k8s-operator/pkg/varnishservice/config"
	"icm-varnish-k8s-operator/pkg/varnishservice/controller"
	"icm-varnish-k8s-operator/pkg/varnishservice/webhooks"
	"log"

	"go.uber.org/zap"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var (
	Version = "undefined" //will be overwritten by the correct version during docker build
)

func main() {
	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.Parse()
	operatorConfig, err := vscfg.LoadConfig()
	if err != nil {
		log.Fatalf("unable to read env vars: %v", err)
	}

	logr := logger.NewLogger(operatorConfig.LogFormat, operatorConfig.LogLevel)

	logr.Infof("Version: %s", Version)
	logr.Infof("Leader election enabled: %t", operatorConfig.LeaderElectionEnabled)
	logr.Infof("Log level: %s", operatorConfig.LogLevel.String())
	logr.Infof("Prometheus metrics exporter port: %d", operatorConfig.ContainerMetricsPort)

	logr = logr.With(logger.FieldOperatorVersion, Version)

	// Get a config to talk to the apiserver
	clientConfig, err := config.GetConfig()
	if err != nil {
		logr.Fatalf("unable to set up client config: %v", err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(clientConfig, manager.Options{
		MetricsBindAddress:      fmt.Sprintf(":%d", operatorConfig.ContainerMetricsPort),
		LeaderElection:          operatorConfig.LeaderElectionEnabled,
		LeaderElectionID:        operatorConfig.LeaderElectionID,
		LeaderElectionNamespace: operatorConfig.Namespace,
	})
	if err != nil {
		logr.With(zap.Error(err)).Fatal("unable to set up overall controller manager")
	}

	logr.Infow("Registering Components")

	// Setup Scheme for all resources
	logr.Infow("Setting up scheme")
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		logr.With(zap.Error(err)).Fatal("unable to add APIs to scheme")
	}

	// Setup all Controllers
	logr.Infow("Setting up controller")
	if err := controller.Add(mgr, operatorConfig, logr); err != nil {
		logr.With(zap.Error(err)).Fatal("unable to register controllers to the manager")
	}

	logr.Infow("Setting up webhooks")
	if err := webhooks.InstallWebhooks(mgr, operatorConfig, logr); err != nil {
		logr.With(zap.Error(err)).Fatal("unable to register webhooks to the manager")
	}

	// Start the Cmd
	logr.Infow("Starting Varnish Operator")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logr.With(zap.Error(err)).Fatal("unable to run the manager")
	}
}
