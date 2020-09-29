package strfmt

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

var InvalidEgressRule = fmt.Errorf("invalid egress rule")

func ParseEgress(egress string) (*Egress, error) {
	if egress == "" {
		return nil, InvalidEgressRule
	}

	u, err := url.Parse(egress)
	if err != nil {
		return nil, err
	}

	hostname := u.Hostname()

	e := &Egress{
		Scheme: strings.ToLower(u.Scheme),
	}

	ip := net.ParseIP(u.Hostname())

	if ip != nil {
		e.IP = ip
	} else {
		e.Hostname = hostname
	}

	port, _ := strconv.ParseUint(u.Port(), 10, 16)
	e.Port = uint16(port)

	return e, nil
}

type Egress struct {
	Scheme   string
	Hostname string
	IP       net.IP
	Port     uint16
}

func (Egress) OpenAPISchemaType() []string { return []string{"string"} }
func (Egress) OpenAPISchemaFormat() string { return "egress" }

func (egress Egress) String() string {
	ret := bytes.NewBuffer(nil)

	ret.WriteString(egress.Scheme)
	ret.WriteString("://")

	if egress.IP != nil {
		ret.WriteString(egress.IP.String())
	} else {
		ret.WriteString(egress.Hostname)
	}

	if egress.Port > 0 {
		_, _ = fmt.Fprintf(ret, ":%d", egress.Port)
	}

	return ret.String()
}

func (egress Egress) MarshalText() ([]byte, error) {
	return []byte(egress.String()), nil
}

func (egress *Egress) UnmarshalText(data []byte) error {
	eg, err := ParseEgress(string(data))
	if err != nil {
		return err
	}
	*egress = *eg
	return nil
}
