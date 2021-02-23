package config

import (
	"github.com/caarlos0/env/v6"
	"reflect"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const VCLConfigDir = "/etc/varnish"

// Config that reads in env variables
type Config struct {
	EndpointSelectorString string        `env:"ENDPOINT_SELECTOR_STRING,required"`
	Namespace              string        `env:"NAMESPACE,required"`
	PodName                string        `env:"POD_NAME,required"`
	NodeName               string        `env:"NODE_NAME,required"`
	VarnishClusterName     string        `env:"VARNISH_CLUSTER_NAME,required"`
	VarnishClusterUID      types.UID     `env:"VARNISH_CLUSTER_UID,required"`
	VarnishClusterGroup    string        `env:"VARNISH_CLUSTER_GROUP,required"`
	VarnishClusterVersion  string        `env:"VARNISH_CLUSTER_VERSION,required"`
	VarnishClusterKind     string        `env:"VARNISH_CLUSTER_KIND,required"`
	VarnishAdmArgs         []string      `env:"VARNISHADM_ARGS" envDefault:"-S /etc/varnish-secret/secret -T 127.0.0.1:6082" envSeparator:" " `
	VarnishPingTimeout     time.Duration `env:"VARNISHADM_PING_TIMEOUT" envDefault:"90s"`
	VarnishPingDelay       time.Duration `env:"VARNISHADM_PING_DELAY" envDefault:"200ms"`
	LogFormat              string        `env:"LOG_FORMAT,required"`
	LogLevel               zapcore.Level `env:"LOG_LEVEL,required"`
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

	parsers := map[reflect.Type]env.ParserFunc{
		int32Type:        int32Parse,
		zapcoreLevelType: zapcoreLevelParse,
	}

	var err error
	if err = env.ParseWithFuncs(&c, parsers); err != nil {
		return &c, errors.WithStack(err)
	}

	c.EndpointSelector, err = labels.Parse(c.EndpointSelectorString)
	if err != nil {
		return &c, errors.Wrapf(err, "could not parse endpoint selector: %s", c.EndpointSelectorString)
	}

	return &c, nil
}
