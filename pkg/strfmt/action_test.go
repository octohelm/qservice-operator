package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestAction(t *testing.T) {
	t.Run("parse & string http", func(t *testing.T) {
		action, err := ParseAction("http://:80/healthy")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(string(action.HTTPGet.Scheme)).To(Equal("HTTP"))
		NewWithT(t).Expect(action.HTTPGet.Port.String()).To(Equal("80"))
		NewWithT(t).Expect(action.HTTPGet.Host).To(Equal(""))
		NewWithT(t).Expect(action.HTTPGet.Path).To(Equal("/healthy"))

		NewWithT(t).Expect(action.String()).To(Equal("http://:80/healthy"))
	})

	t.Run("parse & string tcp", func(t *testing.T) {
		action, err := ParseAction("tcp://:22")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(action.TCPSocket.Port.String()).To(Equal("22"))
		NewWithT(t).Expect(action.TCPSocket.Host).To(Equal(""))

		NewWithT(t).Expect(action.String()).To(Equal("tcp://:22"))
	})
}
