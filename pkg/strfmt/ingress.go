package strfmt

import (
	"fmt"
	"net/url"
	"strconv"
)

func ParseIngress(s string) (*Ingress, error) {
	if s == "" {
		return nil, fmt.Errorf("invalid ingress rule")
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	r := &Ingress{
		Scheme: u.Scheme,
		Host:   u.Hostname(),
		Path:   u.Path,
	}

	if r.Scheme == "" {
		r.Scheme = "http"
	}

	p := u.Port()
	if p == "" {
		r.Port = 80
	} else {
		port, _ := strconv.ParseUint(p, 10, 16)
		r.Port = uint16(port)
	}

	return r, nil
}

// openapi:strfmt ingress-rule
type Ingress struct {
	Scheme string
	Host   string
	Path   string
	Port   uint16
}

func (r Ingress) String() string {
	if r.Scheme == "" {
		r.Scheme = "http"
	}
	if r.Port == 0 {
		r.Port = 80
	}

	return (&url.URL{
		Scheme: r.Scheme,
		Host:   r.Host + ":" + strconv.FormatUint(uint64(r.Port), 10),
		Path:   r.Path,
	}).String()
}

func (r Ingress) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Ingress) UnmarshalText(data []byte) error {
	ir, err := ParseIngress(string(data))
	if err != nil {
		return err
	}
	*r = *ir
	return nil
}
