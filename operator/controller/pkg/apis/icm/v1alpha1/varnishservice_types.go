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

// SharedVolumeSpec is just a wrapper for settings relating to a shared volume
type SharedVolumeSpec struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type VarnishDeployment struct {
	Replicas                 int32                    `json:"replicas"`
	VarnishMemory            string                   `json:"varnishMemory,omitempty"`
	VarnishResources         *v1.ResourceRequirements `json:"varnishResources,omitempty"`
	VarnishExporterResources *v1.ResourceRequirements `json:"varnishExporterResources,omitempty"`
	LivenessProbe            *v1.Probe                `json:"livenessProbe,omitempty"`
	ReadinessProbe           *v1.Probe                `json:"readinessProbe,omitempty"`
	ImagePullSecretName      string                   `json:"imagePullSecretName"`
	VarnishRestartPolicy     *v1.RestartPolicy        `json:"varnishRestartPolicy,omitempty"`
	SharedVolume             SharedVolumeSpec         `json:"sharedVolume,omitempty"`
	BackendsFile             string                   `json:"backendsFile,omitempty"`
	DefaultFile              string                   `json:"defaultFile,omitempty"`
	Affinity                 *v1.Affinity             `json:"affinity,omitempty"`
	Tolerations              []v1.Toleration          `json:"tolerations,omitempty"`
}

// VarnishServiceSpec defines the desired state of VarnishService
// Important: Run "make" to regenerate code after modifying this file
type VarnishServiceSpec struct {
	Service    v1.ServiceSpec `json:"service"`
	Deployment VarnishDeployment
}

// VarnishServiceStatus defines the observed state of VarnishService
type VarnishServiceStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	Deployment appsv1.DeploymentStatus `json:"deployment"`
	// Replicas represents the current number of varnish nodes
	CacheBypassIP string `json:"cacheBypassIP"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishService is the Schema for the varnishservices API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type VarnishService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VarnishServiceSpec   `json:"spec,omitempty"`
	Status VarnishServiceStatus `json:"status,omitempty"`
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