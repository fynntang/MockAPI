package ws

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"mockapi/pkg/script"
)

type MockWS struct {
	mu       sync.Mutex
	handlers map[string]*WSHandler
	upgrader websocket.Upgrader
}

type WSHandler struct {
	Path        string            `json:"path"`
	Description string            `json:"description"`
	OnConnect   string            `json:"on_connect,omitempty"`   // JS to run on connect
	OnMessage   string            `json:"on_message,omitempty"`  // JS to run on message
	AutoReply   string            `json:"auto_reply,omitempty"`  // Static JSON reply
	Delay       int               `json:"delay_ms,omitempty"`
	
	// Stream mode configuration
	StreamEnabled  bool     `json:"stream_enabled,omitempty"`
	StreamMessages []string `json:"stream_messages,omitempty"`
	StreamInterval int      `json:"stream_interval_ms,omitempty"` // Fixed interval in ms
	StreamRandom   bool     `json:"stream_random,omitempty"`      // Use random interval
	StreamMinDelay int      `json:"stream_min_delay_ms,omitempty"`
	StreamMaxDelay int      `json:"stream_max_delay_ms,omitempty"`
	StreamLoop     bool     `json:"stream_loop,omitempty"` // Loop through messages
	StreamFormat   string   `json:"stream_format,omitempty"` // "json" or "text"
}

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

func New() *MockWS {
	return &MockWS{
		handlers: make(map[string]*WSHandler),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (m *MockWS) AddHandler(h WSHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Normalize path: ensure it starts with /
	path := h.Path
	if !strings.HasPrefix(path, "/") && path != "" {
		path = "/" + path
	}
	h.Path = path
	m.handlers[path] = &h
}

func (m *MockWS) GetHandler(path string) *WSHandler {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.handlers[path]
}

func (m *MockWS) ListHandlers() []*WSHandler {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]*WSHandler, 0, len(m.handlers))
	for _, h := range m.handlers {
		result = append(result, h)
	}
	return result
}

func (m *MockWS) DeleteHandler(path string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Normalize path: ensure it starts with /
	normalizedPath := path
	if !strings.HasPrefix(normalizedPath, "/") {
		normalizedPath = "/" + normalizedPath
	}
	
	if _, exists := m.handlers[normalizedPath]; exists {
		delete(m.handlers, normalizedPath)
		return true
	}
	
	// Also try the original path (in case user stored it differently)
	if _, exists := m.handlers[path]; exists {
		delete(m.handlers, path)
		return true
	}
	
	return false
}

// HandleWS upgrades HTTP connection to WebSocket and handles messages
func (m *MockWS) HandleWS(w http.ResponseWriter, r *http.Request, scriptEngine *script.Engine) {
	path := strings.TrimPrefix(r.URL.Path, "/ws")
	
	// Normalize path: ensure it starts with /
	if !strings.HasPrefix(path, "/") && path != "" {
		path = "/" + path
	}
	
	handler := m.GetHandler(path)
	if handler == nil {
		// Try to find handler without leading slash (compatibility)
		if strings.HasPrefix(path, "/") {
			handler = m.GetHandler(path[1:])
		}
	}
	
	if handler == nil {
		http.Error(w, "No WebSocket handler for this path: "+path, 404)
		return
	}

	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Send connect message if defined
	if handler.OnConnect != "" {
		conn.WriteMessage(websocket.TextMessage, []byte(handler.OnConnect))
	}

	// Stream mode: push messages at intervals
	if handler.StreamEnabled && len(handler.StreamMessages) > 0 {
		m.handleStreamMode(conn, handler)
		return
	}

	// Standard mode: respond to messages
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var reply []byte

		// Try script first
		if handler.OnMessage != "" && scriptEngine != nil {
			ctx := script.ScriptContext{
				Method: "WS",
				Path:   path,
				Body:   string(message),
			}
			_, body, _, _ := scriptEngine.Execute(handler.OnMessage, ctx)
			if body != "" {
				reply = []byte(body)
			}
		}

		// Fallback to auto-reply
		if reply == nil && handler.AutoReply != "" {
			reply = []byte(handler.AutoReply)
		}

		// Default echo
		if reply == nil {
			reply = message
		}

		if handler.Delay > 0 {
			time.Sleep(time.Duration(handler.Delay) * time.Millisecond)
		}

		conn.WriteMessage(msgType, reply)
	}
}

// handleStreamMode handles streaming messages to the client
func (m *MockWS) handleStreamMode(conn *websocket.Conn, handler *WSHandler) {
	msgIndex := 0
	msgCount := len(handler.StreamMessages)
	
	for {
		// Get current message
		msg := handler.StreamMessages[msgIndex]
		
		// Format message based on StreamFormat
		var formattedMsg string
		if handler.StreamFormat == "json" {
			// If message is not already JSON object, wrap it
			trimmed := strings.TrimSpace(msg)
			if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
				formattedMsg = fmt.Sprintf(`{"data":%q,"index":%d,"timestamp":%d}`, msg, msgIndex, time.Now().Unix())
			} else {
				formattedMsg = msg
			}
		} else {
			formattedMsg = msg
		}
		
		// Send message
		if err := conn.WriteMessage(websocket.TextMessage, []byte(formattedMsg)); err != nil {
			return // Client disconnected
		}
		
		// Calculate delay
		var delay time.Duration
		if handler.StreamRandom {
			minDelay := handler.StreamMinDelay
			maxDelay := handler.StreamMaxDelay
			if minDelay <= 0 {
				minDelay = 100
			}
			if maxDelay <= minDelay {
				maxDelay = minDelay + 1000
			}
			delay = time.Duration(minDelay+rand.Intn(maxDelay-minDelay)) * time.Millisecond
		} else {
			delay = time.Duration(handler.StreamInterval) * time.Millisecond
			if delay <= 0 {
				delay = 1000 * time.Millisecond // Default 1 second
			}
		}
		
		// Wait for next message or client disconnect
		select {
		case <-time.After(delay):
		case <-waitForClose(conn):
			return
		}
		
		// Move to next message
		msgIndex++
		if msgIndex >= msgCount {
			if handler.StreamLoop {
				msgIndex = 0
			} else {
				// Non-looping: stop after all messages sent
				return
			}
		}
	}
}

// waitForClose returns a channel that signals when the connection is closed
func waitForClose(conn *websocket.Conn) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
	return done
}
