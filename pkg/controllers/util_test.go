package controllers

import (
	"testing"

	"github.com/onsi/gomega"
)

func Test_toDNS1121Safe(t *testing.T) {
	gomega.NewWithT(t).Expect(safeDNS1121Host("prometheus-kube-prometheus-kube-controller-manager---kube-system.auto-internal")).To(gomega.Equal("prometheus-kube-prometheus-kube-controller-manager---k-f6e94d44.auto-internal"))
	gomega.NewWithT(t).Expect(safeDNS1121Host("prometheus-kube---kube-system.auto-internal")).To(gomega.Equal("prometheus-kube---kube-system.auto-internal"))
}
