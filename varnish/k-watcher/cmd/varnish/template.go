package varnish

import (
	"bytes"
	"icm-varnish-k8s-operator/varnish/k-watcher/cmd/util"
	"io"
	"path/filepath"
	"text/template"

	"github.com/juju/errors"
)

// VCLTemplate represents the template used to write new VCL configurations
type VCLTemplate struct {
	t        *template.Template
	filename string
}

// NewVCLTemplate creates a new instance of the VCL template, for generation later
func NewVCLTemplate(dir string, filename string) (*VCLTemplate, error) {
	t := template.New(filepath.Join(dir, filename))
	t, err := t.ParseFiles(filepath.Join(dir, filename))
	if err != nil {
		return nil, errors.Annotatef(err, "could not parse template %s", filename)
	}

	vt := VCLTemplate{
		t:        t,
		filename: filename,
	}
	return &vt, nil
}

// GenerateVCL uses a list of backends to generate a new VCL template, passing back a Reader representing the new template
func (vt *VCLTemplate) GenerateVCL(backends []util.Backend) (io.Reader, error) {
	t, err := vt.t.Clone()
	if err != nil {
		return nil, errors.Annotate(err, "problem initializing template")
	}

	var b bytes.Buffer

	if err := t.ExecuteTemplate(&b, vt.filename, backends); err != nil {
		return nil, errors.Annotate(err, "problem resolving template")
	}
	return &b, nil
}
