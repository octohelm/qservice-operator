package strfmt

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-courier/x/ptr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	ErrInvalidEnvVarSource = fmt.Errorf("invalid env var source")
)

// #/field/:fieldPath?
// #/resourceField/:resource?
// #/secret(Key)?/name/key?
// #/configMap(Key)?/:name/key?
func ParseEnvVarSource(ref string) (*v1.EnvVarSource, error) {
	if strings.Index(ref, "#/") != 0 {
		return nil, ErrInvalidEnvVarSource
	}

	u, err := url.Parse(ref[1:])
	if err != nil {
		return nil, ErrInvalidEnvVarSource
	}

	parts := strings.Split(u.Path[1:], "/")

	s := &v1.EnvVarSource{}

	switch parts[0] {
	case "field":
		if len(parts) != 2 {
			return nil, ErrInvalidEnvVarSource
		}

		v := v1.ObjectFieldSelector{}
		v.FieldPath = parts[1]

		if apiVersion := u.Query().Get("apiVersion"); apiVersion != "" {
			v.APIVersion = apiVersion
		}

		s.FieldRef = &v
	case "resourceField":
		if len(parts) != 2 {
			return nil, ErrInvalidEnvVarSource
		}

		v := v1.ResourceFieldSelector{}
		v.Resource = parts[1]

		q := u.Query()

		if containerName := q.Get("containerName"); containerName != "" {
			v.ContainerName = containerName
		}

		if divisor := q.Get("divisor"); divisor != "" {
			if q, err := resource.ParseQuantity(divisor); err != nil {
				v.Divisor = q
			}
		}

		s.ResourceFieldRef = &v
	case "secret", "secretKey":
		if len(parts) != 3 {
			return nil, ErrInvalidEnvVarSource
		}

		v := v1.SecretKeySelector{}
		v.Name = parts[1]
		v.Key = parts[2]

		if optional := u.Query().Get("optional"); !(optional == "" || optional == "0") {
			v.Optional = ptr.Bool(true)
		}

		s.SecretKeyRef = &v
	case "configMap", "configMapKey":
		if len(parts) != 3 {
			return nil, ErrInvalidEnvVarSource
		}

		v := v1.ConfigMapKeySelector{}
		v.Name = parts[1]
		v.Key = parts[2]

		if optional := u.Query().Get("optional"); !(optional == "" || optional == "0") {
			v.Optional = ptr.Bool(true)
		}

		s.ConfigMapKeyRef = &v
	default:
		return nil, ErrInvalidEnvVarSource
	}

	return s, nil
}
