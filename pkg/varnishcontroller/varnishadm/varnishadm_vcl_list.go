package varnishadm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//VCLConfig represents varnish's config list
type VCLConfig struct {
	Status        string
	Name          string
	State         string
	Busy          int
	Label         bool
	Temperature   string
	ReferencedVCL *string //nil if Label == false
}

// List returns a list of VCL names which had been loaded into the varnish instance.
// it is a wrapper over varnishadm vcl.list command
func (v *VarnishAdm) List() ([]VCLConfig, error) {
	out, err := v.run(append(v.varnishAdmArgs, "vcl.list", "-j"))
	if err != nil {
		if strings.Contains(string(out), "JSON unimplemented") {

			out, err = v.run(append(v.varnishAdmArgs, "vcl.list"))
			if err != nil {
				return []VCLConfig{}, errors.Wrap(err, string(out))
			}
			configs, err := parseVCLConfigsList(out)
			if err != nil {
				return nil, err
			}
			return configs, nil
		}

		return nil, err
	}

	vclList := vclListResponse{}
	err = json.Unmarshal(out, &vclList)
	if err != nil {
		return []VCLConfig{}, errors.Wrap(err, string(out))
	}

	return vclList.VCLs, nil
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

type vclListResponse struct {
	Version      int
	Command      []string
	ResponseTime time.Time
	VCLs         []VCLConfig
}

func (v *vclListResponse) UnmarshalJSON(data []byte) error {
	rawUmarshalledData := make([]interface{}, 0)
	err := json.Unmarshal(data, &rawUmarshalledData)
	if err != nil {
		return err
	}

	version, ok := rawUmarshalledData[0].(float64)
	if !ok {
		return errors.New("unknown version format")
	}

	v.Version = int(version)

	command, ok := rawUmarshalledData[1].([]interface{})
	if !ok {
		return errors.New("unknown command format")
	}

	for _, arg := range command {
		argStr, ok := arg.(string)
		if !ok {
			return errors.New("unknown arg format")
		}
		v.Command = append(v.Command, argStr)
	}

	respTime, ok := rawUmarshalledData[2].(float64)
	if !ok {
		return errors.New("unknown response time format")
	}

	secs := int64(respTime)
	nsecs := int64((respTime - float64(secs)) * 1e9)
	v.ResponseTime = time.Unix(secs, nsecs)

	for i := 3; i < len(rawUmarshalledData); i++ {
		vclMap, ok := rawUmarshalledData[i].(map[string]interface{})
		if !ok {
			return errors.New("unknown VCL config format")
		}

		vcl, err := parseVCLFromMap(vclMap)
		if err != nil {
			return err
		}

		v.VCLs = append(v.VCLs, vcl)
	}

	return nil
}

func parseVCLFromMap(vclMap map[string]interface{}) (VCLConfig, error) {
	vcl := VCLConfig{}
	var ok bool
	vcl.Name, ok = vclMap["name"].(string)
	if !ok {
		return vcl, errors.New("unknown name format")
	}
	vcl.Status, ok = vclMap["status"].(string)
	if !ok {
		return vcl, errors.New("unknown status format")
	}
	vcl.Temperature, ok = vclMap["temperature"].(string)
	if !ok {
		return vcl, errors.New("unknown temperature format")
	}
	vcl.State, ok = vclMap["state"].(string)
	if !ok {
		return vcl, errors.New("unknown state format")
	}
	var busy float64
	busy, ok = vclMap["busy"].(float64)
	if !ok {
		return vcl, errors.New("unknown busy format")
	}
	vcl.Busy = int(busy)

	_, exists := vclMap["label"]
	if exists {
		labelMap, ok := vclMap["label"].(map[string]interface{})
		if !ok {
			return vcl, errors.New("unknown label format")
		}

		referencedLabelName, ok := labelMap["name"].(string)
		if !ok {
			return vcl, errors.New("unknown referenced label format")
		}

		vcl.Label = true
		vcl.ReferencedVCL = &referencedLabelName
	}

	return vcl, nil
}
