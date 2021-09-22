package strfmt

import (
	"net/url"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// http://:80
// tcp://:80
// exec
func ParseAction(s string) (*Action, error) {
	if s == "" {
		return nil, nil
	}

	a := &Action{}

	if strings.HasPrefix(s, "http") || strings.HasPrefix(s, "tcp") {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}

		port, _ := strconv.ParseUint(u.Port(), 10, 64)

		if u.Scheme == "tcp" {
			a.TCPSocket = &v1.TCPSocketAction{}
			a.TCPSocket.Host = u.Hostname()
			a.TCPSocket.Port = intstr.FromInt(int(port))
			return a, nil
		}

		a.HTTPGet = &v1.HTTPGetAction{}
		a.HTTPGet.Port = intstr.FromInt(int(port))
		a.HTTPGet.Host = u.Hostname()
		a.HTTPGet.Path = u.Path
		a.HTTPGet.Scheme = v1.URIScheme(strings.ToUpper(u.Scheme))

		return a, nil
	}

	a.Exec = &v1.ExecAction{
		Command: []string{"sh", "-c", s},
	}

	return a, nil
}

type Action struct {
	v1.ProbeHandler
}

func (Action) OpenAPISchemaType() []string { return []string{"string"} }
func (Action) OpenAPISchemaFormat() string { return "action" }

func (a Action) String() string {
	if a.Exec != nil {
		return a.Exec.Command[2]
	}

	if a.HTTPGet != nil {
		u := &url.URL{}
		u.Scheme = strings.ToLower(string(a.HTTPGet.Scheme))
		u.Path = a.HTTPGet.Path
		u.Host = a.HTTPGet.Host + ":" + a.HTTPGet.Port.String()

		if u.Scheme != "" {
			u.Scheme = "http"
		}
		return u.String()
	}

	if a.TCPSocket != nil {
		u := &url.URL{}
		u.Scheme = "tcp"
		u.Host = a.TCPSocket.Host + ":" + a.TCPSocket.Port.String()

		return u.String()
	}

	return ""
}

func (a Action) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Action) UnmarshalText(data []byte) error {
	action, err := ParseAction(string(data))
	if err != nil {
		return err
	}
	if action != nil {
		*a = *action
	}
	return nil
}
