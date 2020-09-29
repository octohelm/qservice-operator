package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

func TestIngress(t *testing.T) {
	t.Run("parse & string", func(t *testing.T) {
		ingress, err := ParseIngress("http://*.helmx,/v1$,/v2")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(ingress.Port).To(Equal(uint16(80)))
		NewWithT(t).Expect(ingress.Host).To(Equal("*.helmx"))
		NewWithT(t).Expect(ingress.Scheme).To(Equal("http"))

		NewWithT(t).Expect(ingress.String()).To(Equal("http://*.helmx:80,/v1$,/v2"))
	})

	t.Run("parse & string without paths", func(t *testing.T) {
		ingress, err := ParseIngress("http://*.helmx")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(ingress.Port).To(Equal(uint16(80)))
		NewWithT(t).Expect(ingress.Host).To(Equal("*.helmx"))
		NewWithT(t).Expect(ingress.Scheme).To(Equal("http"))

		NewWithT(t).Expect(ingress.String()).To(Equal("http://*.helmx:80"))
	})

	t.Run("yaml marshal & unmarshal", func(t *testing.T) {
		data, err := yaml.Marshal(struct {
			Ingress Ingress `yaml:"ingress"`
		}{
			Ingress: Ingress{
				Port: 80,
				Host: "*.helmx",
				Paths: []PathRule{{
					Path: "/helmx",
				}},
			},
		})
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(data)).To(Equal("ingress: http://*.helmx:80,/helmx\n"))

		v := struct {
			Ingress Ingress `yaml:"ingress"`
		}{}

		err = yaml.Unmarshal(data, &v)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(v.Ingress.String()).To(Equal("http://*.helmx:80,/helmx"))
	})
}
