package converter

import (
	"strings"
)

func ConvertToGatewayName(host string) string {
	return "istio-ingress/" + strings.Join(strings.Split(host, ".")[1:], "--")
}

func GatewaySafeName(gateway string) string {
	return strings.Join(strings.Split(gateway, "."), "--")
}
