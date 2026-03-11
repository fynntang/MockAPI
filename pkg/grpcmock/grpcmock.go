package grpcmock

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Service represents a gRPC service
type Service struct {
	Name    string   `json:"name"`
	Package string   `json:"package"`
	Methods []*Method `json:"methods"`
}

// Method represents a gRPC method
type Method struct {
	Name        string      `json:"name"`
	InputType   string      `json:"input_type"`
	OutputType  string      `json:"output_type"`
	Description string      `json:"description,omitempty"`
	MockResponse interface{} `json:"mock_response,omitempty"`
	Delay       int         `json:"delay_ms,omitempty"`
}

// Request represents a gRPC request
type Request struct {
	Service string      `json:"service"`
	Method  string      `json:"method"`
	Data    interface{} `json:"data"`
}

// Response represents a gRPC response
type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
}

// ParseProto parses a simple proto definition (simplified)
func ParseProto(protoContent string) ([]Service, error) {
	var services []Service
	
	// Simple parser for basic proto3 syntax
	lines := strings.Split(protoContent, "\n")
	var currentService *Service
	var inService bool
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		
		// Parse service definition
		if strings.HasPrefix(line, "service ") {
			name := strings.TrimPrefix(line, "service ")
			name = strings.TrimSuffix(name, " {")
			currentService = &Service{
				Name:    name,
				Methods: []*Method{},
			}
			inService = true
			continue
		}
		
		// Parse rpc method
		if inService && strings.HasPrefix(line, "rpc ") {
			method := parseRPC(line)
			if method != nil && currentService != nil {
				currentService.Methods = append(currentService.Methods, method)
			}
		}
		
		// End of service
		if inService && line == "}" {
			if currentService != nil {
				services = append(services, *currentService)
			}
			inService = false
		}
	}
	
	return services, nil
}

func parseRPC(line string) *Method {
	// Parse: rpc MethodName(InputType) returns (OutputType)
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return nil
	}
	
	method := &Method{
		Name: parts[1],
	}
	
	// Extract input type
	for i, part := range parts {
		if part == "returns" && i > 2 {
			// Input type is before returns
			input := strings.Trim(parts[i-1], "()")
			method.InputType = input
			// Output type is after returns
			if i+1 < len(parts) {
				output := strings.Trim(parts[i+1], "()")
				method.OutputType = output
			}
			break
		}
	}
	
	return method
}

// GenerateMockResponse generates a mock response based on output type
func GenerateMockResponse(outputType string) interface{} {
	// Simple mock generator
	mockData := map[string]interface{}{
		"id":      1,
		"name":    "mock_item",
		"created": true,
		"timestamp": 1234567890,
	}
	
	// Handle common types
	switch {
	case strings.Contains(strings.ToLower(outputType), "user"):
		return map[string]interface{}{
			"id":    1,
			"name":  "John Doe",
			"email": "john@example.com",
		}
	case strings.Contains(strings.ToLower(outputType), "list"):
		return map[string]interface{}{
			"items": []interface{}{mockData, mockData},
			"total": 2,
		}
	default:
		return mockData
	}
}

// FormatMethodFullName returns the full gRPC method name
func (m *Method) FullName(serviceName, packageName string) string {
	return fmt.Sprintf("/%s.%s/%s", packageName, serviceName, m.Name)
}

// ToJSON converts services to JSON
func ToJSON(services []Service) (string, error) {
	data, err := json.MarshalIndent(services, "", "  ")
	return string(data), err
}
