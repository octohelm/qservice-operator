package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestToleration(t *testing.T) {
	toleration, err := ParseToleration("key=value:NoExecute,3600")
	NewWithT(t).Expect(err).To(BeNil())

	NewWithT(t).Expect(toleration.Key).To(Equal("key"))
	NewWithT(t).Expect(toleration.Value).To(Equal("value"))
	NewWithT(t).Expect(toleration.Effect).To(Equal("NoExecute"))
	NewWithT(t).Expect(*toleration.TolerationSeconds).To(Equal(int64(3600)))

	NewWithT(t).Expect(toleration.String()).To(Equal("key=value:NoExecute,3600"))
}
