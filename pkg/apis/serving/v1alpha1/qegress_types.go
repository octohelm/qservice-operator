package v1alpha1

import (
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&QEgress{}, &QEgressList{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type QEgressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QEgress `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type QEgress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QEgressSpec   `json:"spec"`
	Status QEgressStatus `json:"status,omitempty"`
}

type QEgressSpec struct {
	Egress strfmt.Egress `json:"egress"`
}

type QEgressStatus struct {
}
