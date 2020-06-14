package converter

import (
	"encoding/json"
	"net/url"
	"testing"
)

func TestToExternalServiceEntity(t *testing.T) {
	u, _ := url.Parse("redis://127.0.0.1:6379")
	s := ToExternalServiceEntity(u)
	data, _ := json.MarshalIndent(s, "", "  ")

	t.Log(string(data))
}
