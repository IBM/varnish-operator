package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	vsclientset "icm-varnish-k8s-operator/operator/controller/pkg/client/clientset/versioned"
	"icm-varnish-k8s-operator/operator/controller/pkg/controller"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var vsclient *vsclientset.Clientset

// references call site for logError/logAndPanic
func generateErrorStack(err error, msg string) string {
	wrapped := errors.NewErrWithCause(err, msg)
	wrapped.SetLocation(2)
	return errors.ErrorStack(&wrapped)
}

// logError logs the err and message
func logError(err error, msg string) {
	log.Error(generateErrorStack(err, msg))
}

// logAndPanic logs the err and message, then panics (exits the program)
func logAndPanic(err error, msg string) {
	log.Panic(generateErrorStack(err, msg))
}

func init() {
	var err error

	kubecfgFilepath := flag.String("kubecfg", "", "Path to kube config")
	flag.Parse()

	var kubecfg *rest.Config
	if *kubecfgFilepath == "" {
		kubecfg, err = rest.InClusterConfig()
	} else {
		kubecfg, err = clientcmd.BuildConfigFromFlags("", *kubecfgFilepath)
	}
	if err != nil {
		logAndPanic(err, "couldn't get config")
	}

	vsclient, err = vsclientset.NewForConfig(kubecfg)
	if err != nil {
		logAndPanic(err, "couldn't create clientset")
	}
}

func main() {
	c := controller.NewVarnishServiceController(vsclient, 3)
	stopCh := make(chan struct{})
	defer close(stopCh)
	go c.Run(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}
