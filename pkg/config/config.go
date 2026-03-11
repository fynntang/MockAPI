package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Route struct {
	ID          string            `json:"id"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Status      int               `json:"status"`
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers,omitempty"`
	Delay       int               `json:"delay_ms,omitempty"` // response delay in ms
	Description string            `json:"description,omitempty"`
}

type Config struct {
	Port   int     `json:"port"`
	Routes []Route `json:"routes"`
	mu     sync.Mutex
}

func Default() *Config {
	return &Config{
		Port:   8088,
		Routes: []Route{},
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
		if c.Routes[i].Method == method && c.Routes[i].Path == path {
			return &c.Routes[i]
		}
	}
	return nil
}
