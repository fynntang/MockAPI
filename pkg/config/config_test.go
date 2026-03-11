package config

import (
	"testing"
)

func TestMatchPath(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"/users", "/users", true},
		{"/users/:id", "/users/123", true},
		{"/users/:id", "/users/abc", true},
		{"/users/*", "/users/123/posts", true},
		{"/users/*", "/users", false}, // wildcard requires something after
		{"/users", "/posts", false},
		{"/users/:id", "/posts/123", false},
		{"/api/*", "/api/v1/users", true},
		{"/*", "/anything", true},
	}

	for _, tt := range tests {
		got := MatchPath(tt.pattern, tt.path)
		if got != tt.want {
			t.Errorf("MatchPath(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
		}
	}
}

func TestParsePath(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		body    string
		want    string
	}{
		{
			"/users/:id",
			"/users/123",
			`{"id": {{id}}}`,
			`{"id": 123}`,
		},
		{
			"/posts/:postId/comments/:commentId",
			"/posts/1/comments/2",
			`{"postId": {{postId}}, "commentId": {{commentId}}}`,
			`{"postId": 1, "commentId": 2}`,
		},
		{
			"/users",
			"/users",
			`{"static": true}`,
			`{"static": true}`,
		},
	}

	for _, tt := range tests {
		got := ParsePath(tt.pattern, tt.path, tt.body)
		if got != tt.want {
			t.Errorf("ParsePath(%q, %q, %q) = %q, want %q", tt.pattern, tt.path, tt.body, got, tt.want)
		}
	}
}

func TestConfigAddRoute(t *testing.T) {
	cfg := Default()
	
	r1 := Route{Method: "GET", Path: "/users", Status: 200}
	cfg.AddRoute(r1)
	
	if len(cfg.Routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(cfg.Routes))
	}
	
	r2 := Route{Method: "POST", Path: "/users", Status: 201}
	cfg.AddRoute(r2)
	
	if len(cfg.Routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(cfg.Routes))
	}
}

func TestConfigDeleteRoute(t *testing.T) {
	cfg := Default()
	
	r1 := Route{ID: "route1", Method: "GET", Path: "/users"}
	r2 := Route{ID: "route2", Method: "POST", Path: "/users"}
	cfg.AddRoute(r1)
	cfg.AddRoute(r2)
	
	cfg.DeleteRoute("route1")
	
	if len(cfg.Routes) != 1 {
		t.Errorf("Expected 1 route after delete, got %d", len(cfg.Routes))
	}
	
	if cfg.Routes[0].ID != "route2" {
		t.Errorf("Expected route2 to remain, got %s", cfg.Routes[0].ID)
	}
}

func TestConfigFindRoute(t *testing.T) {
	cfg := Default()
	
	cfg.AddRoute(Route{ID: "1", Method: "GET", Path: "/users"})
	cfg.AddRoute(Route{ID: "2", Method: "GET", Path: "/users/:id"})
	cfg.AddRoute(Route{ID: "3", Method: "ALL", Path: "/api/*"})
	
	tests := []struct {
		method string
		path   string
		wantID string
	}{
		{"GET", "/users", "1"},
		{"GET", "/users/123", "2"},
		{"POST", "/api/test", "3"},
		{"DELETE", "/unknown", ""}, // not found
	}
	
	for _, tt := range tests {
		route := cfg.FindRoute(tt.method, tt.path)
		if tt.wantID == "" {
			if route != nil {
				t.Errorf("FindRoute(%q, %q) expected nil, got %s", tt.method, tt.path, route.ID)
			}
		} else {
			if route == nil {
				t.Errorf("FindRoute(%q, %q) expected %s, got nil", tt.method, tt.path, tt.wantID)
			} else if route.ID != tt.wantID {
				t.Errorf("FindRoute(%q, %q) = %s, want %s", tt.method, tt.path, route.ID, tt.wantID)
			}
		}
	}
}