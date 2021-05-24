package serving

import (
	"testing"

	"github.com/ghodss/yaml"
)

func TestCustomResourceDefinition(t *testing.T) {
	data, _ := yaml.Marshal(qserviceCustomResourceDefinition())

	t.Log(string(data))
}
