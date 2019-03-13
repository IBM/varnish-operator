package config

import (
	"reflect"
	"strconv"

	"k8s.io/apimachinery/pkg/labels"

	"go.uber.org/zap/zapcore"

	"k8s.io/apimachinery/pkg/types"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

// Config that reads in env variables
type Config struct {
	EndpointSelectorString string        `env:"ENDPOINT_SELECTOR_STRING,required"`
	ConfigMapName          string        `env:"CONFIGMAP_NAME,required"`
	Namespace              string        `env:"NAMESPACE,required"`
	PodName                string        `env:"POD_NAME,required"`
	VarnishServiceName     string        `env:"VARNISH_SERVICE_NAME,required"`
	VarnishServiceUID      types.UID     `env:"VARNISH_SERVICE_UID,required"`
	VarnishServiceGroup    string        `env:"VARNISH_SERVICE_GROUP,required"`
	VarnishServiceVersion  string        `env:"VARNISH_SERVICE_VERSION,required"`
	VarnishServiceKind     string        `env:"VARNISH_SERVICE_KIND,required"`
	LogFormat              string        `env:"LOG_FORMAT,required"`
	LogLevel               zapcore.Level `env:"LOG_LEVEL,required"`
	VCLDir                 string
	EndpointSelector       labels.Selector
}

// Load uses the caarlos0/env library to read in environment variables into a struct
func Load() (*Config, error) {
	c := Config{}
	int32Type := reflect.TypeOf(int32(0))
	int32Parse := env.ParserFunc(func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, errors.Errorf("%s is not an int32", v)
		}
		return int32(i), nil
	})

	zapcoreLevelType := reflect.TypeOf(zapcore.InfoLevel)
	zapcoreLevelParse := env.ParserFunc(func(v string) (interface{}, error) {
		var l zapcore.Level
		err := l.UnmarshalText([]byte(v))
		return l, errors.Wrapf(err, "%s is not a zap level", v)
	})

	parsers := env.CustomParsers{
		int32Type:        int32Parse,
		zapcoreLevelType: zapcoreLevelParse,
	}

	var err error
	if err = env.ParseWithFuncs(&c, parsers); err != nil {
		return &c, errors.WithStack(err)
	}

	// hard-coded here for convenience
	c.VCLDir = "/etc/varnish"

	c.EndpointSelector, err = labels.Parse(c.EndpointSelectorString)
	if err != nil {
		return &c, errors.Wrapf(err, "could not parse endpoint selector: %s", c.EndpointSelectorString)
	}

	return &c, nil
}
