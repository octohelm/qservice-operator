package controllerutil

import (
	"github.com/onsi/gomega"
	"testing"
)

func TestParseGateways(t *testing.T) {
	gateways := ParseIngressGateways(
		"internal:xxx.internal.com,auto-internal:xxx.auto.internal,external:xxx.external,external-2:xxx.external-2",
	)

	h, _ := gateways.IngressGatewayHost("service.internal")
	gomega.NewWithT(t).Expect(h).To(gomega.Equal("service.xxx.internal.com"))
}
