package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishService describes what a varnish service looks like
type VarnishService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VarnishServiceSpec   `json:"spec"`
	Status VarnishServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VarnishServiceList is a list of VarnishService resources
type VarnishServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []VarnishService `json:"items"`
}

// VarnishServiceSpec describes the spec for a VarnishService
type VarnishServiceSpec struct {
	// Replicas represents the number of varnish nodes
	Replicas int32
}

// VarnishServiceStatus describes the status for a VarnishService
type VarnishServiceStatus struct {
	// Nodes represents the names of the varnish nodes
	Nodes []string `json:"nodes"`
}
