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

	VarnishComponentVarnish             = "varnish"
	VarnishComponentCacheService        = "cache-service"
	VarnishComponentNoCacheService      = "nocache-service"
	VarnishComponentClusterRole         = "clusterrole"
	VarnishComponentClusterRoleBinding  = "clusterrolebinding"
	VarnishComponentRole                = "role"
	VarnishComponentRoleBinding         = "rolebinding"
	VarnishComponentVCLFileConfigMap    = "vcl-file-configmap"
	VarnishComponentPodDisruptionBudget = "poddisruptionbudget"
	VarnishComponentHeadlessService     = "headless-service"
	VarnishComponentServiceAccount      = "serviceaccount"
	VarnishComponentValidatingWebhook   = "validating-webhook"
	VarnishComponentMutatingWebhook     = "mutating-webhook"

	VarnishPort                   = 6081
	VarnishAdminPort              = 6082
	VarnishPrometheusExporterPort = 9131

	VarnishContainerName   = "varnish"
	VarnishMetricsPortName = "metrics"
	VarnishPortName        = "varnish"

	VarnishUpdateStrategyDelayedRollingUpdate = "DelayedRollingUpdate"
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
	Affinity            *v1.Affinity                           `json:"affinity,omitempty"`
	Tolerations         []v1.Toleration                        `json:"tolerations,omitempty"`
	// +kubebuilder:validation:Enum=debug;info;warn;error;dpanic;panic;fatal
	LogLevel string `json:"logLevel,omitempty"`
	// +kubebuilder:validation:Enum=json;console
	LogFormat string `json:"logFormat,omitempty"`
}

type VarnishClusterUpdateStrategyType string

const (
	OnDeleteVarnishClusterStrategyType             VarnishClusterUpdateStrategyType = appsv1.OnDeleteStatefulSetStrategyType
	RollingUpdateVarnishClusterStrategyType        VarnishClusterUpdateStrategyType = appsv1.RollingUpdateStatefulSetStrategyType
	DelayedRollingUpdateVarnishClusterStrategyType VarnishClusterUpdateStrategyType = "DelayedRollingUpdate"
)

type VarnishClusterUpdateStrategy struct {
	// +kubebuilder:validation:Enum=OnDelete;RollingUpdate;DelayedRollingUpdate
	Type                 VarnishClusterUpdateStrategyType         `json:"type,omitempty"`
	RollingUpdate        *appsv1.RollingUpdateStatefulSetStrategy `json:"rollingUpdate,omitempty"`
	DelayedRollingUpdate *UpdateStrategyDelayedRollingUpdate      `json:"delayedRollingUpdate,omitempty"`
}

type UpdateStrategyDelayedRollingUpdate struct {
	// +kubebuilder:validation:Minimum=1
	DelaySeconds int32 `json:"delaySeconds,omitempty"`
}

type VarnishClusterVarnish struct {
	Image           string                   `json:"image,omitempty"`
	ImagePullPolicy v1.PullPolicy            `json:"imagePullPolicy,omitempty"`
	RestartPolicy   v1.RestartPolicy         `json:"restartPolicy,omitempty"`
	Resources       *v1.ResourceRequirements `json:"resources,omitempty"`
	ImagePullSecret *string                  `json:"imagePullSecret,omitempty"`
	Args            []string                 `json:"args,omitempty"`
}

type VarnishClusterVCL struct {
	// +kubebuilder:validation:Required
	ConfigMapName *string `json:"configMapName,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^.+\.vcl$
	EntrypointFileName *string `json:"entrypointFileName,omitempty"`
}

type VarnishClusterBackend struct {
	// +kubebuilder:validation:Required
	Selector map[string]string `json:"selector,omitempty"`
	// +kubebuilder:validation:Required
	Port *intstr.IntOrString `json:"port,omitempty"`
}

type VarnishClusterService struct {
	// +kubebuilder:validation:Required
	Port        *int32 `json:"port,omitempty"`
	MetricsPort int32 `json:"metricsPort,omitempty"`
	// +kubebuilder:validation:Enum=ClusterIP;LoadBalancer;NodePort
	Type        v1.ServiceType    `json:"type,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// VarnishClusterStatus defines the observed state of VarnishCluster
type VarnishClusterStatus struct {
	VCL                 VCLStatus `json:"vcl"`
	VarnishArgs         string    `json:"varnishArgs,omitempty"`
	Replicas            int32     `json:"replicas,omitempty"`
	VarnishPodsSelector string    `json:"varnishPodsSelector,omitempty"`
}

// VCLStatus describes the VCL versions status
type VCLStatus struct {
	Version          *string `json:"version,omitempty"`
	ConfigMapVersion string  `json:"configMapVersion"`
	Availability     string  `json:"availability"`
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
