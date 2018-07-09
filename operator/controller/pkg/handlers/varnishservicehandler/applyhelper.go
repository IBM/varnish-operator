package varnishservicehandler

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/api/meta"
	"github.com/juju/errors"
	metav1unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var metaAccessor = meta.NewAccessor()
var LastAppliedConfigurationKey = "icm.ibm.com/last-applied-configuration"

func addApplyAnnotation(o runtime.Object) error {
	annots, err := metaAccessor.Annotations(o)
	if err != nil {
		return errors.Annotate(err,"object does not have valid metadata")
	}
	if annots == nil {
		annots = map[string]string{}
	}
	lastApplied, err := runtime.Encode(metav1unstructured.UnstructuredJSONScheme, o)
	if err != nil {
		return errors.Annotate(err, "problem encoding obj")
	}
	annots[LastAppliedConfigurationKey] = string(lastApplied)
	if err = metaAccessor.SetAnnotations(o, annots); err != nil {
		return errors.Annotate(err, "could not set apply annotations on object")
	}
	return nil
}
