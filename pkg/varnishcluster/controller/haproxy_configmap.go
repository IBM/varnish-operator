package controller

import (
	"bytes"
	"context"
	"text/template"

	vclabels "github.com/ibm/varnish-operator/pkg/labels"
	"github.com/pkg/errors"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *ReconcileVarnishCluster) reconcileHaproxyConfigMap(ctx context.Context, instance *vcapi.VarnishCluster) (*v1.ConfigMap, error) {
	cmLabels := vclabels.CombinedComponentLabels(instance, vcapi.HaproxyConfigMapName)
	t, err := template.New("haproxy-config").Parse(haproxyConfigTemplate)
	if err != nil {
		return &v1.ConfigMap{}, errors.WithStack(err)
	}

	var b bytes.Buffer
	err = t.Execute(&b, instance.Spec.HaproxySidecar)
	if err != nil {
		return &v1.ConfigMap{}, errors.WithStack(err)
	}

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vcapi.HaproxyConfigMapName,
			Labels:    cmLabels,
			Namespace: instance.Namespace,
		},
		Data: map[string]string{
			"haproxy.cfg": b.String(),
		},
	}
	return cm, nil
}
