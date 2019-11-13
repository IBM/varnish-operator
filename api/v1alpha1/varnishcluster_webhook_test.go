package v1alpha1

import (
	"testing"
)

func TestValidatingWebhook(t *testing.T) {
	cases := []struct {
		name  string
		vc    *VarnishCluster
		valid bool
	}{
		{
			name: "Valid values",
			vc: &VarnishCluster{
				Spec: VarnishClusterSpec{
					Varnish: VarnishClusterVarnish{
						Args: []string{"-s", "malloc,2048M"},
					},
				},
			},
			valid: true,
		},
		{
			name: "Invalid values",
			vc: &VarnishCluster{
				Spec: VarnishClusterSpec{
					Varnish: VarnishClusterVarnish{
						Args: []string{"-#@$s", "malloc,2048M"},
					},
				},
			},
			valid: false,
		},
		{
			name: "Empty values",
			vc: &VarnishCluster{
				Spec: VarnishClusterSpec{
					Varnish: VarnishClusterVarnish{
						Args: []string{},
					},
				},
			},
			valid: true,
		},
		{
			name: "Key pattern should match the whole string",
			vc: &VarnishCluster{
				Spec: VarnishClusterSpec{
					Varnish: VarnishClusterVarnish{
						Args: []string{"invalid-s-invalid", "malloc,2048M"},
					},
				},
			},
			valid: false,
		},
		{
			name: "Disallowed VarnishClusterVarnish arguments",
			vc: &VarnishCluster{
				Spec: VarnishClusterSpec{
					Varnish: VarnishClusterVarnish{
						Args: []string{"-a", "-f", "-F", "-n", "-S"},
					},
				},
			},
			valid: false,
		},
	}

	for _, c := range cases {
		err := c.vc.ValidateCreate()
		if c.valid != (err == nil) {
			t.Fatalf("Test %q failed for Create: Expected to be valid: %t, Actual error: %#v", c.name, c.valid, err)
		}

		err = c.vc.ValidateUpdate(&VarnishCluster{})
		if c.valid != (err == nil) {
			t.Fatalf("Test %q failed for Update: Expected to be valid: %t, Actual error: %#v", c.name, c.valid, err)
		}

		err = c.vc.ValidateDelete()
		if err != nil {
			t.Fatalf("Test %q failed for Delete: the delete validationg webhook should allow any requests", c.name)
		}
	}
}
