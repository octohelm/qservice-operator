package strfmt

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
)

func ParseImagePullSecret(uri string) (*ImagePullSecret, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	endpoint := &ImagePullSecret{}

	endpoint.Name = u.Scheme
	endpoint.Host = u.Host
	endpoint.Prefix = u.Path

	if u.User != nil {
		endpoint.Username = u.User.Username()
		endpoint.Password, _ = u.User.Password()
	}

	return endpoint, nil
}

type ImagePullSecret struct {
	Name     string
	Host     string
	Username string
	Password string
	Prefix   string
}

func (ImagePullSecret) OpenAPISchemaType() []string { return []string{"string"} }
func (ImagePullSecret) OpenAPISchemaFormat() string { return "image-pull-secret" }

func (s *ImagePullSecret) SecretName() string {
	return s.Name
}

func (s *ImagePullSecret) PrefixTag(tag string) string {
	prefix := s.Host + s.Prefix

	if len(strings.Split(tag, "/")) == 1 {
		tag = "library/" + tag
	}

	if strings.HasPrefix(tag, prefix) {
		return tag
	}

	return prefix + tag
}

func (s ImagePullSecret) String() string {
	v := &url.URL{}
	v.Scheme = s.Name
	v.Host = s.Host
	v.Path = s.Prefix

	if s.Username != "" || s.Password != "" {
		v.User = url.UserPassword(s.Username, s.Password)
	}

	return v.String()
}

func (s ImagePullSecret) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *ImagePullSecret) UnmarshalText(data []byte) error {
	imagePullSecret, err := ParseImagePullSecret(string(data))
	if err != nil {
		return err
	}
	*s = *imagePullSecret
	return nil
}

func (s ImagePullSecret) RegistryAuth() string {
	authConfig := AuthConfig{Username: s.Username, Password: s.Password, ServerAddress: s.Host}
	b, _ := json.Marshal(authConfig)
	return base64.StdEncoding.EncodeToString(b)
}

func (s ImagePullSecret) DockerConfigJSON() []byte {
	v := struct {
		Auths map[string]AuthConfig `json:"auths"`
	}{
		Auths: map[string]AuthConfig{
			s.Host: {Username: s.Username, Password: s.Password},
		},
	}
	b, _ := json.Marshal(v)
	return b
}

type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
}
