package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestHostAlias(t *testing.T) {
	customHost, err := ParseHostAlias("127.0.0.1 test1.com,test2.com")

	NewWithT(t).Expect(err).To(BeNil())

	NewWithT(t).Expect(customHost.IP).To(Equal("127.0.0.1"))
	NewWithT(t).Expect(customHost.HostNames).To(Equal([]string{"test1.com", "test2.com"}))

	NewWithT(t).Expect(customHost.String()).To(Equal("127.0.0.1 test1.com,test2.com"))
}
