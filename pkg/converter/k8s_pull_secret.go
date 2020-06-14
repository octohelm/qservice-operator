package converter

import (
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	corev1 "k8s.io/api/core/v1"
)

func ToImagePullSecret(secret *strfmt.ImagePullSecret) *corev1.Secret {
	s := corev1.Secret{}
	s.Type = "kubernetes.io/dockerconfigjson"
	s.Name = secret.Name
	s.Data = map[string][]byte{
		".dockerconfigjson": secret.DockerConfigJSON(),
	}
	return &s
}
