package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestRequestAndLimit(t *testing.T) {
	t.Run("parse & string", func(t *testing.T) {
		r, err := ParseRequestAndLimit("1/500")

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(r.Request).To(Equal(1))
		NewWithT(t).Expect(r.Limit).To(Equal(500))

		NewWithT(t).Expect(r.String()).To(Equal("1/500"))
	})

	t.Run("parse & string with unit", func(t *testing.T) {
		r, err := ParseRequestAndLimit("10/500e6")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(r.Request).To(Equal(10))
		NewWithT(t).Expect(r.Limit).To(Equal(500))
		NewWithT(t).Expect(r.Unit).To(Equal("e6"))

		NewWithT(t).Expect(r.String()).To(Equal("10/500e6"))
	})

	t.Run("parse & string simple", func(t *testing.T) {
		r, err := ParseRequestAndLimit("10")
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(r.Request).To(Equal(10))
		NewWithT(t).Expect(r.Limit).To(Equal(0))

		NewWithT(t).Expect(r.String()).To(Equal("10"))
	})
}
