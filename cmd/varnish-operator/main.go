package main

import (
	"flag"
	"fmt"
	"icm-varnish-k8s-operator/api/v1alpha1"
	"icm-varnish-k8s-operator/pkg/logger"
	vccfg "icm-varnish-k8s-operator/pkg/varnishcluster/config"
	"icm-varnish-k8s-operator/pkg/varnishcluster/controller"
	"log"

	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"go.uber.org/zap"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/go-logr/zapr"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

var (
	Version = "undefined" //will be overwritten by the correct version during docker build
)

const (
	// leaderElectionID defines the name of the ConfigMap acting as the lock for defining the leader
	leaderElectionID = "varnish-operator-leader-election-lock"
)

func main() {
	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.Parse()
	operatorConfig, err := vccfg.LoadConfig()
	if err != nil {
		log.Fatalf("unable to read env vars: %v", err)
	}

	logr := logger.NewLogger(operatorConfig.LogFormat, operatorConfig.LogLevel)

	logr.Infof("Version: %s", Version)
	logr.Infof("Leader election enabled: %t", operatorConfig.LeaderElectionEnabled)
	logr.Infof("Log level: %s", operatorConfig.LogLevel.String())
	logr.Infof("Prometheus metrics exporter port: %d", operatorConfig.MetricsPort)

	logr = logr.With(logger.FieldOperatorVersion, Version)

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

	// Get a config to talk to the apiserver
	clientConfig, err := ctrl.GetConfig()
	if err != nil {
		logr.Fatalf("unable to set up client config: %v", err)
	}

	ctrl.SetLogger(zapr.NewLogger(logr.Desugar())) //set logger for controller-runtime to see internal library logs

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := ctrl.NewManager(clientConfig, ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      fmt.Sprintf(":%d", operatorConfig.MetricsPort),
		LeaderElection:          operatorConfig.LeaderElectionEnabled,
		LeaderElectionID:        leaderElectionID,
		LeaderElectionNamespace: operatorConfig.Namespace,
	})
	if err != nil {
		logr.With(zap.Error(err)).Fatal("unable to set up overall controller manager")
	}

	logr.Infow("Registering Components")

	// Setup all Controllers
	logr.Infow("Setting up controller")
	reconcileChan := make(chan event.GenericEvent)
	if err = controller.SetupVarnishReconciler(mgr, operatorConfig, logr, reconcileChan); err != nil {
		logr.With(zap.Error(err)).Fatalf("unable to setup controller")
	}

	if operatorConfig.WebhooksEnabled {
		logr.Infof("Admission webhooks port: %d", operatorConfig.WebhooksPort)
		mgr.GetWebhookServer().Port = int(operatorConfig.WebhooksPort)
		if err = (&v1alpha1.VarnishCluster{}).SetupWebhookWithManager(mgr); err != nil {
			logr.With(zap.Error(err)).Fatal("unable to create webhook")
		}
		v1alpha1.SetWebhookLogger(logr)
	}

	// +kubebuilder:scaffold:builder

	// Start the Cmd
	logr.Infow("Starting Varnish Operator")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		logr.With(zap.Error(err)).Fatal("unable to run the manager")
	}
}
