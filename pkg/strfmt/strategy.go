package strfmt

import (
	"bytes"
	"sort"
	"strings"
)

// RollingUpdate:maxUnavailable=25%,maxSurge=25%
func ParseStrategy(s string) (*Strategy, error) {
	if s == "" {
		return nil, nil
	}

	t := &Strategy{}

	parts := strings.Split(s, ":")

	t.Type = parts[0]
	t.Flags = map[string]string{}

	if len(parts) > 1 {
		kvs := strings.Split(parts[1], ",")

		for _, kv := range kvs {
			kvParts := strings.Split(kv, "=")

			if len(kvParts) > 1 {
				t.Flags[kvParts[0]] = kvParts[1]
			}
		}
	}

	return t, nil
}

type Strategy struct {
	Type  string
	Flags map[string]string
}

func (Strategy) OpenAPISchemaType() []string { return []string{"string"} }
func (Strategy) OpenAPISchemaFormat() string { return "strategy" }

func (t *Strategy) UnmarshalText(text []byte) error {
	to, err := ParseStrategy(string(text))
	if err != nil {
		return err
	}
	*t = *to
	return nil
}

func (t Strategy) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t Strategy) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(t.Type)

	if len(t.Flags) > 0 {
		buf.WriteRune(':')

		keys := make([]string, 0)

		for k := range t.Flags {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for i, k := range keys {
			if i != 0 {
				buf.WriteRune(',')

			}

			buf.WriteString(k)
			buf.WriteRune('=')
			buf.WriteString(t.Flags[k])
		}
	}

	return buf.String()
}
