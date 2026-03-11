package server

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"mockapi/pkg/config"
	"mockapi/pkg/script"
	"mockapi/pkg/swagger"
	"mockapi/pkg/ws"
)

type Server struct {
	cfg         *config.Config
	configFile  string
	mux         *http.ServeMux
	port        int
	https       bool
	certFile    string
	keyFile     string
	scriptEngine *script.Engine
	wsMock      *ws.MockWS
}

func New(cfg *config.Config, configFile string) *Server {
	s := &Server{
		cfg:         cfg,
		configFile:  configFile,
		mux:         http.NewServeMux(),
		port:        cfg.Port,
		scriptEngine: script.New(),
		wsMock:      ws.New(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) EnableHTTPS(certFile, keyFile string) {
	s.https = true
	s.certFile = certFile
	s.keyFile = keyFile
}

func (s *Server) registerRoutes() {
	// Web UI
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/web/app.js", s.handleJS)
	s.mux.HandleFunc("/web/style.css", s.handleCSS)

	// API management
	s.mux.HandleFunc("/_api/routes", s.handleAPIRoutes)
	s.mux.HandleFunc("/_api/logs", s.handleAPILogs)
	s.mux.HandleFunc("/_api/import", s.handleImport)
	s.mux.HandleFunc("/_api/export", s.handleExport)
	s.mux.HandleFunc("/_api/clear-logs", s.handleClearLogs)
	s.mux.HandleFunc("/_api/templates", s.handleTemplates)
	s.mux.HandleFunc("/_api/config", s.handleAPIConfig)
	s.mux.HandleFunc("/_api/import-swagger", s.handleImportSwagger)
	s.mux.HandleFunc("/_api/ws", s.handleAPIWS)

	// Mock routes
	s.mux.HandleFunc("/mock/", s.handleMock)
	
	// WebSocket routes
	s.mux.HandleFunc("/ws/", s.handleWS)
}

// --- Static files ---

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s.cors(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, indexHTML)
}

func (s *Server) handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	fmt.Fprint(w, appJS)
}

func (s *Server) handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	fmt.Fprint(w, styleCSS)
}

// --- Route CRUD ---

func (s *Server) handleAPIRoutes(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	if r.Method == http.MethodOptions {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(s.cfg.Routes)

	case http.MethodPost:
		var route config.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		route.ID = generateID()
		s.cfg.AddRoute(route)
		s.cfg.Save(s.configFile)
		json.NewEncoder(w).Encode(route)

	case http.MethodPut:
		var route config.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.cfg.UpdateRoute(route)
		s.cfg.Save(s.configFile)
		json.NewEncoder(w).Encode(route)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		s.cfg.DeleteRoute(id)
		s.cfg.Save(s.configFile)
		w.WriteHeader(204)
	}
}

// --- Config ---

func (s *Server) handleAPIConfig(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(map[string]interface{}{
			"port":         s.cfg.Port,
			"cors_enabled": s.cfg.CORSEnabled,
			"proxy_url":    s.cfg.ProxyURL,
			"max_logs":     s.cfg.MaxLogs,
		})
	case http.MethodPut:
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if v, ok := body["cors_enabled"].(bool); ok {
			s.cfg.CORSEnabled = v
		}
		if v, ok := body["proxy_url"].(string); ok {
			s.cfg.ProxyURL = strings.TrimSuffix(v, "/")
		}
		if v, ok := body["max_logs"].(float64); ok {
			s.cfg.MaxLogs = int(v)
		}
		s.cfg.Save(s.configFile)
		w.WriteHeader(204)
	}
}

// --- Logs ---

func (s *Server) handleAPILogs(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.cfg.Logs)
}

func (s *Server) handleClearLogs(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	s.cfg.ClearLogs()
	s.cfg.Save(s.configFile)
	w.WriteHeader(204)
}

// --- Import / Export ---

func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	body, _ := io.ReadAll(r.Body)
	var routes []config.Route
	if err := json.Unmarshal(body, &routes); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), 400)
		return
	}

	for _, route := range routes {
		route.ID = generateID()
		s.cfg.AddRoute(route)
	}
	s.cfg.Save(s.configFile)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"imported": len(routes)})
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=mockapi-routes.json")
	json.NewEncoder(w).Encode(s.cfg.Routes)
}

// --- Templates ---

