package strfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePortForward(s string) (*PortForward, error) {
	if s == "" {
		return nil, fmt.Errorf("missing port value")
	}

	appProtocol := ""
	port := uint16(0)
	targetPort := uint16(0)
	protocol := ""
	isNodePort := false

	parts := strings.Split(s, "/")

	s = parts[0]

	if len(parts) == 2 {
		protocol = strings.ToLower(parts[1])
	}

	if s[0] == '!' {
		isNodePort = true
		s = s[1:]
	}

	ports := strings.Split(s, ":")

	portStr := ports[0]

	appProtocolAndPort := strings.Split(portStr, "-")

	if len(appProtocolAndPort) == 2 {
		portStr = appProtocolAndPort[1]
		appProtocol = strings.ToLower(appProtocolAndPort[0])
	}

	p, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid port %v", ports[0])
	}

	port = uint16(p)

	if len(ports) == 2 {
		if isNodePort {
			if port < 20000 || p > 40000 {
				return nil, fmt.Errorf("invalid value: %d: provided port is not in the valid range. The range of valid ports is 20000-40000", port)
			}
		}
		p, err := strconv.ParseUint(ports[1], 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid target port %v", ports[1])
		}
		targetPort = uint16(p)
	} else {
		targetPort = port
	}

	return &PortForward{
		AppProtocol:   appProtocol,
		Port:          port,
		IsNodePort:    isNodePort,
		ContainerPort: targetPort,
		Protocol:      strings.ToUpper(protocol),
	}, nil
}

type PortForward struct {
	AppProtocol   string
	Port          uint16
	IsNodePort    bool
	ContainerPort uint16
	Protocol      string
}

func (PortForward) OpenAPISchemaType() []string { return []string{"string"} }
func (PortForward) OpenAPISchemaFormat() string { return "port-forward" }

func (s PortForward) String() string {
	v := ""
	if s.IsNodePort {
		v = "!"
	}

	if s.AppProtocol != "" {
		v += s.AppProtocol + "-"
	}

	if s.Port != 0 {
		v += strconv.FormatUint(uint64(s.Port), 10)
	}

	if s.ContainerPort != 0 && s.ContainerPort != s.Port {
		v += ":" + strconv.FormatUint(uint64(s.ContainerPort), 10)
	}

	if s.Protocol != "" {
		v += "/" + strings.ToUpper(s.Protocol)
	}

	return v
}

func (s PortForward) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *PortForward) UnmarshalText(data []byte) error {
	servicePort, err := ParsePortForward(string(data))
	if err != nil {
		return err
	}
	*s = *servicePort
	return nil
}
