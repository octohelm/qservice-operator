package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

func TestPortForward(t *testing.T) {
	t.Run("parse & string", func(t *testing.T) {
		portForward, err := ParsePortForward("http-80:8080/TCP")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(portForward.AppProtocol).To(Equal("http"))
		NewWithT(t).Expect(portForward.Port).To(Equal(uint16(80)))
		NewWithT(t).Expect(portForward.ContainerPort).To(Equal(uint16(8080)))
		NewWithT(t).Expect(portForward.Protocol).To(Equal("TCP"))

		NewWithT(t).Expect(portForward.String()).To(Equal("http-80:8080/TCP"))
	})

	t.Run("parse & string without target port ", func(t *testing.T) {
		portForward, err := ParsePortForward("80/TCP")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(portForward.AppProtocol).To(Equal(""))
		NewWithT(t).Expect(portForward.Port).To(Equal(uint16(80)))
		NewWithT(t).Expect(portForward.ContainerPort).To(Equal(uint16(80)))
		NewWithT(t).Expect(portForward.Protocol).To(Equal("TCP"))

		NewWithT(t).Expect(portForward.String()).To(Equal("80/TCP"))
	})

	t.Run("parse & string without node port", func(t *testing.T) {
		portForward, err := ParsePortForward("!20000:8080")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(portForward.IsNodePort).To(BeTrue())
		NewWithT(t).Expect(portForward.Port).To(Equal(uint16(20000)))
		NewWithT(t).Expect(portForward.ContainerPort).To(Equal(uint16(8080)))

		NewWithT(t).Expect(portForward.String()).To(Equal("!20000:8080"))
	})

	t.Run("parse & string without protocol", func(t *testing.T) {
		portForward, err := ParsePortForward("80:8080")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(portForward.Port).To(Equal(uint16(80)))
		NewWithT(t).Expect(portForward.ContainerPort).To(Equal(uint16(8080)))

		NewWithT(t).Expect(portForward.String()).To(Equal("80:8080"))
	})

	t.Run("yaml marshal & unmarshal", func(t *testing.T) {
		data, err := yaml.Marshal(struct {
			Port PortForward `yaml:"port"`
		}{
			Port: PortForward{
				Port:          80,
				ContainerPort: 8080,
				Protocol:      "TCP",
			},
		})
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(data)).To(Equal("port: 80:8080/TCP\n"))

		v := struct {
			Port PortForward `yaml:"port"`
		}{}

		err = yaml.Unmarshal(data, &v)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(v.Port.String()).To(Equal("80:8080/TCP"))
	})

	t.Run("node port range in 20000-40000", func(t *testing.T) {
		_, ltErr := ParsePortForward("!19999:80")
		NewWithT(t).Expect(ltErr).NotTo(BeNil())

		_, noLtErr := ParsePortForward("!20000:80")
		NewWithT(t).Expect(noLtErr).To(BeNil())

		_, gtErr := ParsePortForward("!40001:80")
		NewWithT(t).Expect(gtErr).NotTo(BeNil())

		_, noGtErr := ParsePortForward("!40000:80")
		NewWithT(t).Expect(noGtErr).To(BeNil())
	})
}
