package ws

import (
	"testing"
)

func TestAddHandler(t *testing.T) {
	m := New()
	
	// Test adding handler with leading slash
	h1 := WSHandler{Path: "/chat", Description: "Chat handler"}
	m.AddHandler(h1)
	
	if m.GetHandler("/chat") == nil {
		t.Error("Expected to find handler at /chat")
	}
	
	// Test adding handler without leading slash (should be normalized)
	h2 := WSHandler{Path: "notifications", Description: "Notifications handler"}
	m.AddHandler(h2)
	
	if m.GetHandler("/notifications") == nil {
		t.Error("Expected to find handler at /notifications (normalized)")
	}
	
	// Original path should also work due to normalization
	if m.GetHandler("notifications") != nil {
		t.Error("Should not find handler at non-normalized path")
	}
}

func TestDeleteHandler(t *testing.T) {
	m := New()
	
	// Add a handler
	h := WSHandler{Path: "/chat", Description: "Chat handler"}
	m.AddHandler(h)
	
	// Verify it exists
	if m.GetHandler("/chat") == nil {
		t.Fatal("Handler should exist before delete")
	}
	
	// Delete with normalized path
	if !m.DeleteHandler("/chat") {
		t.Error("DeleteHandler should return true for existing handler")
	}
	
	// Verify it's deleted
	if m.GetHandler("/chat") != nil {
		t.Error("Handler should be deleted")
	}
	
	// Delete non-existent handler
	if m.DeleteHandler("/nonexistent") {
		t.Error("DeleteHandler should return false for non-existent handler")
	}
}

func TestDeleteHandlerNormalizePath(t *testing.T) {
	m := New()
	
	// Add handler with leading slash
	h := WSHandler{Path: "/chat", Description: "Chat handler"}
	m.AddHandler(h)
	
	// Delete without leading slash (should still work)
	if !m.DeleteHandler("chat") {
		t.Error("DeleteHandler should normalize path and find the handler")
	}
	
	if m.GetHandler("/chat") != nil {
		t.Error("Handler should be deleted")
	}
}

func TestListHandlers(t *testing.T) {
	m := New()
	
	// Add multiple handlers
	m.AddHandler(WSHandler{Path: "/chat", Description: "Chat"})
	m.AddHandler(WSHandler{Path: "/notifications", Description: "Notifications"})
	
	handlers := m.ListHandlers()
	if len(handlers) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(handlers))
	}
}

func TestHandleWSPathMatching(t *testing.T) {
	m := New()
	
	// Test that paths are normalized when added
	h := WSHandler{Path: "chat", Description: "Chat handler"}
	m.AddHandler(h)
	
	// Handler should be stored with normalized path
	stored := m.GetHandler("/chat")
	if stored == nil {
		t.Error("Handler should be accessible with normalized path /chat")
	}
	
	// Description should be preserved
	if stored.Description != "Chat handler" {
		t.Errorf("Expected description 'Chat handler', got '%s'", stored.Description)
	}
}