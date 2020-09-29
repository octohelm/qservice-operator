package controllers

import (
	"encoding/json"
	"testing"

	"github.com/octohelm/qservice-operator/apis/serving/v1alpha1"

	"github.com/octohelm/qservice-operator/pkg/strfmt"
)

func TestToExternalServiceEntity(t *testing.T) {
	t.Run("ip", func(t *testing.T) {
		u, _ := strfmt.ParseEgress("redis://127.0.0.1:6379")
		s := toExternalServiceEntity(&v1alpha1.QEgress{
			Spec: v1alpha1.QEgressSpec{
				Egress: *u,
			},
		})
		data, _ := json.MarshalIndent(s, "", "  ")
		t.Log(string(data))
	})

	t.Run("hostname", func(t *testing.T) {
		u, _ := strfmt.ParseEgress("https://sms.tencentcloudapi.com")
		s := toExternalServiceEntity(&v1alpha1.QEgress{
			Spec: v1alpha1.QEgressSpec{
				Egress: *u,
			},
		})
		data, _ := json.MarshalIndent(s, "", "  ")
		t.Log(string(data))
	})
}
