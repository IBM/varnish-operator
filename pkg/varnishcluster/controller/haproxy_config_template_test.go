package controller

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	vcapi "github.com/ibm/varnish-operator/api/v1alpha1"
)

func templatizeHaproxy(haproxySidecar *vcapi.HaproxySidecar) (string, error) {
	t, err := template.New("haproxy-config").Parse(haproxyConfigTemplate)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := t.Execute(&b, *haproxySidecar); err != nil {
		return "", err
	}

	return b.String(), nil
}

// hack to trick intellij into running tests
func TestHaproxyConfigTemplate(t *testing.T) {
	haproxySidecar := &vcapi.HaproxySidecar{
		Enabled:                 true,
		Image:                   "haproxytech/haproxy-debian:2.4",
		ImagePullPolicy:         "IfNotPresent",
		BackendServers:          []string{"api.mapbox.com"},
		BackendServerHostHeader: "api.mapbox.com",
		EnableFrontendMetrics:   true,
		HttpChk:                 []string{"foo", "bar", "baz"},
		BackendAdditionalFlags:  "this is the worst",
	}
	vcapi.DefaultHaproxySidecar(haproxySidecar)
	data, err := templatizeHaproxy(haproxySidecar)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", data)
	}
	//Expect(err).ToNot(HaveOccurred())
	//Expect(data).ToNot(BeNil())
}
