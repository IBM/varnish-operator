package v1alpha1

import (
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestValidatingWebhook(t *testing.T) {
	cases := []struct {
		name  string
		vs    *VarnishService
		valid bool
	}{
		{
			name: "Valid values",
			vs: &VarnishService{
				Spec: VarnishServiceSpec{
					StatefulSet: VarnishStatefulSet{
						Container: VarnishContainer{
							VarnishArgs: []string{"-s", "malloc,2048M"},
						},
					},
				},
			},
			valid: true,
		},
		{
			name: "Invalid values",
			vs: &VarnishService{
				Spec: VarnishServiceSpec{
					StatefulSet: VarnishStatefulSet{
						Container: VarnishContainer{
							VarnishArgs: []string{"-#@$s", "malloc,2048M"},
						},
					},
				},
			},
			valid: false,
		},
		{
			name: "Empty values",
			vs: &VarnishService{
				Spec: VarnishServiceSpec{
					StatefulSet: VarnishStatefulSet{
						Container: VarnishContainer{
							VarnishArgs: []string{},
						},
					},
				},
			},
			valid: true,
		},
		{
			name: "Key pattern should match the whole string",
			vs: &VarnishService{
				Spec: VarnishServiceSpec{
					StatefulSet: VarnishStatefulSet{
						Container: VarnishContainer{
							VarnishArgs: []string{"invalid-s-invalid", "malloc,2048M"},
						},
					},
				},
			},
			valid: false,
		},
		{
			name: "Disallowed Varnish arguments",
			vs: &VarnishService{
				Spec: VarnishServiceSpec{
					StatefulSet: VarnishStatefulSet{
						Container: VarnishContainer{
							VarnishArgs: []string{"-a", "-f", "-F", "-n", "-S"},
						},
					},
				},
			},
			valid: false,
		},
		{
			name: "Duplicate port names are not allowed for Varnish port",
			vs: &VarnishService{
				Spec: VarnishServiceSpec{
					Service: VarnishServiceService{
						VarnishPort: v1.ServicePort{
							Name: "varnish",
						},
						ServiceSpec: v1.ServiceSpec{
							Ports: []v1.ServicePort{
								{
									Name: "varnish",
								},
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			name: "Duplicate port names are not allowed for Varnish exporter port",
			vs: &VarnishService{
				Spec: VarnishServiceSpec{
					Service: VarnishServiceService{
						VarnishExporterPort: v1.ServicePort{
							Name: "varnishexporter",
						},
						ServiceSpec: v1.ServiceSpec{
							Ports: []v1.ServicePort{
								{
									Name: "varnishexporter",
								},
							},
						},
					},
				},
			},
			valid: false,
		},
	}

	for _, c := range cases {
		err := c.vs.ValidateCreate()
		if c.valid != (err == nil) {
			t.Fatalf("Test %q failed for Create: Expected to be valid: %t, Actual error: %#v", c.name, c.valid, err)
		}

		err = c.vs.ValidateUpdate(&VarnishService{})
		if c.valid != (err == nil) {
			t.Fatalf("Test %q failed for Update: Expected to be valid: %t, Actual error: %#v", c.name, c.valid, err)
		}

		err = c.vs.ValidateDelete()
		if err != nil {
			t.Fatalf("Test %q failed for Delete: the delete validationg webhook should allow any requests", c.name)
		}
	}
}
