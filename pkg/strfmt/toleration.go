package strfmt

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// key=value:NoExecute,3600
// key:NoExecute
func ParseToleration(s string) (*Toleration, error) {
	if s == "" {
		return nil, nil
	}

	t := &Toleration{}

	parts := strings.Split(s, ":")

	kv := strings.Split(parts[0], "=")

	t.Key = kv[0]

	if len(kv) > 1 {
		t.Value = kv[1]
	}

	if len(parts) > 1 {
		effectAndDuration := strings.Split(parts[1], ",")
		t.Effect = effectAndDuration[0]

		if len(effectAndDuration) > 1 {
			d, err := strconv.ParseInt(effectAndDuration[1], 10, 64)
			if err != nil {
				return nil, errors.New("invalid toleration seconds")
			}
			t.TolerationSeconds = &d
		}
	}

	return t, nil
}

type Toleration struct {
	Key               string
	Value             string
	Effect            string
	TolerationSeconds *int64
}

func (Toleration) OpenAPISchemaType() []string { return []string{"string"} }
func (Toleration) OpenAPISchemaFormat() string { return "toleration" }

func (t *Toleration) UnmarshalText(text []byte) error {
	to, err := ParseToleration(string(text))
	if err != nil {
		return err
	}
	*t = *to
	return nil
}

func (t Toleration) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t Toleration) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(t.Key)

	if t.Value != "" {
		buf.WriteRune('=')
		buf.WriteString(t.Value)
	}

	if t.Effect != "" {
		buf.WriteRune(':')
		buf.WriteString(t.Effect)
	}

	if t.TolerationSeconds != nil {
		buf.WriteRune(',')
		buf.WriteString(strconv.FormatInt(int64(*t.TolerationSeconds), 10))
	}

	return buf.String()
}
