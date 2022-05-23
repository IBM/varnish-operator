package v1alpha1

// +kubebuilder:validation:Optional

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	LabelVarnishOwner     = "varnish-owner"
	LabelVarnishComponent = "varnish-component"
	LabelVarnishUID       = "varnish-uid"

	VarnishComponentVarnish                  = "varnish"
	VarnishComponentCacheService             = "cache-service"
	VarnishComponentNoCacheService           = "nocache-service"
	VarnishComponentClusterRole              = "clusterrole"
	VarnishComponentClusterRoleBinding       = "clusterrolebinding"
	VarnishComponentRole                     = "role"
	VarnishComponentRoleBinding              = "rolebinding"
	VarnishComponentVCLFileConfigMap         = "vcl-file-configmap"
	VarnishComponentPodDisruptionBudget      = "poddisruptionbudget"
	VarnishComponentHeadlessService          = "headless-service"
	VarnishComponentServiceAccount           = "serviceaccount"
	VarnishComponentValidatingWebhook        = "validating-webhook"
	VarnishComponentMutatingWebhook          = "mutating-webhook"
	VarnishComponentSecret                   = "secret"
	VarnishComponentPrometheusServiceMonitor = "prometheus-servicemonitor"
	VarnishComponentGrafanaDashboard         = "grafana-dashboard"

	VarnishPort                   = 6081
	VarnishAdminPort              = 6082
	VarnishPrometheusExporterPort = 9131
	VarnishControllerMetricsPort  = 8235
	HealthCheckPort               = 8234

	VarnishContainerName             = "varnish"
	VarnishMetricsExporterName       = "metrics-exporter"
	VarnishMetricsExporterImage      = "-metrics-exporter"
	VarnishControllerName            = "varnish-controller"
	VarnishControllerImage           = "-controller"
	VarnishControllerMetricsPortName = "ctrl-metrics"
	VarnishMetricsPortName           = "metrics"
	VarnishPortName                  = "varnish"
	VarnishSharedVolume              = "workdir"
	VarnishSettingsVolume            = "settings"
	VarnishSecretVolume              = "secret"

	VarnishUpdateStrategyDelayedRollingUpdate = "DelayedRollingUpdate"

	VarnishClusterBackendZoneBalancingTypeDisabled   = "disabled"
	VarnishClusterBackendZoneBalancingTypeAuto       = "auto"
	VarnishClusterBackendZoneBalancingTypeThresholds = "thresholds"

	HaproxyContainerName   = "haproxy-sidecar"
	HaproxyConfigFileName  = "haproxy.cfg"
	HaproxyConfigDir       = "/usr/local/etc/haproxy"
	HaproxyConfigMapName   = "haproxy-configmap"
	HaproxyConfigVolume    = "haproxy-config"
	HaproxyMetricsPort     = 8404
	HaproxyMetricsPortName = "haproxy-metrics"
	HaproxyScriptsVolume   = "haproxy-scripts"
)

// +kubebuilder:object:root=true
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishCluster is the Schema for the varnishclusters API
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced,shortName=vc
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.varnishPodsSelector
type VarnishCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   VarnishClusterSpec   `json:"spec"`
	Status VarnishClusterStatus `json:"status,omitempty"`
}

// VarnishClusterSpec defines the desired state of VarnishCluster
type VarnishClusterSpec struct {
	Replicas       *int32                        `json:"replicas,omitempty"`
	UpdateStrategy *VarnishClusterUpdateStrategy `json:"updateStrategy,omitempty"`
	Varnish        *VarnishClusterVarnish        `json:"varnish,omitempty"`
	// +kubebuilder:validation:Required
	VCL *VarnishClusterVCL `json:"vcl,omitempty"`
	// +kubebuilder:validation:Required
	Backend *VarnishClusterBackend `json:"backend,omitempty"`
	// +kubebuilder:validation:Required
	Service             *VarnishClusterService                 `json:"service,omitempty"`
	PodDisruptionBudget *policyv1beta1.PodDisruptionBudgetSpec `json:"podDisruptionBudget,omitempty"`
	NodeSelector        map[string]string                      `json:"nodeSelector,omitempty"`
	Affinity            *v1.Affinity                           `json:"affinity,omitempty"`
	Tolerations         []v1.Toleration                        `json:"tolerations,omitempty"`
	Monitoring          *VarnishClusterMonitoring              `json:"monitoring,omitempty"`
	// +kubebuilder:validation:Enum=debug;info;warn;error;dpanic;panic;fatal
	LogLevel string `json:"logLevel,omitempty"`
	// +kubebuilder:validation:Enum=json;console
	LogFormat      string          `json:"logFormat,omitempty"`
	HaproxySidecar *HaproxySidecar `json:"haproxySidecar,omitempty"`
}

