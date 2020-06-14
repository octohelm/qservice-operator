package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestImagePullSecret(t *testing.T) {
	ips, err := ParseImagePullSecret("docker-hub://pull-only:123123@docker.io/rk-")
	NewWithT(t).Expect(err).To(BeNil())

	NewWithT(t).Expect(ips.Name).To(Equal("docker-hub"))
	NewWithT(t).Expect(ips.Host).To(Equal("docker.io"))
	NewWithT(t).Expect(ips.Username).To(Equal("pull-only"))
	NewWithT(t).Expect(ips.Password).To(Equal("123123"))
	NewWithT(t).Expect(ips.Prefix).To(Equal("/rk-"))

	NewWithT(t).Expect(ips.String()).To(Equal("docker-hub://pull-only:123123@docker.io/rk-"))

	NewWithT(t).Expect(ips.PrefixTag("nginx:alpine")).To(Equal("docker.io/rk-library/nginx:alpine"))
}
