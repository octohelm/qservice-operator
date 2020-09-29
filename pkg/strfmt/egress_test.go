package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestEgress(t *testing.T) {
	t.Run("parse & string host", func(t *testing.T) {
		egress, err := ParseEgress("https://google.com")
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(egress.String()).To(Equal("https://google.com"))
	})

	t.Run("parse & string ip", func(t *testing.T) {
		egress, err := ParseEgress("tcp://127.0.0.1")
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(egress.String()).To(Equal("tcp://127.0.0.1"))
	})
}
