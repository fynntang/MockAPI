package graphql

import (
	"encoding/json"
	"strings"
)

// Request represents a GraphQL request
type Request struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
}

// Response represents a GraphQL response
type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
}

// Error represents a GraphQL error
type Error struct {
	Message string `json:"message"`
}

// MockHandler represents a mock for a specific GraphQL operation
type MockHandler struct {
	OperationName string      `json:"operationName"`
	Response      interface{} `json:"response"`
	Delay         int         `json:"delay_ms,omitempty"`
}

// ParseRequest parses a GraphQL request from JSON body
func ParseRequest(body string) (*Request, error) {
	var req Request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// ExtractOperation extracts the operation name from a query
func ExtractOperation(query string) string {
	// Simple extraction: look for "queryName" or "mutationName"
	query = strings.TrimSpace(query)
	
	// Remove comments
	lines := strings.Split(query, "\n")
	var cleaned []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			cleaned = append(cleaned, line)
		}
	}
	query = strings.Join(cleaned, "\n")
	
	// Find operation type
	operationTypes := []string{"query", "mutation", "subscription"}
	for _, opType := range operationTypes {
		idx := strings.Index(strings.ToLower(query), opType)
		if idx == -1 {
			continue
		}
		
		// Extract name after operation type
		rest := strings.TrimSpace(query[idx+len(opType):])
		// Skip parentheses if present
		if strings.HasPrefix(rest, "(") {
			end := strings.Index(rest, ")")
			if end != -1 {
				rest = strings.TrimSpace(rest[end+1:])
			}
		}
		if strings.HasPrefix(rest, "{") {
			// Anonymous operation
			return opType
		}
		// Get first word as operation name
		fields := strings.Fields(rest)
		if len(fields) > 0 {
			return fields[0]
		}
	}
	
	// Default: anonymous query
	return "query"
}

// MatchOperation checks if a request matches a mock handler
func MatchOperation(req *Request, handler *MockHandler) bool {
	if handler.OperationName == "" {
		return true // Match all
	}
	
	// Match by operation name from request
	if req.OperationName != "" {
		return req.OperationName == handler.OperationName
	}
	
	// Extract from query
	opName := ExtractOperation(req.Query)
	return opName == handler.OperationName
}

// BuildResponse builds a GraphQL response
func BuildResponse(data interface{}, errors []string) Response {
	resp := Response{Data: data}
	for _, msg := range errors {
		resp.Errors = append(resp.Errors, Error{Message: msg})
	}
	return resp
}