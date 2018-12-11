package main

import (
	"flag"
	"log"

	"icm-varnish-k8s-operator/pkg/apis"
	vscfg "icm-varnish-k8s-operator/pkg/config"
	"icm-varnish-k8s-operator/pkg/controller"
	"icm-varnish-k8s-operator/pkg/logger"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func init() {
	flag.Parse()
}

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		LeaderElection:          vscfg.GlobalConf.LeaderElection,
		LeaderElectionID:        vscfg.GlobalConf.LeaderElectionID,
		LeaderElectionNamespace: vscfg.GlobalConf.Namespace,
	})

	if err != nil {
		log.Fatal(err)
	}

	logger.Infow("Registering Components")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Fatal(err)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Fatal(err)
	}

	logger.Infow("Starting Varnish Operator")

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}
