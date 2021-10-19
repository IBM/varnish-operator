package varnishadm

import (
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

const (
	//VarninshAdmBinary binary to execute to control a varnish instance
	VarninshAdmBinary = "varnishadm"
	//VCLStatusAvailable - VCL configuration status
	VCLStatusAvailable = "available"
	//VCLStatusActive - VCL configuration is active now
	VCLStatusActive = "active"
	//VCLStatusDiscarded - VCL configuration is discarded but still in use by in-flight requests
	VCLStatusDiscarded = "discarded"
	//VCLTemperatureCold - vanish VCL "temperature"
	VCLTemperatureCold = "cold"
	//VCLTemperatureWarm for preloaded varnish's VCL
	VCLTemperatureWarm = "warm"
)

// Commander defines the interface to use for call external utilities to manage varnish instance.
// - Ping() check if a varnish instace ready and reachable
// - Reload() try to load a new varnish VCL configuration
// - List() returns the VCL config currently used in varnish
type Commander interface {
	Ping() error
	Reload(version, entry string) ([]byte, error)
	List() ([]VCLConfig, error)
	Discard(vclConfigName string) error
}

// VarnishAdministrator the Commander interface extension by the funtcion which returns active configuration name.
// - GetActiveConfigurationName() returns an active VCL configuration name or error
type VarnishAdministrator interface {
	GetActiveConfigurationName() (name string, err error)
	Commander
}

// NewVarnishAdministartor returns a wrapper over system command which allows ping, list and relod a varnishd
// it accepts following parameters:
// - timeout: an overall timeout to try to reach a varnish instance
// - delay: a delay between ping retry
// - vclBase: base directory for varnish's VCL files
// - args: varnisadm utility CLI parameters
func NewVarnishAdministartor(timeout, delay time.Duration, vclBase string, args []string) *VarnishAdm {
	return &VarnishAdm{
		binary:         VarninshAdmBinary,
		vclBase:        vclBase,
		varnishAdmArgs: sanitizeVarnishArgs(args),
		pingTimeout:    timeout,
		pingDelay:      delay,
		execute:        execCommandProvider,
	}
}

// VarnishAdm is a structure which implements Commander interface.
// it wraps call to external binaries to use varnish_adm binary from a varnish distribution
// to Commander interface specific methods
type VarnishAdm struct {
	binary         string
	vclBase        string
	pingTimeout    time.Duration
	pingDelay      time.Duration
	varnishAdmArgs []string
	execute        executorProvider
}

// executor an interface compatible with *exec.Cmd method CombinedOutput()
// added for testability
type executor interface {
	CombinedOutput() ([]byte, error)
}

type executorProvider func(name string, arg ...string) executor

// Ping try to reach varnish instance to ensure it is up and running.
// It applies pingDelay value as a maximum time to try to reach the varnish instance
// it is a wrapper over varnishadm ping command.
func (v *VarnishAdm) Ping() error {
	out := make(chan struct{})
	done := make(chan struct{})
	defer close(done)
	args := []string{"ping"}
	go func(out, done chan struct{}) {
		defer close(out)
		for {
			_, err := v.run(args)
			select {
			case <-done:
				return
			default:
			}
			if err == nil {
				return
			}
			time.Sleep(v.pingDelay)
		}
	}(out, done)
	select {
	case <-time.After(v.pingTimeout):
		return errors.New("varnish is unreachable")
	case <-out:
		return nil
	}
}

//Reload loads new VCL configuration into the varnish instance. Accepts two string parameters
// - version string, a version which describes the configuration
// - entrypoint string, a start filename to use as a new VCL configuration
// it is a wrapper over varnishadm vcl.load and vcl.use commands combination
func (v *VarnishAdm) Reload(version, entry string) ([]byte, error) {
	out, err := v.load(version, entry)
	if err != nil {
		return out, err
	}
	return v.use(version)
}

// Discard deletes an existing VCL from the Varnish instance
// it is a wrapper over varnishadm vcl.discard command
func (v *VarnishAdm) Discard(vclConfigName string) error {
	out, err := v.run(append(v.varnishAdmArgs, "vcl.discard", vclConfigName))
	if err != nil {
		return errors.Wrap(err, string(out))
	}

	return nil
}

func (v *VarnishAdm) run(args []string) ([]byte, error) {
	return v.execute(v.binary, args...).CombinedOutput()
}

func (v *VarnishAdm) load(version, entry string) ([]byte, error) {
	args := append(v.varnishAdmArgs, "vcl.load", version, v.vclBase+"/"+entry)
	return v.run(args)
}

func (v *VarnishAdm) use(version string) ([]byte, error) {
	args := append(v.varnishAdmArgs, "vcl.use", version)
	return v.run(args)
}

func sanitizeVarnishArgs(input []string) []string {
	out := make([]string, 0, len(input))
	for _, val := range input {
		if val == "" || val == "\t" || val == " " {
			continue
		}
		out = append(out, val)

	}
	return out
}

func execCommandProvider(name string, args ...string) executor {
	cmd := exec.Command(name, args...)
	return cmd
}

//GetActiveConfigurationName parses varnishadm list output and compute a name of active configuration
func (v *VarnishAdm) GetActiveConfigurationName() (string, error) {
	active, err := v.getActiveVCLConfig()
	if err != nil {
		return "", err
	}
	return active.Name, nil
}

// getActiveVCLConfig returns the VarnishClusterVCL config currently used in VarnishClusterVarnish
func (v *VarnishAdm) getActiveVCLConfig() (*VCLConfig, error) {
	configsList, err := v.List()
	if err != nil {
		return nil, err
	}
	for _, vclConfig := range configsList {
		if vclConfig.Status == VCLStatusActive {
			return &vclConfig, nil //active configuration found
		}
	}
	// That means that VarnishClusterVarnish is in not started/invalid state. Return an error to trigger an another reconcile event
	return nil, errors.Errorf("No active VCL configuration found")
}
