package chats

import (
	"sync"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// StreamHub manages websocket subscribers per conversation.
type StreamHub struct {
	mu          sync.RWMutex
	subscribers map[uuid.UUID]map[*websocket.Conn]struct{}
}

// NewStreamHub constructs a new hub.
func NewStreamHub() *StreamHub {
	return &StreamHub{
		subscribers: make(map[uuid.UUID]map[*websocket.Conn]struct{}),
	}
}

// Register registers a websocket connection for conversation updates.
func (h *StreamHub) Register(conversationID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.subscribers[conversationID]; !ok {
		h.subscribers[conversationID] = make(map[*websocket.Conn]struct{})
	}
	h.subscribers[conversationID][conn] = struct{}{}
}

// Unregister removes a websocket connection from the hub.
func (h *StreamHub) Unregister(conversationID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subs, ok := h.subscribers[conversationID]; ok {
		delete(subs, conn)
		if len(subs) == 0 {
			delete(h.subscribers, conversationID)
		}
	}
}

// Broadcast sends an event to all subscribers of the conversation.
func (h *StreamHub) Broadcast(conversationID uuid.UUID, payload any) {
	h.mu.RLock()
	conns := h.copySubscribers(conversationID)
	h.mu.RUnlock()

	for conn := range conns {
		if err := conn.WriteJSON(payload); err != nil {
			_ = conn.Close()
			h.Unregister(conversationID, conn)
		}
	}
}

func (h *StreamHub) copySubscribers(conversationID uuid.UUID) map[*websocket.Conn]struct{} {
	result := make(map[*websocket.Conn]struct{})
	if subs, ok := h.subscribers[conversationID]; ok {
		for conn := range subs {
			result[conn] = struct{}{}
		}
	}
	return result
}
