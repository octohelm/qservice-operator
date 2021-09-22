package strfmt

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestAction(t *testing.T) {
	t.Run("parse & string http", func(t *testing.T) {
		action, err := ParseAction("http://:80/healthy")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

		gomega.NewWithT(t).Expect(string(action.HTTPGet.Scheme)).To(gomega.Equal("HTTP"))
		gomega.NewWithT(t).Expect(action.HTTPGet.Port.String()).To(gomega.Equal("80"))
		gomega.NewWithT(t).Expect(action.HTTPGet.Host).To(gomega.Equal(""))
		gomega.NewWithT(t).Expect(action.HTTPGet.Path).To(gomega.Equal("/healthy"))

		gomega.NewWithT(t).Expect(action.String()).To(gomega.Equal("http://:80/healthy"))
	})

	t.Run("parse & string tcp", func(t *testing.T) {
		action, err := ParseAction("tcp://:22")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

		gomega.NewWithT(t).Expect(action.TCPSocket.Port.String()).To(gomega.Equal("22"))
		gomega.NewWithT(t).Expect(action.TCPSocket.Host).To(gomega.Equal(""))

		gomega.NewWithT(t).Expect(action.String()).To(gomega.Equal("tcp://:22"))
	})
}
