package script

import (
	"testing"
)

func TestExecute(t *testing.T) {
	engine := New()

	tests := []struct {
		name       string
		script     string
		ctx        ScriptContext
		wantStatus int
		wantBody   string
	}{
		{
			name: "simple response",
			script: `respond({
				status: 200,
				body: JSON.stringify({message: "Hello"})
			})`,
			ctx:        ScriptContext{},
			wantStatus: 200,
			wantBody:   `{"message":"Hello"}`,
		},
		{
			name: "use path params",
			script: `respond({
				status: 200,
				body: JSON.stringify({id: params.id})
			})`,
			ctx: ScriptContext{
				Params: map[string]string{"id": "123"},
			},
			wantStatus: 200,
			wantBody:   `{"id":"123"}`,
		},
		{
			name: "use query params",
			script: `respond({
				status: 200,
				body: JSON.stringify({page: query.page})
			})`,
			ctx: ScriptContext{
				Query: map[string]string{"page": "5"},
			},
			wantStatus: 200,
			wantBody:   `{"page":"5"}`,
		},
		{
			name: "random number",
			script: `var num = Math.floor(Math.random() * 1000);
			respond({
				status: 200,
				body: JSON.stringify({random: num})
			})`,
			ctx:        ScriptContext{},
			wantStatus: 200,
		},
		{
			name: "use request body",
			script: `var data = parseBody();
			respond({
				status: 200,
				body: JSON.stringify({created: true, name: data.name})
			})`,
			ctx: ScriptContext{
				Body: `{"name": "Test"}`,
			},
			wantStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, body, _, err := engine.Execute(tt.script, tt.ctx)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if status != tt.wantStatus {
				t.Errorf("status = %d, want %d", status, tt.wantStatus)
			}
			if tt.wantBody != "" && body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	engine := New()

	// Test invalid syntax
	err := engine.Validate(`respond({status: 200, body: "OK"`)
	if err == nil {
		t.Error("Expected error for invalid syntax")
	}

	// Test empty script (should not error)
	err = engine.Validate(``)
	if err != nil {
		t.Errorf("Empty script should not error: %v", err)
	}
}

func TestContextVariables(t *testing.T) {
	engine := New()

	ctx := ScriptContext{
		Method:  "POST",
		Path:    "/users/123",
		Headers: map[string]string{"Authorization": "Bearer token"},
		Body:    `{"name": "Test"}`,
		Params:  map[string]string{"id": "123"},
		Query:   map[string]string{"include": "posts"},
	}

	script := `
		var result = {
			method: method,
			path: path,
			auth: headers.Authorization,
			name: parseBody().name,
			id: params.id,
			include: query.include
		};
		respond({
			status: 200,
			body: JSON.stringify(result)
		});
	`

	status, body, _, err := engine.Execute(script, ctx)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if status != 200 {
		t.Errorf("Expected status 200, got %d", status)
	}

	// Verify all context variables were accessible
	if !contains(body, "POST") {
		t.Errorf("Expected 'POST' in body, got %s", body)
	}
	if !contains(body, "/users/123") {
		t.Errorf("Expected path in body, got %s", body)
	}
	if !contains(body, "Bearer token") {
		t.Errorf("Expected auth header in body, got %s", body)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
