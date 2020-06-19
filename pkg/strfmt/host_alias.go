package strfmt

import (
	"bytes"
	"strings"
)

type HostAlias struct {
	IP        string
	HostNames []string
}

func (HostAlias) OpenAPISchemaType() []string { return []string{"string"} }
func (HostAlias) OpenAPISchemaFormat() string { return "host-alias" }

// 127.0.0.1 test1.com,test2.com
func ParseHostAlias(s string) (*HostAlias, error) {
	if s == "" {
		return nil, nil
	}

	t := &HostAlias{}

	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		parts = strings.Split(s, " ")
		if len(parts) < 2 {
			return nil, nil
		}
	}

	t.IP = parts[0]

	kv := strings.Split(parts[1], ",")

	if len(kv) > 0 {
		t.HostNames = append(t.HostNames, kv...)
	}

	return t, nil
}

func (t *HostAlias) UnmarshalText(text []byte) error {
	to, err := ParseHostAlias(string(text))
	if err != nil {
		return err
	}
	*t = *to
	return nil
}

func (t HostAlias) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t HostAlias) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(t.IP)
	buf.WriteString(" ")

	if len(t.HostNames) != 0 {
		for index, host := range t.HostNames {
			buf.WriteString(host)
			if index != len(t.HostNames)-1 {
				buf.WriteRune(',')
			}
		}
	}

	return buf.String()
}
