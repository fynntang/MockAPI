package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Route struct {
	ID           string            `json:"id"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	Status       int               `json:"status"`
	Body         string            `json:"body"`
	Headers      map[string]string `json:"headers,omitempty"`
	Delay        int               `json:"delay_ms,omitempty"`
	Description  string            `json:"description,omitempty"`
	// Conditional matching
	MatchHeaders map[string]string `json:"match_headers,omitempty"`
	MatchBody    string            `json:"match_body,omitempty"`
	// Script engine (dynamic response)
	Script string `json:"script,omitempty"` // JavaScript code for dynamic response
}

type RequestLog struct {
	ID        string `json:"id"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Delay     int    `json:"delay_ms"`
	RouteID   string `json:"route_id,omitempty"`
	Timestamp string `json:"timestamp"`
	Body      string `json:"body,omitempty"`
	Proxied   bool   `json:"proxied,omitempty"`
}

// GraphQLHandler represents a mock for GraphQL operations
type GraphQLHandler struct {
	ID            string      `json:"id"`
	OperationName string      `json:"operationName"`
	Response      interface{} `json:"response"`
	Delay         int         `json:"delay_ms,omitempty"`
	Description   string      `json:"description,omitempty"`
}

// GRPCHandler represents a mock for gRPC methods
type GRPCHandler struct {
	ID           string      `json:"id"`
	Service      string      `json:"service"`
	Method       string      `json:"method"`
	MockResponse interface{} `json:"mock_response"`
	Delay        int         `json:"delay_ms,omitempty"`
	Description  string      `json:"description,omitempty"`
}

// WSHandler represents a WebSocket mock handler
type WSHandler struct {
	Path            string   `json:"path"`
	Description     string   `json:"description,omitempty"`
	OnConnect       string   `json:"on_connect,omitempty"`
	OnMessage       string   `json:"on_message,omitempty"`
	AutoReply       string   `json:"auto_reply,omitempty"`
	Delay           int      `json:"delay_ms,omitempty"`
	StreamEnabled   bool     `json:"stream_enabled,omitempty"`
	StreamMessages  []string `json:"stream_messages,omitempty"`
	StreamInterval  int      `json:"stream_interval_ms,omitempty"`
	StreamRandom    bool     `json:"stream_random,omitempty"`
	StreamMinDelay  int      `json:"stream_min_delay_ms,omitempty"`
	StreamMaxDelay  int      `json:"stream_max_delay_ms,omitempty"`
	StreamLoop      bool     `json:"stream_loop,omitempty"`
	StreamFormat    string   `json:"stream_format,omitempty"`
}

type Config struct {
	Port        int              `json:"port"`
	Routes      []Route          `json:"routes"`
	Logs        []RequestLog     `json:"logs,omitempty"`
	CORSEnabled bool             `json:"cors_enabled"`
	MaxLogs     int              `json:"max_logs"`
	ProxyURL    string           `json:"proxy_url,omitempty"`
	GraphQL     []GraphQLHandler `json:"graphql,omitempty"`
	GRPC        []GRPCHandler    `json:"grpc,omitempty"`
	WebSocket   []WSHandler      `json:"websocket,omitempty"`
	mu          sync.Mutex
}

func Default() *Config {
	return &Config{
		Port:        8088,
		Routes:      []Route{},
		Logs:        []RequestLog{},
		CORSEnabled: true,
		MaxLogs:     500,
		ProxyURL:    "",
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}
	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Save(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (c *Config) AddRoute(r Route) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Routes = append(c.Routes, r)
}

func (c *Config) UpdateRoute(r Route) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, existing := range c.Routes {
		if existing.ID == r.ID {
			c.Routes[i] = r
			return
		}
	}
}

func (c *Config) DeleteRoute(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, r := range c.Routes {
		if r.ID == id {
			c.Routes = append(c.Routes[:i], c.Routes[i+1:]...)
			return
		}
	}
}

func (c *Config) FindRoute(method, path string) *Route {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.Routes {
		if c.Routes[i].Method != method && c.Routes[i].Method != "ALL" {
			continue
		}
		if MatchPath(c.Routes[i].Path, path) {
			r := c.Routes[i]
			return &r
		}
	}
	return nil
}

func (c *Config) AddLog(log RequestLog) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logs = append(c.Logs, log)
	if len(c.Logs) > c.MaxLogs && c.MaxLogs > 0 {
		c.Logs = c.Logs[len(c.Logs)-c.MaxLogs:]
	}
}

func (c *Config) ClearLogs() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logs = nil
}

func (c *Config) AddGraphQLHandler(h GraphQLHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.GraphQL = append(c.GraphQL, h)
}

func (c *Config) DeleteGraphQLHandler(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, h := range c.GraphQL {
		if h.ID == id {
			c.GraphQL = append(c.GraphQL[:i], c.GraphQL[i+1:]...)
			return
		}
	}
}

func (c *Config) FindGraphQLHandler(operationName string) *GraphQLHandler {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.GraphQL {
		if c.GraphQL[i].OperationName == operationName || c.GraphQL[i].OperationName == "" {
			return &c.GraphQL[i]
		}
	}
	return nil
}

func (c *Config) AddGRPCHandler(h GRPCHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.GRPC = append(c.GRPC, h)
}

func (c *Config) DeleteGRPCHandler(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, h := range c.GRPC {
		if h.ID == id {
			c.GRPC = append(c.GRPC[:i], c.GRPC[i+1:]...)
			return
		}
	}
}

func (c *Config) FindGRPCHandler(service, method string) *GRPCHandler {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.GRPC {
		if c.GRPC[i].Service == service && c.GRPC[i].Method == method {
			return &c.GRPC[i]
		}
	}
	return nil
}

func (c *Config) AddWSHandler(h WSHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.WebSocket = append(c.WebSocket, h)
}

func (c *Config) UpdateWSHandler(h WSHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, existing := range c.WebSocket {
		if existing.Path == h.Path {
			c.WebSocket[i] = h
			return
		}
	}
	// If not found, add it
	c.WebSocket = append(c.WebSocket, h)
}

func (c *Config) DeleteWSHandler(path string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, h := range c.WebSocket {
		if h.Path == path {
			c.WebSocket = append(c.WebSocket[:i], c.WebSocket[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Config) FindWSHandler(path string) *WSHandler {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.WebSocket {
		if c.WebSocket[i].Path == path {
			return &c.WebSocket[i]
		}
	}
	return nil
}
