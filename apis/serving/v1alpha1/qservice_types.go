package v1alpha1

import (
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&QService{}, &QServiceList{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type QServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QService `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type QService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QServiceSpec   `json:"spec,omitempty"`
	Status QServiceStatus `json:"status,omitempty"`
}

type QServiceStatus struct {
	DeploymentStage         string                      `json:"deploymentStage,omitempty"`
	DeploymentComments      string                      `json:"deploymentComments,omitempty"`
	Ingresses               map[string][]strfmt.Ingress `json:"ingresses,omitempty"`
	appsv1.DeploymentStatus `json:",inline"`
}

type QServiceSpec struct {
	Replicas    *int32 `json:"replicas,omitempty"`
	Pod         `json:",inline"`
	Ports       []strfmt.PortForward `json:"ports,omitempty"`
	Volumes     Volumes              `json:"volumes,omitempty"`
	Resources   Resources            `json:"resources,omitempty"`
	Strategy    *strfmt.Strategy     `json:"strategy,omitempty"`
	Tolerations []strfmt.Toleration  `json:"tolerations,omitempty"`
}

type Pod struct {
	Container `json:",inline"`

	RestartPolicy                 string             `json:"restartPolicy,omitempty"`
	TerminationGracePeriodSeconds *int64             `json:"terminationGracePeriodSeconds,omitempty"`
	ActiveDeadlineSeconds         *int64             `json:"activeDeadlineSeconds,omitempty"`
	DNSConfig                     *v1.PodDNSConfig   `json:"dnsConfig,omitempty"`
	DNSPolicy                     string             `json:"dnsPolicy,omitempty"`
	NodeSelector                  map[string]string  `json:"nodeSelector,omitempty"`
	Hosts                         []strfmt.HostAlias `json:"hosts,omitempty"`
	ServiceAccountName            string             `json:"serviceAccountName,omitempty"`
}

type Container struct {
	Image           string               `json:"image,omitempty"`
	ImagePullSecret string               `json:"imagePullSecret,omitempty"`
	ImagePullPolicy string               `json:"imagePullPolicy,omitempty"`
	WorkingDir      string               `json:"workingDir,omitempty"`
	Command         []string             `json:"command,omitempty"`
	Args            []string             `json:"args,omitempty"`
	Mounts          []strfmt.VolumeMount `json:"mounts,omitempty"`
	Envs            Envs                 `json:"envs,omitempty"`
	TTY             bool                 `json:"tty,omitempty"`
	ReadinessProbe  *Probe               `json:"readinessProbe,omitempty"`
	LivenessProbe   *Probe               `json:"livenessProbe,omitempty"`
	Lifecycle       *Lifecycle           `json:"lifecycle,omitempty"`
}

type Envs map[string]string

func (envs Envs) Merge(srcEnvs Envs) Envs {
	es := Envs{}
	for k, v := range envs {
		es[k] = v
	}
	for k, v := range srcEnvs {
		es[k] = v
	}
	return es
}

type Lifecycle struct {
	PostStart *strfmt.Action `json:"postStart,omitempty" `
	PreStop   *strfmt.Action `json:"preStop,omitempty" `
}

type Probe struct {
	Action    strfmt.Action `json:"action"`
	ProbeOpts `json:",inline"`
}

type ProbeOpts struct {
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty" `
	TimeoutSeconds      int32 `json:"timeoutSeconds,omitempty" `
	PeriodSeconds       int32 `json:"periodSeconds,omitempty" `
	SuccessThreshold    int32 `json:"successThreshold,omitempty" `
	FailureThreshold    int32 `json:"failureThreshold,omitempty" `
}

type Volumes map[string]v1.VolumeSource

type Resources map[string]strfmt.RequestAndLimit
