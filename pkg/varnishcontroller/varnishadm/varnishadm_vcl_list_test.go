package varnishadm

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestListCommand(t *testing.T) {
	cases := []struct {
		errExpected error
		execute     executorProvider
		response    []VCLConfig
		desc        string
	}{
		{
			nil,
			func(name string, args ...string) executor {
				if args[len(args)-1] == "-j" {
					return &mockExecutor{err: errors.New("err"), response: []byte("Command failed with error code 102\nJSON unimplemented")}
				}
				return mockSuccesListResponse(name, args...)
			},
			[]VCLConfig{
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
			"success",
		},
		{
			errors.Wrap(errors.New("some error"), "A response from external program"),
			func(name string, args ...string) executor {
				if args[len(args)-1] == "-j" {
					return &mockExecutor{err: errors.New("err"), response: []byte("Command failed with error code 102\nJSON unimplemented")}
				}
				return mockErrResponse(name, args...)
			},
			[]VCLConfig{},
			"error",
		},
		{
			nil,
			func(name string, args ...string) executor {
				return &mockExecutor{response: []byte(`
[ 2, ["vcl.list", "-j"], 1632389190.659,
  {
    "status": "available",
    "state": "auto",
    "temperature": "warm",
    "busy": 0,
    "name": "boot"
},
  {
    "status": "active",
    "state": "auto",
    "temperature": "warm",
    "busy": 0,
    "name": "v-988-1632389040"
},
{
    "status": "available",
    "state": "label",
    "temperature": "warm",
    "busy": 0,
    "name": "labeledvcl",
    "label": {
      "name": "v-988-1632389040"
      }
}
]`)}
			},
			[]VCLConfig{
				{
					Status:      VCLStatusAvailable,
					Name:        "boot",
					State:       "auto",
					Temperature: VCLTemperatureWarm,
					Busy:        0,
				},
				{
					Status:      VCLStatusActive,
					State:       "auto",
					Name:        "v-988-1632389040",
					Temperature: VCLTemperatureWarm,
					Busy:        0,
				},
				{
					Status:        VCLStatusAvailable,
					State:         "label",
					Name:          "labeledvcl",
					Temperature:   VCLTemperatureWarm,
					Label:         true,
					ReferencedVCL: proto.String("v-988-1632389040"),
					Busy:          0,
				},
			},
			"json",
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(tt *testing.T) {
			p := &VarnishAdm{
				execute: tc.execute,
			}
			data, err := p.List()
			if !cmp.Equal(data, tc.response) {
				tt.Errorf("Unexpected response %v\n Expected: %v", data, tc.response)
			}
			if !cmp.Equal(err, tc.errExpected, equalError) {
				tt.Errorf("Unexpected error return. %s", cmp.Diff(err, tc.errExpected))
			}
		})
	}
}

func TestParseConfigs(t *testing.T) {
	cases := []struct {
		input       string
		expected    []VCLConfig
		expectedErr error
	}{
		{
			input: simpleVCLconfig,
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
			input: labeledVCLconfig,
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
					Name:          "label1",
					Temperature:   VCLTemperatureWarm,
					Label:         true,
					ReferencedVCL: proto.String("v55329"),
				},
			},
			expectedErr: nil,
		},
		{
			input:       unknownVCLconfig,
			expected:    nil,
			expectedErr: errors.New("unknown VCL config format"),
		},
	}

	for _, c := range cases {
		actual, err := parseVCLConfigsList([]byte(c.input))
		if !cmp.Equal(err, c.expectedErr, equalError) {
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
