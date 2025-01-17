package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients  map[string]*websocket.Conn
	mu       *sync.RWMutex
	upgrader *websocket.Upgrader
}

// NewHub initializes a new WebSocket client hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
		mu:      &sync.RWMutex{},
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		},
	}
}

func (h *Hub) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return h.upgrader.Upgrade(w, r, responseHeader)
}

// Add adds a WebSocket connection to the hub.
func (h *Hub) Add(conn *websocket.Conn, username string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[username] = conn
}

// Remove removes a WebSocket connection from the hub.
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
		conn, exists := h.clients[username]
		if !exists {
			// No queuing mechanism is used; if the recipient is offline they will not receive the message.
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

// Close closes a WebSocket connection.
func (h *Hub) Close(username string) error {
	return closeConnection(h.clients[username])
}

func IsNormalCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseNormalClosure)
}
