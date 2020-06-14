package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

func TestIngress(t *testing.T) {
	t.Run("parse & string", func(t *testing.T) {
		ingress, err := ParseIngress("http://*.helmx/helmx")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(ingress.Port).To(Equal(uint16(80)))
		NewWithT(t).Expect(ingress.Host).To(Equal("*.helmx"))
		NewWithT(t).Expect(ingress.Scheme).To(Equal("http"))

		NewWithT(t).Expect(ingress.String()).To(Equal("http://*.helmx:80/helmx"))
	})

	t.Run("yaml marshal & unmarshal", func(t *testing.T) {
		data, err := yaml.Marshal(struct {
			Ingress Ingress `yaml:"ingress"`
		}{
			Ingress: Ingress{
				Port: 80,
				Host: "*.helmx",
				Path: "/helmx",
			},
		})
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(data)).To(Equal("ingress: http://*.helmx:80/helmx\n"))

		v := struct {
			Ingress Ingress `yaml:"ingress"`
		}{}

		err = yaml.Unmarshal(data, &v)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(v.Ingress.String()).To(Equal("http://*.helmx:80/helmx"))
	})
}
