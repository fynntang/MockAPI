package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"mockapi/pkg/config"
)

type Server struct {
	cfg  *config.Config
	mux  *http.ServeMux
	port int
}

func New(cfg *config.Config) *Server {
	s := &Server{
		cfg:  cfg,
		mux:  http.NewServeMux(),
		port: cfg.Port,
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

	// Mock routes - catch all
	s.mux.HandleFunc("/mock/", s.handleMock)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
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

func (s *Server) handleAPIRoutes(w http.ResponseWriter, r *http.Request) {
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
		s.cfg.Save("mockapi.json")
		json.NewEncoder(w).Encode(route)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		s.cfg.DeleteRoute(id)
		s.cfg.Save("mockapi.json")
		w.WriteHeader(204)
	}
}

func (s *Server) handleMock(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	path := strings.TrimPrefix(r.URL.Path, "/mock")

	route := s.cfg.FindRoute(method, path)
	if route == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{
			"error":      "no mock found",
			"method":     method,
			"path":       path,
			"hint":       "add a route via Web UI or /_api/routes",
		})
		return
	}

	if route.Delay > 0 {
		time.Sleep(time.Duration(route.Delay) * time.Millisecond)
	}

	for k, v := range route.Headers {
		w.Header().Set(k, v)
	}
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(route.Status)
	fmt.Fprint(w, route.Body)
}

func (s *Server) Start() error {
	fmt.Printf("🚀 MockAPI Server running at http://localhost:%d\n", s.port)
	fmt.Printf("   Web UI:  http://localhost:%d/\n", s.port)
	fmt.Printf("   Mock base: http://localhost:%d/mock/\n", s.port)
	fmt.Printf("   Routes API: http://localhost:%d/_api/routes\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
