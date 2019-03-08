package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"icm-varnish-k8s-operator/pkg/apis/icm/v1alpha1"
	"icm-varnish-k8s-operator/pkg/kwatcher/events"
	"os/exec"
	"strings"
	"time"

	"k8s.io/api/core/v1"

	"github.com/juju/errors"
)

const (
	VCLStatusAvailable = "available"
	VCLStatusActive    = "active"

	VCLTemperatureCold = "cold"
	VCLTemperatureWarm = "warm"

	// For VCL version name we use config map resource version which is a number.
	// Varnish doesn't accept config name that have numbers in the beginning. Even if it is disguised as strings (e.g. "1243").
	// For that reasons we prepend this prefix.
	VCLVersionPrefix = "v"
)

type VCLConfig struct {
	Status        string
	Name          string
	Label         bool
	Temperature   string
	ReferencedVCL *string //nil if Label == false
}

func (r *ReconcileVarnish) reconcileVarnish(vs *v1alpha1.VarnishService, pod *v1.Pod, cm *v1.ConfigMap) error {
	r.logger.Debugw("Starting varnish reload...")
	start := time.Now()
	out, err := exec.Command("vcl_reload", createVCLConfigName(cm.GetResourceVersion()), vs.Spec.VCLConfigMap.EntrypointFile).CombinedOutput()
	if err != nil {
		if isVCLCompilationError(err) {
			vsEventMsg := "VCL compilation failed for pod " + pod.Name + ". See pod logs for details"
			podEventMsg := "VCL compilation failed. See logs for details"
			r.eventHandler.Warning(pod, events.EventReasonVCLCompilationError, podEventMsg)
			r.eventHandler.Warning(vs, events.EventReasonVCLCompilationError, vsEventMsg)
		} else {
			podEventMsg := "Varnish reload failed for pod " + pod.Name + ". See pod logs for details"
			vsEventMsg := "Varnish reload failed. See logs for details"
			r.eventHandler.Warning(pod, events.EventReasonReloadError, podEventMsg)
			r.eventHandler.Warning(vs, events.EventReasonReloadError, vsEventMsg)
		}
		return errors.Annotate(err, string(out))
	}
	r.logger.Debugf("Varnish successfully reloaded in %f seconds", time.Since(start).Seconds())
	return nil
}

// getActiveVCLConfig returns the VCL config currently used in Varnish
func getActiveVCLConfig() (*VCLConfig, error) {
	configsList, err := getVCLConfigsList()
	if err != nil {
		return nil, err
	}

	var activeVersion *VCLConfig
	for _, vclConfig := range configsList {
		if vclConfig.Status == VCLStatusActive {
			activeVersion = &vclConfig
		}
	}

	if activeVersion == nil {
		// That means that Varnish is in not started/invalid state. Return an error to trigger an another reconcile event
		return nil, errors.NotFoundf("No active VCL configuration found")
	}

	return activeVersion, nil
}

func getVCLConfigsList() ([]VCLConfig, error) {
	out, err := exec.Command("vcl_list").CombinedOutput()
	if err != nil {
		return nil, errors.Annotate(err, string(out))
	}

	configs, err := parseVCLConfigsList(out)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return configs, nil
}

func parseVCLConfigsList(commandOutput []byte) ([]VCLConfig, error) {
	var configs []VCLConfig
	lines := bufio.NewScanner(bytes.NewReader(commandOutput))
	for lines.Scan() {
		columns := strings.Fields(lines.Text())
		switch len(columns) {
		case 0: //empty string
			continue
		case 4: //config without a label
			temp := strings.Split(columns[1], "/")
			configs = append(configs, VCLConfig{Status: columns[0], Name: columns[3], Label: false, Temperature: temp[1]})
		case 6: //labeled config or a label itself
			var refVCL *string
			temp := strings.Split(columns[1], "/")
			isLabel := temp[0] == "label"
			if isLabel {
				refVCL = &columns[5]
			}
			config := VCLConfig{Status: columns[0], Name: columns[3], Label: isLabel, Temperature: temp[1], ReferencedVCL: refVCL}
			configs = append(configs, config)
		default:
			return nil, errors.New("unknown VCL config format")
		}
	}
	return configs, nil
}

func isVCLCompilationError(err error) bool {
	if err == nil {
		return false
	}

	scanner := bufio.NewScanner(strings.NewReader(err.Error()))
	for scanner.Scan() {
		if scanner.Text() == "VCL compilation failed" {
			return true
		}
	}

	return false
}

// creates the VCL config name from config map version
func createVCLConfigName(configMapVersion string) string {
	return fmt.Sprintf("%s-%s-%d", VCLVersionPrefix, configMapVersion, time.Now().Unix())
}

// returns the config name the was used for VCL config name creation
func extractConfigMapVersion(VCLConfigName string) string {
	parts := strings.Split(VCLConfigName, "-")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2]
}
