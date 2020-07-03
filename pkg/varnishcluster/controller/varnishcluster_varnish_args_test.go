package controller

import (
	"fmt"
	"github.com/ibm/varnish-operator/api/v1alpha1"
	"testing"

	"github.com/gogo/protobuf/proto"

	"github.com/google/go-cmp/cmp"
)

func TestGetSanitizedVarnishArgs(t *testing.T) {
	vclConfigMap := &v1alpha1.VarnishClusterVCL{
		ConfigMapName:      proto.String("vcl-files"),
		EntrypointFileName: proto.String("entry-point-file.vcl"),
	}

	cases := []struct {
		name           string
		spec           *v1alpha1.VarnishClusterSpec
		expectedResult []string
	}{
		{
			name: "no specified args",
			spec: &v1alpha1.VarnishClusterSpec{
				VCL: vclConfigMap,
				Varnish: &v1alpha1.VarnishClusterVarnish{
					Args: nil,
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish-secret/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-b", "127.0.0.1:0",
			},
		},
		{
			name: "flag -n should be stripped",
			spec: &v1alpha1.VarnishClusterSpec{
				VCL: vclConfigMap,
				Varnish: &v1alpha1.VarnishClusterVarnish{
					Args: []string{"-n", "custom/work/directory"},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish-secret/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-b", "127.0.0.1:0",
			},
		},
		{
			name: "flag -f should be stripped",
			spec: &v1alpha1.VarnishClusterSpec{
				VCL: vclConfigMap,
				Varnish: &v1alpha1.VarnishClusterVarnish{
					Args: []string{"-f", "/etc/varnish/entry-point-file.vcl"},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish-secret/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-b", "127.0.0.1:0",
			},
		},
		{
			name: "not allowed arguments should be overridden",
			spec: &v1alpha1.VarnishClusterSpec{
				VCL: vclConfigMap,
				Varnish: &v1alpha1.VarnishClusterVarnish{
					Args: []string{"-S", "/etc/varnish/newsecret", "-T", "127.0.0.1:4235", "-a", "0.0.0.0:3425", "-b", "127.0.0.1:3456"},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish-secret/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-b", "127.0.0.1:0",
			},
		},
		{
			name: "options with the same key should be supported",
			spec: &v1alpha1.VarnishClusterSpec{
				VCL: vclConfigMap,
				Varnish: &v1alpha1.VarnishClusterVarnish{
					Args: []string{"-p", "default_ttl=3600", "-p", "default_grace=3600"},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish-secret/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-b", "127.0.0.1:0",
				"-p", "default_grace=3600",
				"-p", "default_ttl=3600",
			},
		},
		{
			name: "the order of arguments doesn't change the end results",
			spec: &v1alpha1.VarnishClusterSpec{
				VCL: vclConfigMap,
				Varnish: &v1alpha1.VarnishClusterVarnish{
					Args: []string{"-p", "default_grace=3600", "-p", "default_ttl=3600"},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish-secret/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-b", "127.0.0.1:0",
				"-p", "default_grace=3600",
				"-p", "default_ttl=3600",
			},
		},
	}

	for _, c := range cases {
		actual := getSanitizedVarnishArgs(c.spec)
		if !cmp.Equal(c.expectedResult, actual) {
			t.Logf("Test %q failed.\nDiff: \n%#v\n%#v", c.name, c.expectedResult, actual)
			t.Fail()
		}
	}
}
