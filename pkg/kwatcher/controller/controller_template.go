package controller

import (
	"bytes"
	"icm-varnish-k8s-operator/pkg/kwatcher/backends"
	"text/template"

	"github.com/juju/errors"
)

func resolveTemplate(tmplBytes []byte, targetPort int32, backends []backends.Backend) ([]byte, error) {
	tmplName := "backends"
	tmpl, err := template.New(tmplName).Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return nil, errors.Annotate(err, "could not parse template")
	}

	data := map[string]interface{}{
		"Backends":   backends,
		"TargetPort": targetPort,
	}

	var b bytes.Buffer
	if err = tmpl.ExecuteTemplate(&b, tmplName, data); err != nil {
		return nil, errors.Annotatef(err, "problem resolving template")
	}
	return b.Bytes(), nil
}
