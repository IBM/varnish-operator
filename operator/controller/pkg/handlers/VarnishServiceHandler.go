package handlers

import (
	icmapiv1alpha1 "icm-varnish-k8s-operator/operator/controller/pkg/apis/icm/v1alpha1"
	"reflect"

	"github.com/juju/errors"

	log "github.com/sirupsen/logrus"
)

// VarnishServiceHandler describes the functions that handle events coming in for the VarnishService CRD
type VarnishServiceHandler struct {
}

// ObjectAdded prints out the VarnishService
func (h *VarnishServiceHandler) ObjectAdded(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}
	log.Infof("adding %+v", vs)
	return nil
}

// ObjectUpdated prints out the VarnishService
func (h *VarnishServiceHandler) ObjectUpdated(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}
	log.Infof("adding %+v", vs)
	return nil
}

// ObjectDeleted prints out the VarnishService
func (h *VarnishServiceHandler) ObjectDeleted(obj interface{}) error {
	vs, ok := obj.(*icmapiv1alpha1.VarnishService)
	if !ok {
		return errors.NotValidf("object was not of type VarnishService. Found %s", reflect.TypeOf(vs))
	}
	log.Infof("adding %+v", vs)
	return nil
}
