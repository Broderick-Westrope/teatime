package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ClientManager is the interface that wraps the helper methods for managing WebSocket clients.
type ClientManager interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error)
	Add(conn Conn)
	Remove(conn Conn)
	Broadcast(message []byte) error
}

type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type ClientManagerImpl struct {
	clients  map[Conn]struct{}
	mu       *sync.RWMutex
	upgrader *websocket.Upgrader
}

// NewClientManager initializes a new WebSocket client manager.
func NewClientManager() *ClientManagerImpl {
	return &ClientManagerImpl{
		clients: make(map[Conn]struct{}),
		mu:      &sync.RWMutex{},
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (m *ClientManagerImpl) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error) {
	return m.upgrader.Upgrade(w, r, responseHeader)
}

// Add adds a WebSocket connection to the manager.
func (m *ClientManagerImpl) Add(conn Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[conn] = struct{}{}
}

// Remove removes a WebSocket connection from the manager.
func (m *ClientManagerImpl) Remove(conn Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, conn)
}

// Broadcast sends a message to all connected clients.
func (m *ClientManagerImpl) Broadcast(message []byte) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for conn := range m.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			m.Remove(conn)
			conn.Close()
			return fmt.Errorf("error writing message: %w", err)
		}
	}
	return nil
}
