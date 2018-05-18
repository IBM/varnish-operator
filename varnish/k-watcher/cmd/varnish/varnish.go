package varnish

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/juju/errors"
)

// Configurator controls the Varnish configuration, namely its VCL files
type Configurator struct {
	vclDir          string
	vclFile         string
	vclFullPath     string
	vclReloadScript string
}

// NewConfigurator creates a new instance of configurator
func NewConfigurator(VCLDir, VCLFile, VCLReloadScript string) *Configurator {
	return &Configurator{
		vclDir:          VCLDir,
		vclFile:         VCLFile,
		vclFullPath:     filepath.Join(VCLDir, VCLFile),
		vclReloadScript: VCLReloadScript,
	}
}

// WriteNewVCL takes the contents of a new VCL file and writes it to disk
func (vc *Configurator) WriteNewVCL(VCL io.Reader) error {
	file, err := os.Create(vc.vclFullPath)
	if err != nil {
		return errors.Annotate(err, "couldn't open backends file for writing")
	}
	defer file.Close()

	if _, err = io.Copy(file, VCL); err != nil {
		return errors.Annotate(err, "couldn't write backends to file")
	}
	return nil
}

// Reload tells Varnish to reload its VCL files given a path to a script that does so
func (vc *Configurator) Reload() error {
	out, err := exec.Command(vc.vclReloadScript).CombinedOutput()
	if err != nil {
		return errors.Annotate(err, string(out))
	}
	return nil
}

// ReloadWithVCL combines WriteNewVCL and Reload into one step
func (vc *Configurator) ReloadWithVCL(VCL io.Reader) error {
	if err := vc.WriteNewVCL(VCL); err != nil {
		return errors.Trace(err)
	}
	if err := vc.Reload(); err != nil {
		return errors.Trace(err)
	}
	return nil
}
