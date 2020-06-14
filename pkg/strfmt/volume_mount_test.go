package strfmt

import (
	"testing"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/gomega"
)

func TestMount(t *testing.T) {
	t.Run("parse & string simple", func(t *testing.T) {
		r, err := ParseVolumeMount("data:/html")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(r.Name).To(Equal("data"))
		NewWithT(t).Expect(r.MountPath).To(Equal("/html"))
		NewWithT(t).Expect(r.ReadOnly).To(Equal(false))
		NewWithT(t).Expect(r.SubPath).To(Equal(""))

		NewWithT(t).Expect(r.String()).To(Equal("data:/html"))
	})

	t.Run("parse & string", func(t *testing.T) {
		r, err := ParseVolumeMount("data/html:/html:ro")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(r.Name).To(Equal("data"))
		NewWithT(t).Expect(r.MountPath).To(Equal("/html"))
		NewWithT(t).Expect(r.ReadOnly).To(Equal(true))
		NewWithT(t).Expect(r.SubPath).To(Equal("html"))

		NewWithT(t).Expect(r.String()).To(Equal("data/html:/html:ro"))
	})

	t.Run("VolumeMount yaml marshal & unmarshal", func(t *testing.T) {
		data, err := yaml.Marshal(struct {
			Mount VolumeMount `yaml:"volumeMount"`
		}{
			Mount: VolumeMount{
				MountPath: "/html",
				Name:      "data",
				ReadOnly:  true,
				SubPath:   "html",
			},
		})

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(data)).To(Equal("volumeMount: data/html:/html:ro\n"))

		v := struct {
			Mount VolumeMount `yaml:"volumeMount"`
		}{}

		err = yaml.Unmarshal(data, &v)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(v.Mount.String()).To(Equal("data/html:/html:ro"))
	})
}
