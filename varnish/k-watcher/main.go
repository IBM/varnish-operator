package main

import (
	"flag"
	"fmt"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/util"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/varnish"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/watch"
	"reflect"
	"time"

	"github.com/codingconcepts/env"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// config that reads in env variables
type config struct {
	AppName       string `env:"APP_NAME" required:"true"`
	BackendsFile  string `env:"BACKENDS_FILE" required:"true"`
	Namespace     string `env:"NAMESPACE" required:"true"`
	PortWatchName string `env:"PORT_WATCH_NAME" required:"true"`
	VCLDir        string `env:"VCL_DIR" required:"true"`
}

var clientset *kubernetes.Clientset
var conf *config
var vc *varnish.Configurator
var vtmpl *varnish.VCLTemplate

// logError logs the err and message
func logError(err error, msg string) {
	log.WithField("error", errors.Details(err)).Error(msg)
}

// logAndPanic logs the err and message, then panics (exits the program)
func logAndPanic(err error, msg string) {
	log.WithField("error", errors.Details(err)).Panic(msg)
}

// loadConfig uses the codingconcepts/env library to read in environment variables into a struct
func loadConfig() (*config, error) {
	c := config{}
	if err := env.Set(&c); err != nil {
		return &c, errors.Trace(err)
	}
	return &c, nil
}

// getBackends pulls out the valid backends from the endpoints list.
// This entails pulling out only those address/port combos whose ports have the PortWatchName
func getBackends(ep *v1.Endpoints) (backendList []util.Backend) {
	if ep == nil {
		return
	}
	for _, endpoint := range ep.Subsets {
		for _, address := range endpoint.Addresses {
			for _, port := range endpoint.Ports {
				if port.Name == conf.PortWatchName {
					backendList = append(backendList, util.Backend{Host: address.IP, Port: port.Port})
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
		logError(err, "oldEp was not the correct type")
		return
	}

	newEp, err := castToEp(newObj)
	if err != nil {
		logError(err, "newEp was not the correct type")
		return
	}

	oldBackends := getBackends(oldEp)
	newBackends := getBackends(newEp)

	removed, added := util.DiffBackends(oldBackends, newBackends)

	VCL, err := vtmpl.GenerateVCL(newBackends)
	if err != nil {
		logError(errors.Trace(err), "Could not generate VCL")
		return
	}

	if err = vc.ReloadWithVCL(VCL); err != nil {
		logError(errors.Trace(err), "could not reload varnish with new VCL")
		return
	}

	log.WithFields(log.Fields{
		"added":   added,
		"removed": removed,
	}).Info("backends updated")
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
		logAndPanic(errors.Trace(err), "missing environment variables")
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
		logAndPanic(errors.Trace(err), "couldn't get config")
	}

	clientset, err = kubernetes.NewForConfig(kubecfg)
	if err != nil {
		logAndPanic(errors.Trace(err), "couldn't create clientset")
	}

	templateFilename := fmt.Sprintf("%s.tmpl", conf.BackendsFile)
	if vtmpl, err = varnish.NewVCLTemplate(conf.VCLDir, templateFilename); err != nil {
		logAndPanic(errors.Trace(err), "couldn't create VCL Template")
	}

	vc = varnish.NewConfigurator(conf.VCLDir, conf.BackendsFile, "vcl_reload")
}

// starts watch of the endpoints backend in the cluster, triggering a new VCL + varnish reload on changes
func main() {
	_, controller, err := watch.WatchResource(
		clientset.CoreV1().RESTClient(),
		"endpoints",
		conf.Namespace,
		fields.OneTermEqualSelector("metadata.name", conf.AppName),
		onChange,
	)
	if err != nil {
		logAndPanic(errors.Trace(err), "could not initialize endpoints watcher")
	}

	stopCh := make(chan struct{})
	go controller.Run(stopCh)

	log.WithFields(log.Fields{
		"resource":      "endpoints",
		"namespace":     conf.Namespace,
		"metadata.name": conf.AppName,
		"portWatchName": conf.PortWatchName,
	}).Info("endpoints watcher has started")

	for {
		time.Sleep(time.Second)
	}
}