type HaproxySidecar struct {
	Enabled bool `json:"enabled,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9.-]+$`
	ConfigMapName string `json:"configMapName,omitempty"` // mount under /usr/local/etc/haproxy/haproxy.cfg
	Image         string `json:"image,omitempty"`
	// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
	ImagePullPolicy           v1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	ImagePullSecret           string                  `json:"imagePullSecret,omitempty"`
	Resources                 v1.ResourceRequirements `json:"resources,omitempty"`
	MaxConnections            *int32                  `json:"maxConnections,omitempty"`
	ConnectTimeout            *int32                  `json:"connectTimeout,omitempty"`  // in millis, 5000 default
	ClientTimeout             *int32                  `json:"clientTimeout,omitempty"`   // in millis, 50000 default
	ServerTimeout             *int32                  `json:"serverTimeout,omitempty"`   // in millis, 50000 default
	StatRefreshRate           *int32                  `json:"statRefreshRate,omitempty"` // in seconds, 10 default
	BackendServerPort         *int32                  `json:"backendServerPort,omitempty"`
	BackendServerHostHeader   string                  `json:"backendServerHostHeader,omitempty"`
	BackendServerMaxAgeHeader *int32                  `json:"backendServerMaxAgeHeader,omitempty"`
	BackendServers            []string                `json:"backendServers,omitempty"`
}

type VarnishClusterUpdateStrategyType string

const (
	OnDeleteVarnishClusterStrategyType             = VarnishClusterUpdateStrategyType(appsv1.OnDeleteStatefulSetStrategyType)
	RollingUpdateVarnishClusterStrategyType        = VarnishClusterUpdateStrategyType(appsv1.RollingUpdateStatefulSetStrategyType)
	DelayedRollingUpdateVarnishClusterStrategyType = VarnishClusterUpdateStrategyType("DelayedRollingUpdate")
)

type VarnishClusterUpdateStrategy struct {
	// +kubebuilder:validation:Enum=OnDelete;RollingUpdate;DelayedRollingUpdate
	Type                 VarnishClusterUpdateStrategyType         `json:"type,omitempty"`
	RollingUpdate        *appsv1.RollingUpdateStatefulSetStrategy `json:"rollingUpdate,omitempty"`
	DelayedRollingUpdate *UpdateStrategyDelayedRollingUpdate      `json:"delayedRollingUpdate,omitempty"`
}

type UpdateStrategyDelayedRollingUpdate struct {
	DelaySeconds int32 `json:"delaySeconds,omitempty"`
}

type VarnishClusterVarnish struct {
	Image string `json:"image,omitempty"`
	// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
	ImagePullPolicy           v1.PullPolicy                         `json:"imagePullPolicy,omitempty"`
	ImagePullSecret           string                                `json:"imagePullSecret,omitempty"`
	Resources                 *v1.ResourceRequirements              `json:"resources,omitempty"`
	Args                      []string                              `json:"args,omitempty"`
	Controller                *VarnishClusterVarnishController      `json:"controller,omitempty"`
	MetricsExporter           *VarnishClusterVarnishMetricsExporter `json:"metricsExporter,omitempty"`
	Secret                    *VarnishClusterVarnishSecret          `json:"admAuth,omitempty"`
	EnvFrom                   []v1.EnvFromSource                    `json:"envFrom,omitempty"`
	ExtraInitContainers       []v1.Container                        `json:"extraInitContainers,omitempty"`
	ExtraVolumeClaimTemplates []PVC                                 `json:"extraVolumeClaimTemplates,omitempty"`
	ExtraVolumes              []v1.Volume                           `json:"extraVolumes,omitempty"`
	ExtraVolumeMounts         []v1.VolumeMount                      `json:"extraVolumeMounts,omitempty"`
}

type PVC struct {
	Metadata ObjectMetadata `json:"metadata,omitempty"`
	// +kubebuilder:validation:Required
	Spec v1.PersistentVolumeClaimSpec `json:"spec"`
}

