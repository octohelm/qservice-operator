package converter

import (
	"encoding/hex"
	"net"
	"net/url"
	"strconv"
	"strings"

	"istio.io/api/networking/v1alpha3"
	istiov1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func ToExternalServiceEntity(u *url.URL) *istiov1alpha3.ServiceEntry {
	hostname := u.Hostname()

	if len(strings.Split(hostname, ".")) <= 2 {
		// in cluster
		return nil
	}

	se := &istiov1alpha3.ServiceEntry{}
	portNumber, _ := strconv.ParseUint(u.Port(), 10, 64)

	protocol := ""
	prefix := "ext-"

	scheme := strings.ToLower(u.Scheme)

	switch scheme {
	case "http":
		protocol = "HTTP"
		if portNumber == 0 {
			portNumber = 80
		}
	case "https":
		protocol = "HTTPS"
		if portNumber == 0 {
			portNumber = 443
		}
	case "grpc":
		protocol = "GRPC"
	default:
		protocol = "TCP"
		prefix = "ext-" + scheme + "-"
	}

	ort := v1alpha3.Port{
		Name:     strings.ToLower(protocol) + "-" + strconv.FormatUint(portNumber, 10),
		Number:   uint32(portNumber),
		Protocol: protocol,
	}

	se.Name = prefix + ort.Name + "--" + hostname

	se.Spec.Location = v1alpha3.ServiceEntry_MESH_EXTERNAL
	se.Spec.ExportTo = []string{"."}

	ip := net.ParseIP(hostname)

	if ip == nil {
		// host
		se.Spec.Hosts = []string{hostname}
		se.Spec.Resolution = v1alpha3.ServiceEntry_DNS
	} else {
		se.Spec.Hosts = []string{
			prefix + ort.Name + "-" + hexIP(ip),
		}

		ipv4 := ip.To4().String()

		se.Spec.Addresses = []string{ipv4}
		se.Spec.Endpoints = []*v1alpha3.WorkloadEntry{{Address: ipv4}}
		se.Spec.Resolution = v1alpha3.ServiceEntry_STATIC
	}

	se.Spec.Ports = []*v1alpha3.Port{&ort}

	return se
}

func hexIP(ip net.IP) string {
	v4 := ip.To4()
	return hex.EncodeToString([]byte{v4[0], v4[1], v4[2], v4[3]})
}
