package strfmt

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestParseEnvVarSource(t *testing.T) {
	t.Run("parse fieldRef", func(t *testing.T) {
		envvarResource, err := ParseEnvVarSource("#/field/spec.nodeName")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(envvarResource.FieldRef.FieldPath).To(gomega.Equal("spec.nodeName"))
	})

	t.Run("parse resourceFieldRef", func(t *testing.T) {
		envvarResource, err := ParseEnvVarSource("#/resourceField/limits.cpu")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(envvarResource.ResourceFieldRef.Resource).To(gomega.Equal("limits.cpu"))
	})

	t.Run("parse secretKeyRef", func(t *testing.T) {
		envvarResource, err := ParseEnvVarSource("#/secret/pg.octohelm/username")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(envvarResource.SecretKeyRef.Name).To(gomega.Equal("pg.octohelm"))
		gomega.NewWithT(t).Expect(envvarResource.SecretKeyRef.Key).To(gomega.Equal("username"))
	})

	t.Run("parse configMapKeyRef", func(t *testing.T) {
		envvarResource, err := ParseEnvVarSource("#/configMap/pg.octohelm/username")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(envvarResource.ConfigMapKeyRef.Name).To(gomega.Equal("pg.octohelm"))
		gomega.NewWithT(t).Expect(envvarResource.ConfigMapKeyRef.Key).To(gomega.Equal("username"))
	})
}