type ObjectMetadata struct {
	Name        string            `json:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type VarnishClusterVarnishController struct {
	Image string `json:"image,omitempty"`
	// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
	ImagePullPolicy v1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	Resources       v1.ResourceRequirements `json:"resources,omitempty"`
}

type VarnishClusterVarnishMetricsExporter struct {
	Image string `json:"image,omitempty"`
	// +kubebuilder:validation:Enum=Always;Never;IfNotPresent
	ImagePullPolicy v1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	Resources       v1.ResourceRequirements `json:"resources,omitempty"`
}

type VarnishClusterVCL struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9.-]+$`
	ConfigMapName *string `json:"configMapName,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^.+\.vcl$
	EntrypointFileName *string `json:"entrypointFileName,omitempty"`
}

// Defines the type and parameters for backend traffic distribution
// in multi-zone clusters
type VarnishClusterBackendZoneBalancing struct {
	// +kubebuilder:validation:Enum=auto;thresholds;disabled
	Type       string                                        `json:"type,omitempty"`
	Thresholds []VarnishClusterBackendZoneBalancingThreshold `json:"thresholds,omitempty"`
}

// Defines one or more conditions and respective weights for backends
// located in the same or remote zone
type VarnishClusterBackendZoneBalancingThreshold struct {
	// +kubebuilder:validation:Required
	Local *int `json:"local"`
	// +kubebuilder:validation:Required
	Remote *int `json:"remote"`
	// +kubebuilder:validation:Required
	Threshold *int `json:"threshold"`
}

type VarnishClusterBackend struct {
	// +kubebuilder:validation:Required
	Selector map[string]string `json:"selector,omitempty"`
	// +kubebuilder:validation:Required
	Port          *intstr.IntOrString                 `json:"port,omitempty"`
	Namespaces    []string                            `json:"namespaces,omitempty"`
	OnlyReady     bool                                `json:"onlyReady,omitempty"`
	ZoneBalancing *VarnishClusterBackendZoneBalancing `json:"zoneBalancing,omitempty"`
}

type VarnishClusterVarnishSecret struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9.-]+$`
	SecretName *string `json:"secretName,omitempty"`
	//+kubebuilder:validation:Pattern=`^[a-zA-Z0-9._-]+$`
	Key *string `json:"key,omitempty"`
}

type VarnishClusterService struct {
	// +kubebuilder:validation:Required
	Port        *int32 `json:"port,omitempty"`
	MetricsPort *int32 `json:"metricsPort,omitempty"`
	// +kubebuilder:validation:Enum=ClusterIP;LoadBalancer;NodePort
	Type        v1.ServiceType    `json:"type,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type VarnishClusterMonitoring struct {
	PrometheusServiceMonitor *VarnishClusterMonitoringPrometheusServiceMonitor `json:"prometheusServiceMonitor,omitempty"`
	GrafanaDashboard         *VarnishClusterMonitoringGrafanaDashboard         `json:"grafanaDashboard,omitempty"`
}

type VarnishClusterMonitoringPrometheusServiceMonitor struct {
	Enabled        bool              `json:"enabled"`
	Namespace      string            `json:"namespace"`
	Labels         map[string]string `json:"labels,omitempty"`
	ScrapeInterval string            `json:"scrapeInterval,omitempty"`
}

type VarnishClusterMonitoringGrafanaDashboard struct {
	Enabled   bool              `json:"enabled"`
	Title     string            `json:"title,omitempty"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels,omitempty"`
	// +kubebuilder:validation:Required
	DatasourceName *string `json:"datasourceName,omitempty"`
}

// VarnishClusterStatus defines the observed state of VarnishCluster
type VarnishClusterStatus struct {
	VCL                 VCLStatus            `json:"vcl"`
	HAProxy             HaproxySidecarStatus `json:"haproxy"`
	VarnishArgs         string               `json:"varnishArgs,omitempty"`
	Replicas            int32                `json:"replicas,omitempty"`
	VarnishPodsSelector string               `json:"varnishPodsSelector,omitempty"`
}

type ConfigMapStatus struct {
	Version          *string `json:"version,omitempty"`
	ConfigMapVersion string  `json:"configMapVersion"`
	Availability     string  `json:"availability"`
}

// VCLStatus describes the VCL versions status
type VCLStatus struct {
	ConfigMapStatus `json:",inline"`
}

type HaproxySidecarStatus struct {
	ConfigMapStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishClusterList contains a list of VarnishCluster
type VarnishClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VarnishCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VarnishCluster{}, &VarnishClusterList{})
}
