package strfmt

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMetrics(t *testing.T) {
	metrics, err := ParseMetrics(`Object(target = apps/v1.Deployment#srv-test, metricName = http_requests, targetValue = 100, selector = "app = test") Resource(name = cpu, targetAverageUtilization = 70)`)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(metrics.String()).To(Equal(`Object(metricName = "http_requests", selector = "app=test", target = "apps/v1.Deployment#srv-test", targetValue = "100") Resource(name = "cpu", targetAverageUtilization = "70")`))
}
