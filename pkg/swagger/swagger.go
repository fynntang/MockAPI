package swagger

import (
	"encoding/json"
	"fmt"
	"strings"

	"mockapi/pkg/config"
	"gopkg.in/yaml.v3"
)

type OpenAPI struct {
	Swagger string                 `yaml:"swagger" json:"swagger"`
	OpenAPI string                 `yaml:"openapi" json:"openapi"`
	Paths   map[string]PathItem    `yaml:"paths" json:"paths"`
	Components Components          `yaml:"components" json:"components"`
	Definitions map[string]Schema  `yaml:"definitions" json:"definitions"`
}

type PathItem struct {
	Get    *Operation `yaml:"get" json:"get"`
	Post   *Operation `yaml:"post" json:"post"`
	Put    *Operation `yaml:"put" json:"put"`
	Patch  *Operation `yaml:"patch" json:"patch"`
	Delete *Operation `yaml:"delete" json:"delete"`
}

type Operation struct {
	Summary     string                 `yaml:"summary" json:"summary"`
	Description string                 `yaml:"description" json:"description"`
	OperationID string                 `yaml:"operationId" json:"operationId"`
	Responses   map[string]Response    `yaml:"responses" json:"responses"`
	Parameters  []Parameter            `yaml:"parameters" json:"parameters"`
	RequestBody *RequestBody           `yaml:"requestBody" json:"requestBody"`
}

type Response struct {
	Description string      `yaml:"description" json:"description"`
	Content     interface{} `yaml:"content" json:"content"`
	Schema      interface{} `yaml:"schema" json:"schema"`
}

type Parameter struct {
	Name        string      `yaml:"name" json:"name"`
	In          string      `yaml:"in" json:"in"`
	Description string      `yaml:"description" json:"description"`
	Required    bool        `yaml:"required" json:"required"`
	Schema      interface{} `yaml:"schema" json:"schema"`
}

type RequestBody struct {
	Content interface{} `yaml:"content" json:"content"`
}

type Components struct {
	Schemas map[string]Schema `yaml:"schemas" json:"schemas"`
}

type Schema struct {
	Type       string             `yaml:"type" json:"type"`
	Properties map[string]Schema  `yaml:"properties" json:"properties"`
	Items      *Schema            `yaml:"items" json:"items"`
	Ref        string             `yaml:"$ref" json:"$ref"`
	Example    interface{}        `yaml:"example" json:"example"`
}

// ParseOpenAPI parses OpenAPI 2.0 or 3.x spec and generates mock routes
func ParseOpenAPI(data []byte) ([]config.Route, error) {
	var spec OpenAPI
	
	// Try JSON first
	if err := json.Unmarshal(data, &spec); err != nil {
		// Try YAML
		if err := yaml.Unmarshal(data, &spec); err != nil {
			return nil, fmt.Errorf("failed to parse OpenAPI spec: %v", err)
		}
	}

	var routes []config.Route

	for path, pathItem := range spec.Paths {
		ops := []struct {
			method string
			op     *Operation
		}{
			{"GET", pathItem.Get},
			{"POST", pathItem.Post},
			{"PUT", pathItem.Put},
			{"PATCH", pathItem.Patch},
			{"DELETE", pathItem.Delete},
		}

		for _, op := range ops {
			if op.op == nil {
				continue
			}

			route := config.Route{
				Method:      op.method,
				Path:        convertPath(path),
				Description: op.op.Summary,
			}

			// Find best response (200, 201, 204, or first)
			status, resp := findBestResponse(op.op.Responses)
			route.Status = status

			// Generate response body
			route.Body = generateBody(resp, spec)

			// Extract path params for template
			if strings.Contains(route.Path, ":") {
				route.Body = injectPathParams(route.Path, route.Body)
			}

			routes = append(routes, route)
		}
	}

	return routes, nil
}

func convertPath(path string) string {
	// OpenAPI uses {param}, we use :param
	path = strings.ReplaceAll(path, "{", ":")
	path = strings.ReplaceAll(path, "}", "")
	return path
}

func findBestResponse(responses map[string]Response) (int, Response) {
	// Priority: 200 > 201 > 204 > first
	priority := []string{"200", "201", "204", "202", "default"}
	
	for _, code := range priority {
		if resp, ok := responses[code]; ok {
			status := 200
			fmt.Sscanf(code, "%d", &status)
			if code == "default" {
				status = 200
			}
			return status, resp
		}
	}
	
	// Return first response
	for code, resp := range responses {
		status := 200
		fmt.Sscanf(code, "%d", &status)
		return status, resp
	}
	
	return 200, Response{}
}

func generateBody(resp Response, spec OpenAPI) string {
	// Try to get schema from response
	var schema *Schema
	
	// OpenAPI 3.x: content -> application/json -> schema
	if resp.Content != nil {
		if content, ok := resp.Content.(map[string]interface{}); ok {
			if jsonContent, ok := content["application/json"].(map[string]interface{}); ok {
				if schemaData, ok := jsonContent["schema"].(map[string]interface{}); ok {
					schema = &Schema{}
					jsonData, _ := json.Marshal(schemaData)
					json.Unmarshal(jsonData, schema)
				}
			}
		}
	}
	
	// OpenAPI 2.0: schema directly
	if resp.Schema != nil {
		schema = &Schema{}
		jsonData, _ := json.Marshal(resp.Schema)
		json.Unmarshal(jsonData, schema)
	}

	if schema != nil {
		return schemaToJSON(schema, spec)
	}

	// Default empty object
	return `{}`
}

func schemaToJSON(schema *Schema, spec OpenAPI) string {
	if schema == nil {
		return "{}"
	}

	// Handle $ref
	if schema.Ref != "" {
		refName := strings.TrimPrefix(schema.Ref, "#/definitions/")
		refName = strings.TrimPrefix(refName, "#/components/schemas/")
		if def, ok := spec.Definitions[refName]; ok {
			return schemaToJSON(&def, spec)
		}
		if comp, ok := spec.Components.Schemas[refName]; ok {
			return schemaToJSON(&comp, spec)
		}
	}

	switch schema.Type {
	case "object":
		if schema.Example != nil {
			data, _ := json.Marshal(schema.Example)
			return string(data)
		}
		result := "{\n"
		props := []string{}
		for name, prop := range schema.Properties {
			props = append(props, fmt.Sprintf(`  "%s": %s`, name, schemaToJSON(&prop, spec)))
		}
		result += strings.Join(props, ",\n")
		result += "\n}"
		return result

	case "array":
		if schema.Items != nil {
			item := schemaToJSON(schema.Items, spec)
			return fmt.Sprintf("[\n  %s\n]", item)
		}
		return "[]"

	case "string":
		if schema.Example != nil {
			return fmt.Sprintf("%q", schema.Example)
		}
		return `"string"`

	case "integer", "number":
		if schema.Example != nil {
			return fmt.Sprintf("%v", schema.Example)
		}
		return "0"

	case "boolean":
		return "true"

	default:
		return "{}"
	}
}

func injectPathParams(path, body string) string {
	// Extract :param names
	params := []string{}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			params = append(params, strings.TrimPrefix(part, ":"))
		}
	}

	result := body
	for _, param := range params {
		placeholder := fmt.Sprintf("{{%s}}", param)
		if !strings.Contains(result, placeholder) {
			// Add param to body if not already referenced
			result = strings.Replace(result, "{}", fmt.Sprintf(`{"%s": {{%s}}}`, param, param), 1)
		}
	}

	return result
}
