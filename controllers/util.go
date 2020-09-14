package controllers

import (
	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	"os"
	"strings"
)

func getGateway(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts[1:], ".")
}

var IngressGateways = controllerutil.ParseIngressGateways(os.Getenv("INGRESS_GATEWAYS"))
