package ws

import (
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
	m.handlers[h.Path] = &h
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

// HandleWS upgrades HTTP connection to WebSocket and handles messages
func (m *MockWS) HandleWS(w http.ResponseWriter, r *http.Request, scriptEngine *script.Engine) {
	path := strings.TrimPrefix(r.URL.Path, "/ws")
	
	handler := m.GetHandler(path)
	if handler == nil {
		http.Error(w, "No WebSocket handler for this path", 404)
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
