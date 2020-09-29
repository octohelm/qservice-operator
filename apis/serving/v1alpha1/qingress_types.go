package v1alpha1

import (
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&QIngress{}, &QIngressList{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type QIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QIngress `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type QIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QIngressSpec   `json:"spec"`
	Status QIngressStatus `json:"status,omitempty"`
}

type QIngressSpec struct {
	Ingress strfmt.Ingress                   `json:"ingress"`
	Backend extensionsv1beta1.IngressBackend `json:"backend"`
}

type QIngressStatus struct {
}
