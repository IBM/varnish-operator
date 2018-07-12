package config

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/caarlos0/env"
	"github.com/juju/errors"
	v1 "k8s.io/api/core/v1"
)

// Config describes constant values that will be applied to all varnish services, but may change per-cluster
type Config struct {
	ImagePullSecret           string           `env:"IMAGE_PULL_SECRET,required"`
	OperatorRetryCount        int              `env:"OPERATOR_RETRY_COUNT" envDefault:"3"`
	RestartPolicy             v1.RestartPolicy `env:"RESTART_POLICY" envDefault:"Always"`
	VarnishExporterPort       int32            `env:"VARNISH_EXPORTER_PORT" envDefault:"2034"`
	VarnishExporterTargetPort int32            `env:"VARNISH_EXPORTER_TARGET_PORT" envDefault:"2034"`
	VarnishPort               int32            `env:"VARNISH_PORT" envDefault:"2035"`
	VarnishTargetPort         int32            `env:"VARNISH_TARGET_PORT" envDefault:"2035"`
	VarnishImageHost          string           `env:"VARNISH_IMAGE_HOST,required"`
	VarnishImageNamespace     string           `env:"VARNISH_IMAGE_NAMESPACE" envDefault:"icm-varnish"`
	VarnishImageName          string           `env:"VARNISH_IMAGE_NAME" envDefault:"varnish"`
	VarnishImageTag           string           `env:"VARNISH_IMAGE_TAG,required"`
	VarnishImagePullPolicy    v1.PullPolicy    `env:"VARNISH_IMAGE_PULL_POLICY" envDefault:"Always"`
	VarnishName               string           `env:"VARNISH_NAME" envDefault:"varnish"`
	VCLDir                    string           `env:"VCL_DIR" envDefault:"/etc/varnish"`
	VarnishImageFullPath      string
	VarnishExporterName       string
	VarnishCommonLabels map[string]string
}

func verifyImagePullPolicy(v v1.PullPolicy) error {
	switch v {
	case v1.PullAlways:
		return nil
	case v1.PullNever:
		return nil
	case v1.PullIfNotPresent:
		return nil
	default:
		return errors.NotSupportedf("ImagePullPolicy %s not supported", v)
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
		return errors.NotSupportedf("RestartPolicy %s not supported", v)
	}
}

func int32Parser(v string) (interface{}, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, errors.NotSupportedf("%s is not an int32", v)
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
		return &c, errors.Trace(err)
	}
	if err := verifyImagePullPolicy(c.VarnishImagePullPolicy); err != nil {
		return &c, errors.Trace(err)
	}
	if err := verifyRestartPolicy(c.RestartPolicy); err != nil {
		return &c, errors.Trace(err)
	}
	c.VarnishImageFullPath = c.fullImagePath()
	c.VarnishExporterName = fmt.Sprintf("%s-exporter", c.VarnishName)
	c.VarnishCommonLabels = map[string]string{
		"owner": c.VarnishName,
	}
	return &c, nil
}

// FullImagePath compiles the full path to the image
func (c *Config) fullImagePath() string {
	return fmt.Sprintf("%s/%s/%s:%s", c.VarnishImageHost, c.VarnishImageNamespace, c.VarnishImageName, c.VarnishImageTag)
}
