package swagger

import (
	"strings"
	"testing"

	"mockapi/pkg/config"
)

func TestParseOpenAPIv3(t *testing.T) {
	spec := `{
  "openapi": "3.0.0",
  "info": {"title": "Test API", "version": "1.0.0"},
  "paths": {
    "/users": {
      "get": {
        "summary": "Get all users",
        "responses": {
          "200": {
            "description": "Success",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {"type": "object", "properties": {"id": {"type": "integer"}, "name": {"type": "string"}}}
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create user",
        "responses": {
          "201": {"description": "Created"}
        }
      }
    },
    "/users/{id}": {
      "get": {
        "summary": "Get user by ID",
        "parameters": [{"name": "id", "in": "path", "required": true}],
        "responses": {
          "200": {"description": "Success"}
        }
      }
    }
  }
}`

	routes, err := ParseOpenAPI([]byte(spec))
	if err != nil {
		t.Fatalf("ParseOpenAPI failed: %v", err)
	}

	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Find the GET /users route (order not guaranteed due to map iteration)
	var getUsersRoute *config.Route
	for i := range routes {
		if routes[i].Method == "GET" && routes[i].Path == "/users" {
			getUsersRoute = &routes[i]
			break
		}
	}

	if getUsersRoute == nil {
		t.Fatalf("GET /users route not found")
	}

	if getUsersRoute.Description != "Get all users" {
		t.Errorf("Expected description 'Get all users', got %s", getUsersRoute.Description)
	}
}

func TestParseOpenAPIv2(t *testing.T) {
	spec := `{
  "swagger": "2.0",
  "info": {"title": "Test API", "version": "1.0.0"},
  "paths": {
    "/posts": {
      "get": {
        "summary": "Get posts",
        "responses": {
          "200": {"description": "Success"}
        }
      }
    }
  }
}`

	routes, err := ParseOpenAPI([]byte(spec))
	if err != nil {
		t.Fatalf("ParseOpenAPI failed: %v", err)
	}

	if len(routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(routes))
	}

	if routes[0].Method != "GET" {
		t.Errorf("Expected GET, got %s", routes[0].Method)
	}
}

func TestConvertPath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"/users/{id}", "/users/:id"},
		{"/posts/{postId}/comments/{commentId}", "/posts/:postId/comments/:commentId"},
		{"/users", "/users"},
	}

	for _, tt := range tests {
		got := convertPath(tt.input)
		if got != tt.want {
			t.Errorf("convertPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFindBestResponse(t *testing.T) {
	responses := map[string]Response{
		"404": {Description: "Not found"},
		"200": {Description: "Success"},
		"500": {Description: "Error"},
	}

	status, resp := findBestResponse(responses)
	if status != 200 {
		t.Errorf("Expected status 200, got %d", status)
	}
	if resp.Description != "Success" {
		t.Errorf("Expected 'Success', got %s", resp.Description)
	}
}

func TestSchemaToJSON(t *testing.T) {
	spec := OpenAPI{}

	// Test object schema
	objSchema := &Schema{
		Type: "object",
		Properties: map[string]Schema{
			"id":   {Type: "integer"},
			"name": {Type: "string"},
		},
	}

	jsonStr := schemaToJSON(objSchema, spec)
	if !strings.Contains(jsonStr, "id") {
		t.Errorf("Expected 'id' in JSON, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, "name") {
		t.Errorf("Expected 'name' in JSON, got %s", jsonStr)
	}

	// Test array schema
	arrSchema := &Schema{
		Type: "array",
		Items: &Schema{Type: "string"},
	}

	jsonStr = schemaToJSON(arrSchema, spec)
	if !strings.Contains(jsonStr, "[") {
		t.Errorf("Expected array in JSON, got %s", jsonStr)
	}
}

func TestInjectPathParams(t *testing.T) {
	body := `{}`
	result := injectPathParams("/users/:id", body)

	if !strings.Contains(result, "{{id}}") {
		t.Errorf("Expected {{id}} placeholder in result, got %s", result)
	}
}

func TestParseOpenAPIWithExample(t *testing.T) {
	spec := `{
  "openapi": "3.0.0",
  "paths": {
    "/health": {
      "get": {
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {"type": "string", "example": "ok"},
                    "uptime": {"type": "integer", "example": 3600}
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}`

	routes, err := ParseOpenAPI([]byte(spec))
	if err != nil {
		t.Fatalf("ParseOpenAPI failed: %v", err)
	}

	if len(routes) != 1 {
		t.Fatalf("Expected 1 route, got %d", len(routes))
	}

	// Check if example values are used
	if !strings.Contains(routes[0].Body, "ok") {
		t.Errorf("Expected example value 'ok' in body, got %s", routes[0].Body)
	}
}
