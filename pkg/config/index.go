package config

import (
	"strings"
	"sync"
)

// RouteIndex provides fast route lookup
type RouteIndex struct {
	mu        sync.RWMutex
	exact     map[string]*Route // method:path -> route
	param     map[string][]*Route // method:prefix -> routes with params
	wildcard  map[string][]*Route // method:prefix -> routes with wildcards
	all       []*Route // routes with ALL method
}

// NewRouteIndex creates a new route index
func NewRouteIndex() *RouteIndex {
	return &RouteIndex{
		exact:    make(map[string]*Route),
		param:    make(map[string][]*Route),
		wildcard: make(map[string][]*Route),
		all:      make([]*Route, 0),
	}
}

// Add adds a route to the index
func (ri *RouteIndex) Add(route *Route) {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	key := route.Method + ":" + route.Path

	if strings.Contains(route.Path, "*") {
		// Wildcard route: /api/*
		prefix := strings.TrimSuffix(route.Path, "*")
		wildKey := route.Method + ":" + prefix
		ri.wildcard[wildKey] = append(ri.wildcard[wildKey], route)
	} else if strings.Contains(route.Path, ":") {
		// Param route: /users/:id
		parts := strings.Split(strings.Trim(route.Path, "/"), "/")
		prefix := ""
		for _, part := range parts {
			if strings.HasPrefix(part, ":") {
				break
			}
			if prefix != "" {
				prefix += "/"
			}
			prefix += part
		}
		paramKey := route.Method + ":" + prefix
		ri.param[paramKey] = append(ri.param[paramKey], route)
	} else if route.Method == "ALL" {
		// ALL method
		ri.all = append(ri.all, route)
	} else {
		// Exact route: /users
		ri.exact[key] = route
	}
}

// Remove removes a route from the index
func (ri *RouteIndex) Remove(route *Route) {
	ri.mu.Lock()
	defer ri.mu.Unlock()

	key := route.Method + ":" + route.Path

	// Try exact
	delete(ri.exact, key)

	// Try param
	if strings.Contains(route.Path, ":") {
		parts := strings.Split(strings.Trim(route.Path, "/"), "/")
		prefix := ""
		for _, part := range parts {
			if strings.HasPrefix(part, ":") {
				break
			}
			if prefix != "" {
				prefix += "/"
			}
			prefix += part
		}
		paramKey := route.Method + ":" + prefix
		routes := ri.param[paramKey]
		for i, r := range routes {
			if r.ID == route.ID {
				ri.param[paramKey] = append(routes[:i], routes[i+1:]...)
				break
			}
		}
	}

	// Try wildcard
	if strings.Contains(route.Path, "*") {
		prefix := strings.TrimSuffix(route.Path, "*")
		wildKey := route.Method + ":" + prefix
		routes := ri.wildcard[wildKey]
		for i, r := range routes {
			if r.ID == route.ID {
				ri.wildcard[wildKey] = append(routes[:i], routes[i+1:]...)
				break
			}
		}
	}

	// Try ALL
	if route.Method == "ALL" {
		for i, r := range ri.all {
			if r.ID == route.ID {
				ri.all = append(ri.all[:i], ri.all[i+1:]...)
				break
			}
		}
	}
}

// Find finds a matching route
func (ri *RouteIndex) Find(method, path string, headers map[string]string, body string) *Route {
	ri.mu.RLock()
	defer ri.mu.RUnlock()

	// 1. Try exact match
	key := method + ":" + path
	if route, ok := ri.exact[key]; ok {
		if matchConditions(route, headers, body) {
			return route
		}
	}

	// 2. Try param routes
	parts := strings.Split(strings.Trim(path, "/"), "/")
	prefix := ""
	for i, part := range parts {
		if i > 0 {
			prefix += "/"
		}
		prefix += part

		paramKey := method + ":" + prefix
		if routes, ok := ri.param[paramKey]; ok {
			for _, route := range routes {
				if MatchPath(route.Path, path) && matchConditions(route, headers, body) {
					return route
				}
			}
		}
	}

	// 3. Try wildcard routes
	for i := len(parts); i > 0; i-- {
		prefix := strings.Join(parts[:i], "/") + "/"
		wildKey := method + ":" + prefix
		if routes, ok := ri.wildcard[wildKey]; ok {
			for _, route := range routes {
				if MatchPath(route.Path, path) && matchConditions(route, headers, body) {
					return route
				}
			}
		}
	}

	// 4. Try ALL method routes
	for _, route := range ri.all {
		if MatchPath(route.Path, path) && matchConditions(route, headers, body) {
			return route
		}
	}

	return nil
}

// matchConditions checks if route conditions match the request
func matchConditions(route *Route, headers map[string]string, body string) bool {
	// Check header conditions
	for k, v := range route.MatchHeaders {
		if headers[k] != v {
			return false
		}
	}

	// Check body conditions
	if route.MatchBody != "" && !strings.Contains(body, route.MatchBody) {
		return false
	}

	return true
}

// Count returns total number of routes
func (ri *RouteIndex) Count() int {
	ri.mu.RLock()
	defer ri.mu.RUnlock()

	count := len(ri.exact)
	for _, routes := range ri.param {
		count += len(routes)
	}
	for _, routes := range ri.wildcard {
		count += len(routes)
	}
	count += len(ri.all)
	return count
}
