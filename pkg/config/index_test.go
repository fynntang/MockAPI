package config

import (
	"testing"
)

func BenchmarkFindRouteLinear(b *testing.B) {
	cfg := Default()
	
	// Add 100 routes
	for i := 0; i < 100; i++ {
		cfg.AddRoute(Route{
			ID:     string(rune(i)),
			Method: "GET",
			Path:   "/api/resource" + string(rune(i)),
			Status: 200,
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg.FindRoute("GET", "/api/resource50")
	}
}

func BenchmarkFindRouteIndexed(b *testing.B) {
	idx := NewRouteIndex()
	
	// Add 100 routes
	for i := 0; i < 100; i++ {
		idx.Add(&Route{
			ID:     string(rune(i)),
			Method: "GET",
			Path:   "/api/resource" + string(rune(i)),
			Status: 200,
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Find("GET", "/api/resource50", nil, "")
	}
}

func BenchmarkFindRouteWithParams(b *testing.B) {
	idx := NewRouteIndex()
	
	// Add routes with params
	for i := 0; i < 100; i++ {
		idx.Add(&Route{
			ID:     string(rune(i)),
			Method: "GET",
			Path:   "/users/:id/posts/:postId",
			Status: 200,
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Find("GET", "/users/123/posts/456", nil, "")
	}
}

func BenchmarkFindRouteWildcard(b *testing.B) {
	idx := NewRouteIndex()
	
	// Add wildcard routes
	for i := 0; i < 50; i++ {
		idx.Add(&Route{
			ID:     string(rune(i)),
			Method: "GET",
			Path:   "/api/v" + string(rune(i)) + "/*",
			Status: 200,
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Find("GET", "/api/v5/resources/123", nil, "")
	}
}
