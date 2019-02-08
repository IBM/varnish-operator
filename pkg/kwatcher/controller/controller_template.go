package controller

import (
	"bytes"
	"text/template"

	"github.com/juju/errors"
)

func (r *ReconcileVarnish) resolveTemplate(tmplBytes []byte, targetPort int32, backends []Backend, varnishNodes []VarnishNode) ([]byte, error) {
	tmplName := "backends"
	tmpl, err := template.New(tmplName).Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return nil, errors.Annotate(err, "could not parse template")
	}

	data := map[string]interface{}{
		"Backends":     backends,
		"TargetPort":   targetPort,
		"VarnishNodes": varnishNodes,
	}

	var b bytes.Buffer
	if err = tmpl.ExecuteTemplate(&b, tmplName, data); err != nil {
		return nil, errors.Annotatef(err, "problem resolving template")
	}
	return b.Bytes(), nil
}
