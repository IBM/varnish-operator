package controller

import (
	"icm-varnish-k8s-operator/operator/controller/pkg/controller/varnishservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, varnishservice.Add)
}
