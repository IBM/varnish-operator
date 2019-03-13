package controller

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func (r *ReconcileVarnish) resolveTemplate(tmplStr string, targetPort, varnishPort int32, backends, varnishNodes []PodInfo) (string, error) {
	tmplName := "backends"
	tmpl, err := template.New(tmplName).Option("missingkey=error").Parse(tmplStr)
	if err != nil {
		return "", errors.Wrap(err, "could not parse template")
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
		return "", errors.Wrapf(err, "problem resolving template")
	}
	return b.String(), nil
}

func (r *ReconcileVarnish) resolveTemplates(tmplStrs map[string]string, targetPort, varnishPort int32, backends, varnishNodes []PodInfo) (map[string]string, error) {
	data := map[string]interface{}{
		"Backends":     backends,
		"TargetPort":   targetPort,
		"VarnishNodes": varnishNodes,
		"VarnishPort":  varnishPort,
	}

	out := make(map[string]string, len(tmplStrs))
	for tmplFileName, tmplStr := range tmplStrs {
		tmpl, err := template.New(tmplFileName).Option("missingkey=error").Parse(tmplStr)
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse template %s", tmplFileName)
		}

		var b bytes.Buffer
		b.WriteString("// This file is generated. Do not edit manually, as changes will be destroyed\n\n")
		if err = tmpl.ExecuteTemplate(&b, tmplFileName, data); err != nil {
			return nil, errors.Wrapf(err, "problem resolving template %s", tmplFileName)
		}
		fileName := strings.TrimSuffix(tmplFileName, ".tmpl")
		out[fileName] = b.String()
	}
	return out, nil
}
