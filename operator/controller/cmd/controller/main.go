package main

import (
	"flag"

	vsclientset "icm-varnish-k8s-operator/operator/controller/pkg/client/clientset/versioned"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	list, err := vsclient.IcmV1alpha1().VarnishServices("default").List(metav1.ListOptions{})
	if err != nil {
		logAndPanic(err, "could not list varnish services")
	}

	for _, vs := range list.Items {
		log.WithFields(log.Fields{
			"replicas": vs.Spec.Replicas,
			"name":     vs.Name,
		}).Info("varnish service")
	}
}
