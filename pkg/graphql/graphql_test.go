package graphql

import (
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantOp  string
		wantErr bool
	}{
		{
			name:   "simple query",
			body:   `{"query": "{ users { id name } }"}`,
			wantOp: "",
		},
		{
			name:   "named query",
			body:   `{"query": "query GetUsers { users { id } }", "operationName": "GetUsers"}`,
			wantOp: "GetUsers",
		},
		{
			name:   "with variables",
			body:   `{"query": "query($id: ID!){ user(id: $id) { name } }", "variables": {"id": 1}}`,
			wantOp: "",
		},
		{
			name:    "invalid json",
			body:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := ParseRequest(tt.body)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && req.OperationName != tt.wantOp {
				t.Errorf("OperationName = %q, want %q", req.OperationName, tt.wantOp)
			}
		})
	}
}

func TestExtractOperation(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{
			query: "{ users { id } }",
			want:  "query",
		},
		{
			query: "query GetUsers { users { id } }",
			want:  "GetUsers",
		},
		{
			query: "mutation CreateUser { createUser { id } }",
			want:  "CreateUser",
		},
		{
			query: "subscription OnUserCreated { userCreated { id } }",
			want:  "OnUserCreated",
		},
		{
			query: "query($id: ID!) { user(id: $id) { name } }",
			want:  "query",
		},
		{
			query: "# comment\nquery GetPosts { posts { title } }",
			want:  "GetPosts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := ExtractOperation(tt.query)
			if got != tt.want {
				t.Errorf("ExtractOperation(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchOperation(t *testing.T) {
	handler := &MockHandler{
		OperationName: "GetUsers",
		Response:      map[string]interface{}{"users": []interface{}{}},
	}

	// Test specific match
	if MatchOperation(&Request{OperationName: "GetUsers"}, handler) != true {
		t.Error("Expected match for GetUsers")
	}

	if MatchOperation(&Request{OperationName: "GetPosts"}, handler) != false {
		t.Error("Expected no match for GetPosts")
	}

	// Test catch-all (empty operation name)
	catchAll := &MockHandler{OperationName: ""}
	if MatchOperation(&Request{OperationName: "Anything"}, catchAll) != true {
		t.Error("Expected catch-all to match any operation")
	}
}

func TestBuildResponse(t *testing.T) {
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice"},
		},
	}

	resp := BuildResponse(data, nil)
	if resp.Data == nil {
		t.Error("Expected data in response")
	}
	if len(resp.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(resp.Errors))
	}

	// Test with errors
	respErr := BuildResponse(data, []string{"Error 1", "Error 2"})
	if len(respErr.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(respErr.Errors))
	}
	if respErr.Errors[0].Message != "Error 1" {
		t.Errorf("First error message = %q, want %q", respErr.Errors[0].Message, "Error 1")
	}
}

func TestRequestWithVariables(t *testing.T) {
	body := `{"query": "query($id: ID!){ user(id: $id) { name } }", "variables": {"id": 123}}`
	req, err := ParseRequest(body)
	if err != nil {
		t.Fatalf("ParseRequest failed: %v", err)
	}

	if req.Variables == nil {
		t.Error("Expected variables to be parsed")
	}

	id, ok := req.Variables["id"].(float64)
	if !ok || id != 123 {
		t.Errorf("Expected id=123, got %v", req.Variables["id"])
	}
}
