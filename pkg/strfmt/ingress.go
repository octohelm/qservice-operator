package strfmt

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var InvalidIngressRule = fmt.Errorf("invalid ingress rule")

// https://octohelm.tech,/v1$,/v2
func ParseIngress(ingress string) (*Ingress, error) {
	if ingress == "" {
		return nil, InvalidIngressRule
	}

	origin := ingress
	paths := ""

	if i := strings.Index(ingress, ","); i != -1 {
		origin = ingress[0:i]
		paths = ingress[i+1:]
	}

	parts := strings.Split(origin, "://")
	if len(parts) != 2 {
		return nil, InvalidIngressRule
	}

	r := &Ingress{
		Scheme: parts[0],
	}

	if r.Scheme == "" {
		r.Scheme = "http"
	}

	hostPort := strings.Split(parts[1], ":")

	if len(hostPort) == 2 {
		port, _ := strconv.ParseUint(hostPort[1], 10, 16)
		r.Port = uint16(port)
	}

	if r.Port == 0 {
		r.Port = 80
	}

	r.Host = hostPort[0]

	if len(paths) != 0 {
		parts := strings.Split(paths, ",")

		for i := range parts {
			p := parts[i]

			if len(p) == 0 {
				continue
			}

			if strings.HasSuffix(p, "$") {
				r.Paths = append(r.Paths, PathRule{
					Path:    p[0 : len(p)-1],
					Exactly: true,
				})
			} else {
				r.Paths = append(r.Paths, PathRule{
					Path: p,
				})
			}
		}
	}

	return r, nil
}

type Ingress struct {
	Scheme string
	Host   string
	Port   uint16
	Paths  []PathRule
}

type PathRule struct {
	Path    string
	Exactly bool
}

func (Ingress) OpenAPISchemaType() []string { return []string{"string"} }
func (Ingress) OpenAPISchemaFormat() string { return "ingress" }

func (ingress Ingress) String() string {
	ret := bytes.NewBuffer(nil)

	defaultPort := 0

	if ingress.Scheme == "" {
		ret.WriteString("http://")
		defaultPort = 80
	} else {
		ret.WriteString(ingress.Scheme)
		ret.WriteString("://")
	}

	ret.WriteString(ingress.Host)

	if ingress.Port == 0 {
		_, _ = fmt.Fprintf(ret, ":%d", defaultPort)
	} else {
		_, _ = fmt.Fprintf(ret, ":%d", ingress.Port)
	}

	for _, path := range ingress.Paths {
		ret.WriteString(",")
		ret.WriteString(path.Path)

		if path.Exactly {
			ret.WriteString("$")
		}
	}

	return ret.String()
}

func (ingress Ingress) MarshalText() ([]byte, error) {
	return []byte(ingress.String()), nil
}

func (ingress *Ingress) UnmarshalText(data []byte) error {
	ir, err := ParseIngress(string(data))
	if err != nil {
		return err
	}
	*ingress = *ir
	return nil
}
