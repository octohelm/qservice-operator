package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestStrategy(t *testing.T) {
	strategy, err := ParseStrategy("RollingUpdate:maxUnavailable=25%,maxSurge=25%")
	NewWithT(t).Expect(err).To(BeNil())

	NewWithT(t).Expect(strategy.Type).To(Equal("RollingUpdate"))
	NewWithT(t).Expect(strategy.Flags).To(Equal(map[string]string{
		"maxUnavailable": "25%",
		"maxSurge":       "25%",
	}))

	NewWithT(t).Expect(strategy.String()).To(Equal("RollingUpdate:maxSurge=25%,maxUnavailable=25%"))
}
