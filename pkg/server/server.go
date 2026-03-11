package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"mockapi/pkg/config"
)

type Server struct {
	cfg        *config.Config
	configFile string
	mux        *http.ServeMux
	port       int
}

func New(cfg *config.Config, configFile string) *Server {
	s := &Server{
		cfg:        cfg,
		configFile: configFile,
		mux:        http.NewServeMux(),
		port:       cfg.Port,
	}
	s.registerRoutes()
	return s
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

	// Mock routes - catch all
	s.mux.HandleFunc("/mock/", s.handleMock)
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

// --- Mock handler ---

func (s *Server) handleMock(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	method := strings.ToUpper(r.Method)
	path := strings.TrimPrefix(r.URL.Path, "/mock")

	route := s.cfg.FindRoute(method, path)

	var status int
	var body string
	var routeID string

	if route == nil {
		status = 404
		body = fmt.Sprintf(`{"error":"no mock found","method":"%s","path":"%s","hint":"add a route via Web UI or /_api/routes"}`, method, path)
	} else {
		status = route.Status
		routeID = route.ID
		// Parse path params into body
		body = config.ParsePath(route.Path, path, route.Body)

		if route.Delay > 0 {
			time.Sleep(time.Duration(route.Delay) * time.Millisecond)
		}

		for k, v := range route.Headers {
			w.Header().Set(k, v)
		}
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)
	fmt.Fprint(w, body)

	// Log the request
	elapsed := time.Since(start)
	reqBody := ""
	if r.Body != nil {
		b, _ := io.ReadAll(io.LimitReader(r.Body, 1024))
		reqBody = string(b)
	}
	s.cfg.AddLog(config.RequestLog{
		ID:        generateID(),
		Method:    method,
		Path:      path,
		Status:    status,
		Delay:     int(elapsed.Milliseconds()),
		RouteID:   routeID,
		Timestamp: config.Now(),
		Body:      reqBody,
	})
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
	fmt.Printf(`
🦞 MockAPI Server running!
   Web UI:    http://localhost:%d/
   Mock base: http://localhost:%d/mock/
   API:       http://localhost:%d/_api/routes
   Config:    %s
`, s.port, s.port, s.port, s.configFile)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
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
		Description: "Login",
		Headers: map[string]string{"X-Auth-Token": "mock-token-abc123"},
		Body: `{"token": "mock-token-abc123", "user": {"id": 1, "name": "Alice"}}`,
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
}
