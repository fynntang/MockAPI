package config

import (
	"strings"
)

// matchPath supports:
//   - exact: /users
//   - param: /users/:id
//   - wildcard: /users/*
func matchPath(pattern, path string) bool {
	pParts := strings.Split(strings.Trim(pattern, "/"), "/")
	uParts := strings.Split(strings.Trim(path, "/"), "/")

	for i := 0; i < len(pParts); i++ {
		if i >= len(uParts) {
			return false
		}
		if pParts[i] == "*" {
			return true
		}
		if pParts[i] == "" {
			continue
		}
		if strings.HasPrefix(pParts[i], ":") {
			continue
		}
		if pParts[i] != uParts[i] {
			return false
		}
	}

	return len(pParts) == len(uParts)
}
