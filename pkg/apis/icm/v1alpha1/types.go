package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VarnishServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []VarnishService `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VarnishService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              VarnishServiceSpec   `json:"spec"`
	Status            VarnishServiceStatus `json:"status,omitempty"`
}

type VarnishServiceSpec struct {
	// Size represents the number of varnish nodes
	Replicas int32
}
type VarnishServiceStatus struct {
	// Nodes represents the names of the varnish nodes
	Nodes []string `json:"nodes"`
}
