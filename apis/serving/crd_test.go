package serving

import (
	"testing"

	"github.com/ghodss/yaml"
)

func TestCustomResourceDefinition(t *testing.T) {
	data, _ := yaml.Marshal(QServiceCustomResourceDefinition())

	t.Log(string(data))
}
