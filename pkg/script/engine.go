package script

import (
	"encoding/json"

	"github.com/dop251/goja"
)

type Engine struct {
	vm *goja.Runtime
}

type ScriptContext struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Params  map[string]string `json:"params"`
	Query   map[string]string `json:"query"`
}

func New() *Engine {
	return &Engine{
		vm: goja.New(),
	}
}

// Execute runs a JavaScript snippet and returns the response
// The script should return an object with status, body, headers
func (e *Engine) Execute(script string, ctx ScriptContext) (status int, body string, headers map[string]string, err error) {
	vm := e.vm

	// Set up context
	vm.Set("method", ctx.Method)
	vm.Set("path", ctx.Path)
	vm.Set("headers", ctx.Headers)
	vm.Set("body", ctx.Body)
	vm.Set("params", ctx.Params)
	vm.Set("query", ctx.Query)

	// Helper to parse JSON body
	vm.Set("parseBody", func(call goja.FunctionCall) goja.Value {
		var result interface{}
		json.Unmarshal([]byte(ctx.Body), &result)
		return vm.ToValue(result)
	})

	// Helper to set response
	var responseResult map[string]interface{}
	vm.Set("respond", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			val := call.Arguments[0].Export()
			if m, ok := val.(map[string]interface{}); ok {
				responseResult = m
			}
		}
		return goja.Undefined()
	})

	// Wrap script in function
	wrapped := `(function() {
		` + script + `
	})()`

	_, err = vm.RunString(wrapped)
	if err != nil {
		return 200, "", nil, err
	}

	// Check if respond() was called
	if responseResult != nil {
		if statusVal, ok := responseResult["status"].(int); ok {
			status = statusVal
		} else {
			status = 200
		}

		if bodyVal, ok := responseResult["body"].(string); ok {
			body = bodyVal
		} else if bodyVal, ok := responseResult["body"]; ok {
			data, _ := json.Marshal(bodyVal)
			body = string(data)
		}

		if headersVal, ok := responseResult["headers"].(map[string]interface{}); ok {
			headers = make(map[string]string)
			for k, v := range headersVal {
				headers[k] = v.(string)
			}
		}

		return status, body, headers, nil
	}

	// Fallback: use return value
	return 200, "", nil, nil
}

// Validate checks if a script is valid JavaScript
func (e *Engine) Validate(script string) error {
	_, err := e.vm.RunString(script)
	return err
}
