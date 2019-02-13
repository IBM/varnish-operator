package controller

import (
	"bytes"
	"text/template"

	"github.com/juju/errors"
)

func (r *ReconcileVarnish) resolveTemplate(tmplBytes []byte, targetPort, varnishPort int32, backends, varnishNodes []PodInfo) ([]byte, error) {
	tmplName := "backends"
	tmpl, err := template.New(tmplName).Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return nil, errors.Annotate(err, "could not parse template")
	}

	data := map[string]interface{}{
		"Backends":     backends,
		"TargetPort":   targetPort,
		"VarnishNodes": varnishNodes,
		"VarnishPort":  varnishPort,
	}

	var b bytes.Buffer
	b.WriteString("// This file is generated. Do not edit manually, as changes will be destroyed\n\n")
	if err = tmpl.ExecuteTemplate(&b, tmplName, data); err != nil {
		return nil, errors.Annotatef(err, "problem resolving template")
	}
	return b.Bytes(), nil
}
