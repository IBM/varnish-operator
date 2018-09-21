package config

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Config describes constant values that will be applied to all varnish services, but may change per-cluster
type Config struct {
	VarnishImageHost               string           `env:"VARNISH_IMAGE_HOST" envDefault:"registry.ng.bluemix.net"`
	VarnishImageNamespace          string           `env:"VARNISH_IMAGE_NAMESPACE" envDefault:"icm-varnish"`
	VarnishImageName               string           `env:"VARNISH_IMAGE_NAME" envDefault:"varnish"`
	VarnishImageTag                string           `env:"VARNISH_IMAGE_TAG,required"`
	VarnishImagePullPolicy         v1.PullPolicy    `env:"VARNISH_IMAGE_PULL_POLICY" envDefault:"Always"`
	ImagePullSecret                string           `env:"IMAGE_PULL_SECRET" envDefault:"docker-reg-secret"`
	VarnishExporterPort            int32            `env:"VARNISH_EXPORTER_PORT" envDefault:"9131"`
	VarnishExporterTargetPort      int32            `env:"VARNISH_EXPORTER_TARGET_PORT" envDefault:"9131"`
	VarnishPort                    int32            `env:"VARNISH_PORT" envDefault:"2035"`
	VarnishTargetPort              int              `env:"VARNISH_TARGET_PORT" envDefault:"2035"`
	VarnishName                    string           `env:"VARNISH_NAME" envDefault:"varnish"`
	VCLDir                         string           `env:"VCL_DIR" envDefault:"/etc/varnish"`
	DefaultVarnishMemory           string           `env:"DEFAULT_VARNISH_MEMORY" envDefault:"1024M"`
	DefaultBackendsFile            string           `env:"DEFAULT_BACKENDS_FILE" envDefault:"backends.vcl"`
	DefaultDefaultFile             string           `env:"DEFAULT_DEFAULT_FILE" envDefault:"default.vcl"`
	DefaultVarnishResourceLimitCPU string           `env:"DEFAULT_VARNISH_RESOURCE_LIMIT_CPU" envDefault:"1"`
	DefaultVarnishResourceLimitMem string           `env:"DEFAULT_VARNISH_RESOURCE_LIMIT_MEM" envDefault:"2048Mi"`
	DefaultVarnishResourceReqCPU   string           `env:"DEFAULT_VARNISH_RESOURCE_REQ_CPU" envDefault:"1"`
	DefaultVarnishResourceReqMem   string           `env:"DEFAULT_VARNISH_RESOURCE_REQ_MEM" envDefault:"2048Mi"`
	DefaultVarnishRestartPolicy    v1.RestartPolicy `env:"DEFAULT_VARNISH_RESTART_POLICY" envDefault:"Always"`
	DefaultLivenessProbeHTTPPath   string           `env:"DEFAULT_LIVENESS_PROBE_HTTP_PATH"`
	DefaultLivenessProbePort       int              `env:"DEFAULT_LIVENESS_PROBE_PORT"`
	DefaultReadinessProbeCommand   []string         `env:"DEFAULT_READINESS_PROBE_COMMAND" envDefault:"/usr/bin/varnishadm,ping"`
	VarnishImageFullPath           string
	VarnishCommonLabels            map[string]string
	DefaultVarnishResources        v1.ResourceRequirements
	DefaultLivenessProbe           *v1.Probe
	DefaultReadinessProbe          v1.Probe
}

// GlobalConf is config that affects the operator directly, as well as provides default values for varnish instances
var GlobalConf *Config

func verifyImagePullPolicy(v v1.PullPolicy) error {
	switch v {
	case v1.PullAlways:
		return nil
	case v1.PullNever:
		return nil
	case v1.PullIfNotPresent:
		return nil
	default:
		return errors.Errorf("ImagePullPolicy %s not supported", v)
	}
}

func verifyRestartPolicy(v v1.RestartPolicy) error {
	switch v {
	case v1.RestartPolicyAlways:
		return nil
	case v1.RestartPolicyNever:
		return nil
	case v1.RestartPolicyOnFailure:
		return nil
	default:
		return errors.Errorf("RestartPolicy %s not supported", v)
	}
}

func int32Parser(v string) (interface{}, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, errors.Errorf("%s is not an int32", v)
	}
	return int32(i), nil
}

var (
	int32Type    = reflect.TypeOf(int32(0))
	parseFuncMap = env.CustomParsers{
		int32Type: int32Parser,
	}
)

// LoadConfig uses the env library to read in environment variables and return an instance of Config
func LoadConfig() (*Config, error) {
	c := Config{}
	if err := env.ParseWithFuncs(&c, parseFuncMap); err != nil {
		return &c, errors.WithStack(err)
	}
	if err := verifyImagePullPolicy(c.VarnishImagePullPolicy); err != nil {
		return &c, errors.WithStack(err)
	}
	if err := verifyRestartPolicy(c.DefaultVarnishRestartPolicy); err != nil {
		return &c, errors.WithStack(err)
	}
	c.VarnishImageFullPath = c.fullImagePath()
	c.VarnishCommonLabels = map[string]string{
		"owner": c.VarnishName,
	}

	varnishResourceLimitCPU, err := resource.ParseQuantity(c.DefaultVarnishResourceLimitCPU)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	varnishResourceLimitMem, err := resource.ParseQuantity(c.DefaultVarnishResourceLimitMem)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	varnishResourceReqCPU, err := resource.ParseQuantity(c.DefaultVarnishResourceReqCPU)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	varnishResourceReqMem, err := resource.ParseQuantity(c.DefaultVarnishResourceReqMem)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	c.DefaultVarnishResources = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    varnishResourceLimitCPU,
			v1.ResourceMemory: varnishResourceLimitMem,
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    varnishResourceReqCPU,
			v1.ResourceMemory: varnishResourceReqMem,
		},
	}

	if c.DefaultLivenessProbeHTTPPath != "" && c.DefaultLivenessProbePort != 0 {
		c.DefaultLivenessProbe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path: c.DefaultLivenessProbeHTTPPath,
					Port: intstr.FromInt(c.DefaultLivenessProbePort),
				},
			},
		}
	}

	c.DefaultReadinessProbe = v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: c.DefaultReadinessProbeCommand,
			},
		},
	}

	return &c, nil
}

// FullImagePath compiles the full path to the image
func (c *Config) fullImagePath() string {
	return fmt.Sprintf("%s/%s/%s:%s", c.VarnishImageHost, c.VarnishImageNamespace, c.VarnishImageName, c.VarnishImageTag)
}

func init() {
	var err error
	if GlobalConf, err = LoadConfig(); err != nil {
		log.Fatal(err)
	}

}
