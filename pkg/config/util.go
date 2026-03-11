package config

import (
	"fmt"
	"strings"
	"time"
)

// ParsePath replaces :param placeholders with actual values
func ParsePath(pattern, path string, body string) string {
	pParts := splitTrim(pattern)
	uParts := splitTrim(path)
	if len(uParts) != len(pParts) {
		return body
	}
	result := body
	for i, p := range pParts {
		if strings.HasPrefix(p, ":") {
			paramName := strings.TrimPrefix(p, ":")
			placeholder := fmt.Sprintf("{{%s}}", paramName)
			result = strings.ReplaceAll(result, placeholder, uParts[i])
		}
	}
	return result
}

func splitTrim(s string) []string {
	return strings.Split(strings.Trim(s, "/"), "/")
}

func Now() string {
	return time.Now().Format("15:04:05")
}
