package controller

import (
	"reflect"
	"testing"
)

func TestParseConfigs(t *testing.T) {
	createStrRef := func(str string) *string { return &str }

	cases := []struct {
		input    string
		expected []VCLConfig
	}{
		{
			input: `
available   cold/cold          0 boot
active      auto/warm          0 v55329

`,
			expected: []VCLConfig{
				{
					Status:      VCLStatusAvailable,
					Name:        "boot",
					Temperature: VCLTemperatureCold,
				},
				{
					Status:      VCLStatusActive,
					Name:        "v55329",
					Temperature: VCLTemperatureWarm,
				},
			},
		},
		{
			input: `
available   cold/cold          0 boot
active      auto/warm          0 v55329 (1 label)
available  label/warm          0 lable1 -> v55329

`,
			expected: []VCLConfig{
				{
					Status:      VCLStatusAvailable,
					Name:        "boot",
					Temperature: VCLTemperatureCold,
				},
				{
					Status:      VCLStatusActive,
					Name:        "v55329",
					Temperature: VCLTemperatureWarm,
				},
				{
					Status:        VCLStatusAvailable,
					Name:          "lable1",
					Temperature:   VCLTemperatureWarm,
					Label:         true,
					ReferencedVCL: createStrRef("v55329"),
				},
			},
		},
	}

	for _, c := range cases {
		actual := parseVCLConfigsList([]byte(c.input))
		if !reflect.DeepEqual(actual, c.expected) {
			t.Logf(`
Input: 
%s
Parsed config:
%#v
Expected config:
%#v
`, c.input, actual, c.expected)
			t.Fail()
		}
	}
}
