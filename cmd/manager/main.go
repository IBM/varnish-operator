package main

import (
	"flag"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"log"

	"icm-varnish-k8s-operator/pkg/apis"
	"icm-varnish-k8s-operator/pkg/logger"
	vscfg "icm-varnish-k8s-operator/pkg/varnishservice/config"
	"icm-varnish-k8s-operator/pkg/varnishservice/controller"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func init() {
	flag.Parse()
}

func main() {
	operatorConfig, err := vscfg.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Get a config to talk to the apiserver
	clientConfig, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	logr := logger.NewLogger(operatorConfig.LogFormat, operatorConfig.LogLevel)

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(clientConfig, manager.Options{
		LeaderElection:          operatorConfig.LeaderElection,
		LeaderElectionID:        operatorConfig.LeaderElectionID,
		LeaderElectionNamespace: operatorConfig.Namespace,
	})
	if err != nil {
		logr.Fatal(err)
	}

	logr.Infow("Registering Components")

	// Setup Scheme for all resources
	v1alpha1.Init(operatorConfig)
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		logr.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.Add(mgr, operatorConfig, logr); err != nil {
		logr.Fatal(err)
	}

	logr.Infow("Starting Varnish Operator")

	// Start the Cmd
	logr.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
