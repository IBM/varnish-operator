package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceWrapper is just a way to generate a path "service.spec"
type ServiceWrapper struct {
	Spec v1.ServiceSpec `json:"spec"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishService is the Schema for the varnishservices API
// +k8s:openapi-gen=true
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
	VCLConfigMap VarnishVCLConfigMap `json:"vclConfigMap"`
	Deployment   VarnishDeployment   `json:"deployment"`
	Service      v1.ServiceSpec      `json:"service"`
}

type VarnishVCLConfigMap struct {
	Name         string `json:"name"`
	BackendsFile string `json:"backendsFile,omitempty"`
	DefaultFile  string `json:"defaultFile,omitempty"`

	BackendsTmplFile string
}

type VarnishDeployment struct {
	Replicas             *int32                   `json:"replicas,omitempty"`
	VarnishMemory        string                   `json:"varnishMemory,omitempty"`
	VarnishResources     *v1.ResourceRequirements `json:"varnishResources,omitempty"`
	LivenessProbe        *v1.Probe                `json:"livenessProbe,omitempty"`
	ReadinessProbe       *v1.Probe                `json:"readinessProbe,omitempty"`
	ImagePullPolicy      *v1.PullPolicy           `json:"imagePullPolicy,omitempty"`
	ImagePullSecretName  string                   `json:"imagePullSecretName"`
	VarnishRestartPolicy v1.RestartPolicy         `json:"varnishRestartPolicy,omitempty"`
	Affinity             *v1.Affinity             `json:"affinity,omitempty"`
	Tolerations          []v1.Toleration          `json:"tolerations,omitempty"`
}

// VarnishServiceStatus defines the observed state of VarnishService
type VarnishServiceStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	Deployment appsv1.DeploymentStatus     `json:"deployment,omitempty"`
	Service    VarnishServiceServiceStatus `json:"service,omitempty"`
}

// VarnishServiceSingleServiceStatus describes the status of one service as it exists within a VarnishService
type VarnishServiceSingleServiceStatus struct {
	v1.ServiceStatus `json:",inline"`
	Name             string `json:"name,omitempty"`
	IP               string `json:"ip,omitempty"`
}

// VarnishServiceServiceStatus defines the observed state of the Service portion of VarnishService
type VarnishServiceServiceStatus struct {
	Cached   VarnishServiceSingleServiceStatus `json:"cached,omitempty"`
	NoCached VarnishServiceSingleServiceStatus `json:"noCached,omitempty"`
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
