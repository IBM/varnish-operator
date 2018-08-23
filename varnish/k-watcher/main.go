package main

import (
	"flag"
	"fmt"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/logger"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/util"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/varnish"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/watch"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/juju/errors"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// config that reads in env variables
type config struct {
	AppSelectorString string `env:"APP_SELECTOR_STRING,required"`
	BackendsFile      string `env:"BACKENDS_FILE,required"`
	Namespace         string `env:"NAMESPACE,required"`
	TargetPort        int32  `env:"TARGET_PORT,required"`
	VCLDir            string `env:"VCL_DIR,required"`
}

var clientset *kubernetes.Clientset
var conf *config
var vc *varnish.Configurator
var vtmpl *varnish.VCLTemplate

// loadConfig uses the codingconcepts/env library to read in environment variables into a struct
func loadConfig() (*config, error) {
	c := config{}
	int32Type := reflect.TypeOf(int32(0))
	int32Parse := func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, errors.Errorf("%s is not an int32", v)
		}
		return int32(i), nil
	}

	parsers := env.CustomParsers{
		int32Type: int32Parse,
	}

	if err := env.ParseWithFuncs(&c, parsers); err != nil {
		return &c, errors.Trace(err)
	}
	return &c, nil
}

// getBackends pulls out the valid backends from the endpoints list.
// This entails pulling out only those address/port combos whose ports match the incoming TargetPort
func getBackends(ep *v1.Endpoints, targetPort int32) (backendList []string) {
	if ep == nil {
		return
	}
	for _, endpoint := range ep.Subsets {
		for _, address := range endpoint.Addresses {
			for _, port := range endpoint.Ports {
				if port.Port == targetPort {
					backendList = append(backendList, address.IP)
					break
				}
			}
		}
	}
	return
}

// castToEp takes an interface, and attempts to cast it to a v1.Endpoints
func castToEp(obj interface{}) (*v1.Endpoints, error) {
	var ep *v1.Endpoints
	var ok bool
	if obj != nil {
		if ep, ok = obj.(*v1.Endpoints); !ok {
			return nil, errors.Errorf("invalid type: %s is not of type *v1.Endpoints. (contents are %+v)", reflect.TypeOf(obj).String(), obj)
		}
	}
	return ep, nil
}

// onChange runs whenever kubernetes reports a change in the backend endpoints.
// It rewrites the varnish VCL file and then reloads varnish with the new file.
func onChange(oldObj, newObj interface{}) {
	oldEp, err := castToEp(oldObj)
	if err != nil {
		logger.Error(err, "oldEp was not the correct type")
		return
	}

	newEp, err := castToEp(newObj)
	if err != nil {
		logger.Error(err, "newEp was not the correct type")
		return
	}

	oldBackends := getBackends(oldEp, conf.TargetPort)
	newBackends := getBackends(newEp, conf.TargetPort)

	removed, added := util.DiffBackends(oldBackends, newBackends)

	VCL, err := vtmpl.GenerateVCL(newBackends, conf.TargetPort)
	if err != nil {
		logger.Error(err, "Could not generate VCL")
		return
	}

	if err = vc.ReloadWithVCL(VCL); err != nil {
		logger.Error(err, "could not reload varnish with new VCL")
		return
	}

	logger.Info("backends updated", "added", added, "removed", removed)
}

// initializes:
// config (as read from environment variables)
// kubernetes client
// varnish vcl template file
// varnish configurator (which writes the template file and reloads varnish)
func init() {
	var err error

	conf, err = loadConfig()
	if err != nil {
		logger.Panic(err, "missing environment variables")
	}

	kubecfgFilepath := flag.String("kubecfg", "", "Path to kube config")
	flag.Parse()

	var kubecfg *rest.Config
	if *kubecfgFilepath == "" {
		kubecfg, err = rest.InClusterConfig()
	} else {
		kubecfg, err = clientcmd.BuildConfigFromFlags("", *kubecfgFilepath)
	}
	if err != nil {
		logger.Panic(err, "couldn't get config")
	}

	clientset, err = kubernetes.NewForConfig(kubecfg)
	if err != nil {
		logger.Panic(err, "couldn't create clientset")
	}

	templateFilename := fmt.Sprintf("%s.tmpl", conf.BackendsFile)
	if vtmpl, err = varnish.NewVCLTemplate(conf.VCLDir, templateFilename); err != nil {
		logger.Panic(err, "couldn't create VCL Template")
	}

	vc = varnish.NewConfigurator(conf.VCLDir, conf.BackendsFile, "vcl_reload")
}

// starts watch of the endpoints backend in the cluster, triggering a new VCL + varnish reload on changes
func main() {
	defer logger.Sync()
	_, controller, err := watch.Resource(
		clientset.CoreV1().RESTClient(),
		"endpoints",
		conf.Namespace,
		conf.AppSelectorString,
		onChange,
	)
	if err != nil {
		logger.Panic(err, "could not initialize endpoints watcher")
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(stopCh)

	logger.Info("endpoints watcher has started",
		"resource", "endpoints",
		"namespace", conf.Namespace,
		"appSelector", conf.AppSelectorString,
		"watchedPort", conf.TargetPort)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}
