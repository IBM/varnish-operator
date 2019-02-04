package config

import (
	"reflect"
	"strconv"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Config describes constant values that will be applied to all varnish services, but may change per-cluster
type Config struct {
	Namespace                                       string           `env:"NAMESPACE" envDefault:"varnish-operator-system"`
	OperatorLeaderElectionEnabled                   bool             `env:"OPERATOR_LEADERELECTION_ENABLED" envDefault:"true"`
	OperatorLeaderElectionID                        string           `env:"OPERATOR_LEADERELECTION_ID" envDefault:"varnish-operator-lock"`
	OperatorLogLevel                                zapcore.Level    `env:"OPERATOR_LOGLEVEL" envDefault:"info"`
	OperatorLogFormat                               string           `env:"OPERATOR_LOGFORMAT" envDefault:"console"`
	VclConfigMapBackendsFile                        string           `env:"VCLCONFIGMAP_BACKENDSFILE" envDefault:"backends.vcl"`
	VclConfigMapDefaultFile                         string           `env:"VCLCONFIGMAP_DEFAULTFILE" envDefault:"default.vcl"`
	DeploymentReplicas                              int32            `env:"DEPLOYMENT_REPLICAS" envDefault:"2"`
	DeploymentContainerImage                        string           `env:"DEPLOYMENT_CONTAINER_IMAGE,required"`
	DeploymentContainerImagePullPolicy              v1.PullPolicy    `env:"DEPLOYMENT_CONTAINER_IMAGEPULLPOLICY" envDefault:"Always"`
	DeploymentContainerRestartPolicy                v1.RestartPolicy `env:"DEPLOYMENT_CONTAINER_RESTARTPOLICY" envDefault:"Always"`
	DeploymentContainerResourcesLimitsCpu           string           `env:"DEPLOYMENT_CONTAINER_RESOURCES_LIMITS_CPU" envDefault:"1"`
	DeploymentContainerResourcesLimitsMemory        string           `env:"DEPLOYMENT_CONTAINER_RESOURCES_LIMITS_MEMORY" envDefault:"2048Mi"`
	DeploymentContainerResourcesRequestsCpu         string           `env:"DEPLOYMENT_CONTAINER_RESOURCES_REQUESTS_CPU" envDefault:"1"`
	DeploymentContainerResourcesRequestsMemory      string           `env:"DEPLOYMENT_CONTAINER_RESOURCES_REQUESTS_MEMORY" envDefault:"2048Mi"`
	DeploymentContainerLivenessProbeHTTPGetHTTPPath string           `env:"DEPLOYMENT_CONTAINER_LIVENESSPROBE_HTTPGET_HTTPPATH"`
	DeploymentContainerLivenessProbeHTTPGetPort     int              `env:"DEPLOYMENT_CONTAINER_LIVENESSPROBE_HTTPGET_PORT"`
	DeploymentContainerReadinessProbeExecCommand    []string         `env:"DEPLOYMENT_CONTAINER_READINESSPROBE_EXEC_COMMAND" envDefault:"/usr/bin/varnishadm,ping"`
	DeploymentContainerImagePullSecret              string           `env:"DEPLOYMENT_CONTAINER_IMAGEPULLSECRET" envDefault:"docker-reg-secret"`
	ServicePrometheusAnnotations                    bool             `env:"SERVICE_PROMETHEUSANNOTATIONS" envDefault:"true"`

	DeploymentContainerResources      v1.ResourceRequirements
	DeploymentContainerLivenessProbe  *v1.Probe
	DeploymentContainerReadinessProbe *v1.Probe
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
	parseFuncMap = env.CustomParsers{
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
	if err := verifyImagePullPolicy(c.DeploymentContainerImagePullPolicy); err != nil {
		return &c, errors.WithStack(err)
	}
	if err := verifyRestartPolicy(c.DeploymentContainerRestartPolicy); err != nil {
		return &c, errors.WithStack(err)
	}

	varnishResourceLimitCPU, err := resource.ParseQuantity(c.DeploymentContainerResourcesLimitsCpu)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	varnishResourceLimitMem, err := resource.ParseQuantity(c.DeploymentContainerResourcesLimitsMemory)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	varnishResourceReqCPU, err := resource.ParseQuantity(c.DeploymentContainerResourcesRequestsCpu)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	varnishResourceReqMem, err := resource.ParseQuantity(c.DeploymentContainerResourcesRequestsMemory)
	if err != nil {
		return &c, errors.WithStack(err)
	}
	c.DeploymentContainerResources = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    varnishResourceLimitCPU,
			v1.ResourceMemory: varnishResourceLimitMem,
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    varnishResourceReqCPU,
			v1.ResourceMemory: varnishResourceReqMem,
		},
	}

	if c.DeploymentContainerLivenessProbeHTTPGetHTTPPath != "" && c.DeploymentContainerLivenessProbeHTTPGetPort != 0 {
		c.DeploymentContainerLivenessProbe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path: c.DeploymentContainerLivenessProbeHTTPGetHTTPPath,
					Port: intstr.FromInt(c.DeploymentContainerLivenessProbeHTTPGetPort),
				},
			},
		}
	}

	if len(c.DeploymentContainerReadinessProbeExecCommand) > 0 {
		c.DeploymentContainerReadinessProbe = &v1.Probe{
			Handler: v1.Handler{
				Exec: &v1.ExecAction{
					Command: c.DeploymentContainerReadinessProbeExecCommand,
				},
			},
		}
	}

	return &c, nil
}
