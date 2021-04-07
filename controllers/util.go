package controllers

import (
	"os"
	"strings"

	"github.com/octohelm/qservice-operator/pkg/controllerutil"
	"github.com/octohelm/qservice-operator/pkg/converter"
)

func getGateway(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts[1:], ".")
}

var IngressGateways = controllerutil.ParseIngressGateways(os.Getenv("INGRESS_GATEWAYS"))

var DNS1123LabelMaxLength = 63

func safeDNS1121Host(hostname string) string {
	parts := strings.SplitN(hostname, ".", 2)
	if len(parts) == 2 && len(parts[0]) > DNS1123LabelMaxLength {
		n := converter.HashID(hostname)
		return parts[0][0:DNS1123LabelMaxLength-len("-"+n)] + "-" + n + "." + parts[1]
	}
	return hostname
}
