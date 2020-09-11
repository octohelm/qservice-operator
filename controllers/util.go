package controllers

import "strings"

func BaseHost(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts[1:], ".")
}
