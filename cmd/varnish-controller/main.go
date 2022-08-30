// The app watches kubernetes endpoints api for changes in deployment
// and then re-writes varnish vcl file with any new/removed backends.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ibm/varnish-operator/pkg/names"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/ibm/varnish-operator/api/v1alpha1"
	"github.com/ibm/varnish-operator/pkg/logger"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/config"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/controller"
	varnishMetrics "github.com/ibm/varnish-operator/pkg/varnishcontroller/metrics"
	"github.com/ibm/varnish-operator/pkg/varnishcontroller/varnishadm"

	controllerMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

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
	logr.Infof("Prometheus metrics exporter port: %d", v1alpha1.VarnishControllerMetricsPort)

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

	if _, ok := os.LookupEnv("INIT_CONTAINER"); ok {
		initHaproxyConfig(clientConfig, logr, varnishControllerConfig)
		return
	}

	vMetrics := varnishMetrics.NewVarnishControllerMetrics()
	controllerMetrics.Registry.MustRegister(vMetrics.VCLCompilationError)

	mgr, err := ctrl.NewManager(clientConfig, ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: fmt.Sprintf(":%d", v1alpha1.HealthCheckPort),
		MetricsBindAddress:     fmt.Sprintf(":%d", v1alpha1.VarnishControllerMetricsPort),
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

	if err = controller.SetupVarnishReconciler(mgr, varnishControllerConfig, varnishAdm, vMetrics, logr); err != nil {
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

func initHaproxyConfig(clientConfig *rest.Config, logr *logger.Logger, varnishControllerConfig *config.Config) {
	cs, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		logr.With(zap.Error(err)).Fatalf("unable to create k8s client")
	}
	haproxyConfigMapName := names.HaproxyConfigMap(varnishControllerConfig.VarnishClusterName)
	configFound := false
	for !configFound {
		logr.Infof("attempting to retrieve configmap: %s", haproxyConfigMapName)
		configMaps := cs.CoreV1().ConfigMaps(varnishControllerConfig.Namespace)
		if haproxyConfigMap, err := configMaps.Get(context.Background(), haproxyConfigMapName, metav1.GetOptions{}); err == nil {
			cfgData := haproxyConfigMap.Data[v1alpha1.HaproxyConfigFileName]
			haproxyConfigFileName := v1alpha1.HaproxyConfigDir + "/" + v1alpha1.HaproxyConfigFileName
			if err := os.WriteFile(haproxyConfigFileName, []byte(cfgData), 0644); err != nil {
				logr.With(zap.Error(err)).Fatalf("unable to write haproxy config file: %s", haproxyConfigFileName)
			}
			logr.Infof("haproxy config has been written\n%s", cfgData)
			configFound = true
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}
