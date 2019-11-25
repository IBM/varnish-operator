package controller

import (
	"github.com/gogo/protobuf/proto"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseConfigs(t *testing.T) {
	cases := []struct {
		input       string
		expected    []VCLConfig
		expectedErr error
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
			expectedErr: nil,
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
					ReferencedVCL: proto.String("v55329"),
				},
			},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		actual, err := parseVCLConfigsList([]byte(c.input))
		if !cmp.Equal(err, c.expectedErr) {
			t.Logf("Unexpected error values: %#v. Expected: %#v", err, c.expectedErr)
			t.Fail()
		}
		if !cmp.Equal(actual, c.expected) {
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
