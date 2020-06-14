package strfmt

import (
	"fmt"
	"regexp"
	"strconv"
)

var reRequestAndLimit = regexp.MustCompile(`^([+\-]?[0-9.]+)?(\/([+\-]?[0-9.]+))?([eEinumkKMGTP]*[-+]?[0-9]*)$`)

func ParseRequestAndLimit(s string) (*RequestAndLimit, error) {
	if s == "" || !reRequestAndLimit.MatchString(s) {
		return nil, fmt.Errorf("missing request and limit")
	}

	parts := reRequestAndLimit.FindAllStringSubmatch(s, 1)[0]

	rl := &RequestAndLimit{}

	if parts[1] != "" {
		i, err := strconv.ParseInt(parts[1], 10, 64)
		if err == nil {
			rl.Request = int(i)
		}
	}

	if parts[3] != "" {
		i, err := strconv.ParseInt(parts[3], 10, 64)
		if err == nil {
			rl.Limit = int(i)
		}
	}

	if parts[4] != "" {
		rl.Unit = parts[4]
	}

	return rl, nil
}

// openapi:strfmt request-and-limit
type RequestAndLimit struct {
	Request int
	Limit   int
	Unit    string
}

func (s RequestAndLimit) String() string {
	v := ""
	if s.Request != 0 {
		v = strconv.FormatInt(int64(s.Request), 10)
	}
	if s.Limit != 0 {
		v += "/" + strconv.FormatInt(int64(s.Limit), 10)
	}
	if s.Unit != "" {
		v += s.Unit
	}
	return v
}

func (s RequestAndLimit) RequestString() string {
	v := ""
	if s.Request != 0 {
		v = strconv.FormatInt(int64(s.Request), 10)
	}
	if s.Unit != "" {
		v += s.Unit
	}
	return v
}

func (s RequestAndLimit) LimitString() string {
	v := ""
	if s.Limit != 0 {
		v = strconv.FormatInt(int64(s.Limit), 10)
	}
	if s.Unit != "" {
		v += s.Unit
	}
	return v
}

func (s RequestAndLimit) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *RequestAndLimit) UnmarshalText(data []byte) error {
	servicePort, err := ParseRequestAndLimit(string(data))
	if err != nil {
		return err
	}
	*s = *servicePort
	return nil
}