func (s *Server) handleTemplates(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

// --- Mock handler with proxy fallback ---

func (s *Server) handleMock(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	method := strings.ToUpper(r.Method)
	path := strings.TrimPrefix(r.URL.Path, "/mock")

	// Read request body for conditional matching & logging
	reqBodyBytes, _ := io.ReadAll(io.LimitReader(r.Body, 64*1024))
	reqBody := string(reqBodyBytes)

	// Restore body for proxy
	r.Body = io.NopCloser(strings.NewReader(reqBody))

	route := findRoute(s.cfg, method, path, r.Header, reqBody)

	var status int
	var body string
	var routeID string
	var proxied bool

	if route == nil {
		// Try proxy
		if s.cfg.ProxyURL != "" {
			proxied = s.proxyRequest(w, r, path)
			if !proxied {
				status = 502
				body = `{"error":"proxy failed","message":"backend unreachable"}`
			}
		}
		if !proxied {
			status = 404
			body = fmt.Sprintf(`{"error":"no mock found","method":"%s","path":"%s","hint":"add a route via Web UI or /_api/routes"}`, method, path)
		}
	} else {
		routeID = route.ID
		
		// Check for script execution
		if route.Script != "" {
			ctx := script.ScriptContext{
				Method:  method,
				Path:    path,
				Headers: headersToMap(r.Header),
				Body:    reqBody,
				Params:  extractParams(route.Path, path),
				Query:   queryToMap(r.URL.Query()),
			}
			scriptStatus, scriptBody, scriptHeaders, err := s.scriptEngine.Execute(route.Script, ctx)
			if err == nil && scriptBody != "" {
				status = scriptStatus
				body = scriptBody
				for k, v := range scriptHeaders {
					w.Header().Set(k, v)
				}
			} else {
				// Fallback to static response
				status = route.Status
				body = config.ParsePath(route.Path, path, route.Body)
			}
		} else {
			status = route.Status
			body = config.ParsePath(route.Path, path, route.Body)
		}

		if route.Delay > 0 {
			time.Sleep(time.Duration(route.Delay) * time.Millisecond)
		}

		for k, v := range route.Headers {
			w.Header().Set(k, v)
		}
	}

	if !proxied {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}

	// Log
	elapsed := time.Since(start)
	logStatus := status
	if proxied {
		logStatus = 200 // simplified
	}
	s.cfg.AddLog(config.RequestLog{
		ID:        generateID(),
		Method:    method,
		Path:      path,
		Status:    logStatus,
		Delay:     int(elapsed.Milliseconds()),
		RouteID:   routeID,
		Timestamp: config.Now(),
		Body:      truncate(reqBody, 500),
		Proxied:   proxied,
	})
}

// --- Proxy ---

func (s *Server) proxyRequest(w http.ResponseWriter, r *http.Request, path string) bool {
	target, err := url.Parse(s.cfg.ProxyURL)
	if err != nil {
		return false
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		req.URL.Path = path
		req.Host = target.Host
	}

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	proxy.ServeHTTP(w, r)
	return true
}

// --- CORS ---

func (s *Server) cors(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.CORSEnabled {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// --- Start ---

func (s *Server) Start() error {
	protocol := "http"
	if s.https {
		protocol = "https"
	}

	proxyNote := ""
	if s.cfg.ProxyURL != "" {
		proxyNote = fmt.Sprintf("   Proxy:    %s\n", s.cfg.ProxyURL)
	}

	fmt.Printf(`
🦞 MockAPI Server running!
   Web UI:    %s://localhost:%d/
   Mock base: %s://localhost:%d/mock/
   API:       %s://localhost:%d/_api/routes
   Config:    %s
%s
`, protocol, s.port, protocol, s.port, protocol, s.port, s.configFile, proxyNote)

	if s.https {
		cfg := &tls.Config{MinVersion: tls.VersionTLS12}
		server := &http.Server{
			Addr:      fmt.Sprintf(":%d", s.port),
			Handler:   s.mux,
			TLSConfig: cfg,
		}
		return server.ListenAndServeTLS(s.certFile, s.keyFile)
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

// --- Helpers ---

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// --- FindRoute with conditions (server-level) ---

func findRoute(cfg *config.Config, method, path string, headers http.Header, reqBody string) *config.Route {
	// Use the simple FindRoute first for path/method match, then filter by conditions
	// We iterate directly since we need conditional checks
	for _, r := range cfg.Routes {
		if r.Method != method && r.Method != "ALL" {
			continue
		}
		if !config.MatchPath(r.Path, path) {
			continue
		}
		if !matchHeaders(r.MatchHeaders, headers) {
			continue
		}
		if r.MatchBody != "" && !strings.Contains(reqBody, r.MatchBody) {
			continue
		}
		cp := r
		return &cp
	}
	return nil
}

func matchHeaders(required map[string]string, actual http.Header) bool {
	for k, v := range required {
		if actual.Get(k) != v {
			return false
		}
	}
	return true
}

// --- Preset templates ---

var templates = []config.Route{
	{
		Method: "GET", Path: "/users", Status: 200,
		Description: "User list",
		Body: `[
  {"id": 1, "name": "Alice", "email": "alice@example.com"},
  {"id": 2, "name": "Bob", "email": "bob@example.com"}
]`,
	},
	{
		Method: "GET", Path: "/users/:id", Status: 200,
		Description: "User by ID (with path param)",
		Body: `{"id": {{id}}, "name": "Alice", "email": "alice@example.com"}`,
	},
	{
		Method: "POST", Path: "/users", Status: 201,
		Description: "Create user",
		Body: `{"id": 3, "name": "New User", "created": true}`,
	},
	{
		Method: "PUT", Path: "/users/:id", Status: 200,
		Description: "Update user",
		Body: `{"id": {{id}}, "name": "Updated", "updated": true}`,
	},
	{
		Method: "DELETE", Path: "/users/:id", Status: 204,
		Description: "Delete user",
		Body: ``,
	},
	{
		Method: "GET", Path: "/posts", Status: 200,
		Description: "Post list",
		Body: `[
  {"id": 1, "title": "First Post", "author": "Alice"},
  {"id": 2, "title": "Second Post", "author": "Bob"}
]`,
	},
	{
		Method: "GET", Path: "/posts/:id", Status: 200,
		Description: "Post by ID",
		Body: `{"id": {{id}}, "title": "First Post", "content": "Hello world!", "author": "Alice"}`,
	},
	{
		Method: "GET", Path: "/posts/:id/comments", Status: 200,
		Description: "Post comments",
		Body: `[
  {"id": 1, "postId": {{id}}, "text": "Great post!", "author": "Bob"},
  {"id": 2, "postId": {{id}}, "text": "Thanks!", "author": "Alice"}
]`,
	},
	{
		Method: "POST", Path: "/auth/login", Status: 200,
		Description: "Login success",
		Headers: map[string]string{"X-Auth-Token": "mock-token-abc123"},
		Body: `{"token": "mock-token-abc123", "user": {"id": 1, "name": "Alice"}}`,
	},
	{
		Method: "POST", Path: "/auth/login", Status: 401,
		Description: "Login fail (when body contains 'wrong')",
		MatchBody: "wrong",
		Body: `{"error": "invalid credentials"}`,
	},
	{
		Method: "GET", Path: "/health", Status: 200,
		Description: "Health check",
		Body: `{"status": "ok", "uptime": 3600}`,
	},
	{
		Method: "ALL", Path: "/slow/*", Status: 200,
		Description: "Slow response (2s delay)",
		Delay: 2000,
		Body: `{"message": "This took 2 seconds!"}`,
	},
	{
		Method: "GET", Path: "/api/*", Status: 404,
		Description: "Generic 404 for /api/*",
		Body: `{"error": "not found", "code": 404}`,
	},
	{
		Method: "GET", Path: "/items", Status: 200,
		Description: "Paginated list",
		Body: `{"items": [{"id":1},{"id":2}],"total": 42,"page": 1,"per_page": 20}`,
	},
	{
		Method: "GET", Path: "/random", Status: 200,
		Description: "Dynamic response (JS script)",
		Script: `// Dynamic response with JavaScript
var random = Math.floor(Math.random() * 1000);
respond({
  status: 200,
  body: JSON.stringify({ id: random, timestamp: Date.now() })
});`,
	},
}

// --- Swagger Import ---

func (s *Server) handleImportSwagger(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	data, _ := io.ReadAll(r.Body)
	routes, err := swagger.ParseOpenAPI(data)
	if err != nil {
		http.Error(w, "Failed to parse OpenAPI: "+err.Error(), 400)
		return
	}

	for _, route := range routes {
		route.ID = generateID()
		s.cfg.AddRoute(route)
	}
	s.cfg.Save(s.configFile)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"imported": len(routes),
		"routes":   routes,
	})
}

// --- WebSocket API ---

func (s *Server) handleAPIWS(w http.ResponseWriter, r *http.Request) {
	s.cors(w, r)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		handlers := s.wsMock.ListHandlers()
		json.NewEncoder(w).Encode(handlers)

	case http.MethodPost:
		var h ws.WSHandler
		if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.wsMock.AddHandler(h)
		json.NewEncoder(w).Encode(h)

	case http.MethodDelete:
		path := r.URL.Query().Get("path")
		_ = path // for now
		w.WriteHeader(204)
	}
}

// --- WebSocket Handler ---

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	s.wsMock.HandleWS(w, r, s.scriptEngine)
}

// --- Helpers for script context ---

func headersToMap(h http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}

func extractParams(pattern, path string) map[string]string {
	params := make(map[string]string)
	pParts := strings.Split(strings.Trim(pattern, "/"), "/")
	uParts := strings.Split(strings.Trim(path, "/"), "/")
	for i, p := range pParts {
		if i < len(uParts) && strings.HasPrefix(p, ":") {
			params[strings.TrimPrefix(p, ":")] = uParts[i]
		}
	}
	return params
}

func queryToMap(q url.Values) map[string]string {
	m := make(map[string]string)
	for k, v := range q {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}
