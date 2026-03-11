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

type Config struct {
	Port       int          `json:"port"`
	Routes     []Route      `json:"routes"`
	Logs       []RequestLog `json:"logs,omitempty"`
	CORSEnabled bool        `json:"cors_enabled"`
	MaxLogs    int          `json:"max_logs"`
	ProxyURL   string       `json:"proxy_url,omitempty"` // forward unmatched requests
	mu         sync.Mutex
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
