package controller

import (
	"icm-varnish-k8s-operator/pkg/controller/varnishservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, varnishservice.Add)
}
