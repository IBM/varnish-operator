package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	LabelVarnishOwner     = "varnish-owner"
	LabelVarnishComponent = "varnish-component"
	LabelVarnishUID       = "varnish-uid"

	VarnishComponentVarnishes           = "varnishes"
	VarnishComponentCacheService        = "cache-service"
	VarnishComponentNoCacheService      = "nocache-service"
	VarnishComponentClusterRole         = "clusterrole"
	VarnishComponentClusterRoleBinding  = "clusterrolebinding"
	VarnishComponentRole                = "role"
	VarnishComponentRoleBinding         = "rolebinding"
	VarnishComponentVCLFileConfigMap    = "vcl-file-configmap"
	VarnishComponentPodDisruptionBudget = "poddisruptionbudget"
	VarnishComponentServiceAccount      = "serviceaccount"

	VarnishPort                   = 6081
	VarnishAdminPort              = 6082
	VarnishPrometheusExporterPort = 9131
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishService is the Schema for the varnishservices API
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".status.deployment.replicas",description="desired number of varnish pods",format="int32",priority="0"
// +kubebuilder:printcolumn:name="Current",type="integer",JSONPath=".status.deployment.readyReplicas",description="current number of varnish pods",format="int32",priority="0"
// +kubebuilder:printcolumn:name="Up-To-Date",type="integer",JSONPath=".status.deployment.updatedReplicas",description="number of varnish pods that are up to date",format="int32",priority="0"
// +kubebuilder:printcolumn:name="VCL-Version",type="string",JSONPath=".status.vcl.configMapVersion",description="version of the ConfigMap containing the VCL",priority="0"
// +kubebuilder:printcolumn:name="Service-IP",type="string",JSONPath=".status.service.ip",description="IP Address of the service backed by Varnish",priority="0"
// +kubebuilder:printcolumn:name="ServiceNoCache-IP",type="string",JSONPath=".status.serviceNoCache.ip",description="IP Address of the service without any caching",priority="0"
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.deployment.replicas,statuspath=.status.deployment.replicas,selectorpath=
type VarnishService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VarnishServiceSpec   `json:"spec,omitempty"`
	Status VarnishServiceStatus `json:"status,omitempty"`
}

// VarnishServiceSpec defines the desired state of VarnishService
// Important: Run "make" to regenerate code after modifying this file
type VarnishServiceSpec struct {
	VCLConfigMap        VarnishVCLConfigMap                    `json:"vclConfigMap"`
	Deployment          VarnishDeployment                      `json:"deployment,omitempty"`
	PodDisruptionBudget *policyv1beta1.PodDisruptionBudgetSpec `json:"podDisruptionBudget,omitempty"`
	Service             VarnishServiceService                  `json:"service"`
	LogLevel            string                                 `json:"logLevel,omitempty"`
	LogFormat           string                                 `json:"logFormat,omitempty"`
}

type VarnishVCLConfigMap struct {
	Name           string `json:"name"`
	EntrypointFile string `json:"entrypointFile"`
}

type VarnishDeployment struct {
	Replicas    *int32           `json:"replicas,omitempty"`
	Container   VarnishContainer `json:"container,omitempty"`
	Affinity    *v1.Affinity     `json:"affinity,omitempty"`
	Tolerations []v1.Toleration  `json:"tolerations,omitempty"`
}

type VarnishContainer struct {
	Image           string                  `json:"image,omitempty"`
	ImagePullPolicy v1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	RestartPolicy   v1.RestartPolicy        `json:"restartPolicy,omitempty"`
	Resources       v1.ResourceRequirements `json:"resources,omitempty"`
	ImagePullSecret *string                 `json:"imagePullSecret,omitempty"`
	VarnishArgs     []string                `json:"varnishArgs,omitempty"`
}

type VarnishServiceService struct {
	v1.ServiceSpec
	VarnishPort           v1.ServicePort `json:"varnishPort"`
	VarnishExporterPort   v1.ServicePort `json:"varnishExporterPort"`
	PrometheusAnnotations bool           `json:"prometheusAnnotations,omitempty"`
}

// VarnishServiceStatus defines the observed state of VarnishService
type VarnishServiceStatus struct {
	Deployment     VarnishServiceDeploymentStatus `json:"deployment,omitempty"`
	Service        VarnishServiceServiceStatus    `json:"service,omitempty"`
	ServiceNoCache VarnishServiceServiceStatus    `json:"serviceNoCache,omitempty"`
	VCL            VCLStatus                      `json:"vcl"`
}

// VCLStatus describes the VCL versions status
type VCLStatus struct {
	Version          *string `json:"version,omitempty"`
	ConfigMapVersion string  `json:"configMapVersion"`
	Availability     string  `json:"availability"`
}

// VarnishServiceSingleServiceStatus describes the status of a service as it exists within a VarnishService
type VarnishServiceServiceStatus struct {
	v1.ServiceStatus `json:",inline"`
	Name             string `json:"name,omitempty"`
	IP               string `json:"ip,omitempty"`
}

// VarnishServiceDeploymentStatus comprises information about the deployment and its status.
type VarnishServiceDeploymentStatus struct {
	Name        string `json:"name,omitempty"`
	VarnishArgs string `json:"varnishArgs,omitempty"`
	appsv1.DeploymentStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishServiceList contains a list of VarnishService
type VarnishServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VarnishService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VarnishService{}, &VarnishServiceList{})
}
