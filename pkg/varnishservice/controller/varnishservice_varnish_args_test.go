package controller

import (
	"fmt"
	"icm-varnish-k8s-operator/api/v1alpha1"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSanitizedVarnishArgs(t *testing.T) {
	vclConfigMap := v1alpha1.VarnishVCLConfigMap{
		Name:           "vcl-files",
		EntrypointFile: "entry-point-file.vcl",
	}

	cases := []struct {
		name           string
		spec           *v1alpha1.VarnishServiceSpec
		expectedResult []string
	}{
		{
			name: "no specified args",
			spec: &v1alpha1.VarnishServiceSpec{
				VCLConfigMap: vclConfigMap,
				StatefulSet: v1alpha1.VarnishStatefulSet{
					Container: v1alpha1.VarnishContainer{
						VarnishArgs: nil,
					},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-f", "/etc/varnish/entry-point-file.vcl",
			},
		},
		{
			name: "flag -n should be stripped",
			spec: &v1alpha1.VarnishServiceSpec{
				VCLConfigMap: vclConfigMap,
				StatefulSet: v1alpha1.VarnishStatefulSet{
					Container: v1alpha1.VarnishContainer{
						VarnishArgs: []string{"-n", "custom/work/directory"},
					},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-f", "/etc/varnish/entry-point-file.vcl",
			},
		},
		{
			name: "not allowed arguments should be overridden",
			spec: &v1alpha1.VarnishServiceSpec{
				VCLConfigMap: vclConfigMap,
				StatefulSet: v1alpha1.VarnishStatefulSet{
					Container: v1alpha1.VarnishContainer{
						VarnishArgs: []string{"-S", "/etc/varnish/newsecret", "-T", "127.0.0.1:4235", "-a", "0.0.0.0:3425", "-f", "/custom-entry-point.vcl"},
					},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-f", "/etc/varnish/entry-point-file.vcl",
			},
		},
		{
			name: "options with the same key should be supported",
			spec: &v1alpha1.VarnishServiceSpec{
				VCLConfigMap: vclConfigMap,
				StatefulSet: v1alpha1.VarnishStatefulSet{
					Container: v1alpha1.VarnishContainer{
						VarnishArgs: []string{"-p", "default_ttl=3600", "-p", "default_grace=3600"},
					},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-f", "/etc/varnish/entry-point-file.vcl",
				"-p", "default_grace=3600",
				"-p", "default_ttl=3600",
			},
		},
		{
			name: "the order of arguments doesn't change the end results",
			spec: &v1alpha1.VarnishServiceSpec{
				VCLConfigMap: vclConfigMap,
				StatefulSet: v1alpha1.VarnishStatefulSet{
					Container: v1alpha1.VarnishContainer{
						VarnishArgs: []string{"-p", "default_grace=3600", "-p", "default_ttl=3600"},
					},
				},
			},
			expectedResult: []string{
				"-F",
				"-S", "/etc/varnish/secret",
				"-T", fmt.Sprintf("127.0.0.1:%d", v1alpha1.VarnishAdminPort),
				"-a", fmt.Sprintf("0.0.0.0:%d", v1alpha1.VarnishPort),
				"-f", "/etc/varnish/entry-point-file.vcl",
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
