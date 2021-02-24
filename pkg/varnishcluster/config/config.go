package config

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
	dockerref "github.com/docker/distribution/reference"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// Config contains the Varnish Operator configs
type Config struct {
	Namespace             string        `env:"NAMESPACE" envDefault:"default"`
	LeaderElectionEnabled bool          `env:"LEADERELECTION_ENABLED" envDefault:"true"`
	ContainerImage        string        `env:"CONTAINER_IMAGE,required"`
	LogLevel              zapcore.Level `env:"LOGLEVEL" envDefault:"info"`
	LogFormat             string        `env:"LOGFORMAT" envDefault:"json"`
	WebhooksPort          int32         `env:"WEBHOOKS_PORT" envDefault:"7340"`
	MetricsPort           int32         `env:"METRICS_PORT" envDefault:"8329"`
	WebhooksEnabled       bool          `env:"WEBHOOKS_ENABLED" envDefault:"true"`

	CoupledVarnishImage string
}

func int32Parser(v string) (interface{}, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, errors.Errorf("%s is not an int32", v)
	}
	return int32(i), nil
}

func levelParser(v string) (interface{}, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(v))
	if err != nil {
		return nil, errors.Errorf("%s is not an zapcore.Level", v)
	}
	return level, nil
}

var (
	int32Type    = reflect.TypeOf(int32(0))
	levelType    = reflect.TypeOf(zapcore.DebugLevel)
	parseFuncMap = map[reflect.Type]env.ParserFunc{
		int32Type: int32Parser,
		levelType: levelParser,
	}
)

// LoadConfig uses the env library to read in environment variables and return an instance of Config
func LoadConfig() (*Config, error) {
	c := Config{}
	if err := env.ParseWithFuncs(&c, parseFuncMap); err != nil {
		return &c, errors.WithStack(err)
	}

	ref, err := dockerref.Parse(c.ContainerImage)
	if err != nil {
		return &c, errors.Wrap(err, "image is not properly formatted")
	}
	nt, ok := ref.(dockerref.NamedTagged)
	if !ok {
		return &c, errors.New("image name does not include tag")
	}
	repo := nt.Name()
	if idx := strings.LastIndexByte(repo, '/'); idx != -1 {
		repo = repo[:idx] // chop off `/<image-name>`
	}

	varnishImageName, err := dockerref.WithName(repo + "/varnish")
	if err != nil {
		return &c, errors.Wrap(err, "could not initialize varnish image name")
	}
	varnishImage, err := dockerref.WithTag(varnishImageName, nt.Tag())
	if err != nil {
		return &c, errors.Wrap(err, "could not include tag to varnish image name")
	}
	c.CoupledVarnishImage = varnishImage.String()

	return &c, nil
}
