package websocket

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients  map[string]*websocket.Conn
	mu       *sync.RWMutex
	upgrader *websocket.Upgrader
}

// NewHub initializes a new WebSocket client manager.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
		mu:      &sync.RWMutex{},
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *Hub) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return h.upgrader.Upgrade(w, r, responseHeader)
}

// Add adds a WebSocket connection to the manager.
func (h *Hub) Add(conn *websocket.Conn, username string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[username] = conn
}

// Remove removes a WebSocket connection from the manager.
func (h *Hub) Remove(username string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, username)
}

// Send sends a message to the provided clients.
func (h *Hub) Send(message []byte, usernames []string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, username := range usernames {
		slog.Info("sending message", slog.String("recipient", username))
		conn, exists := h.clients[username]
		if !exists {
			// TODO: handle user offline
			continue
		}

		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			h.Remove(username)
			conn.Close()
			return fmt.Errorf("error writing message: %w", err)
		}
	}
	return nil
}

func IsNormalCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseNormalClosure)
}
