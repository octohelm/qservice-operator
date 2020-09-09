package converter

import (
	"github.com/octohelm/qservice-operator/pkg/strfmt"
	corev1 "k8s.io/api/core/v1"
)

func ToImagePullSecret(secret *strfmt.ImagePullSecret, namespace string) *corev1.Secret {
	s := corev1.Secret{}
	s.Namespace = namespace
	s.Name = secret.Name

	s.Type = "kubernetes.io/dockerconfigjson"
	s.Data = map[string][]byte{
		".dockerconfigjson": secret.DockerConfigJSON(),
	}
	return &s
}
