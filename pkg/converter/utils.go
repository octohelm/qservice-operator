package converter

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

func toProtocol(protocol string) v1.Protocol {
	switch strings.ToUpper(protocol) {
	case "SCTP":
		return v1.ProtocolSCTP
	case "UDP":
		return v1.ProtocolUDP
	case "TCP":
		return v1.ProtocolTCP
	default:
		return v1.ProtocolTCP
	}
}

func cloneKV(from map[string]string) map[string]string {
	m := map[string]string{}
	for k, v := range from {
		m[k] = v

	}
	return m
}
