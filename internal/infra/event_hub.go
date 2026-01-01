package infra

import (
	"sync"
)

// EventHub manages SSE client connections and broadcasts events
type EventHub struct {
	clients map[chan string]bool
	mu      sync.RWMutex
}

// NewEventHub creates a new EventHub
func NewEventHub() *EventHub {
	return &EventHub{
		clients: make(map[chan string]bool),
	}
}

// Register adds a new client channel
func (h *EventHub) Register(client chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
}

// Unregister removes a client channel
func (h *EventHub) Unregister(client chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client)
	close(client)
}

// Broadcast sends an event to all connected clients
func (h *EventHub) Broadcast(event string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client <- event:
		default:
			// Client channel is full, skip
		}
	}
}

// ClientCount returns the number of connected clients
func (h *EventHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
