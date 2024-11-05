package websocket

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/gorilla/websocket"
)

// Client is a struct that represents the websocket client.
type Client struct {
	conn      *websocket.Conn
	mu        *sync.RWMutex
	uri       string
	sessionID string
}

// NewClient is a function used to create a new websocket client.
func NewClient(uri, sessionID string) (*Client, error) {
	uri = strings.Replace(uri, "http", "ws", 1)

	c := &Client{
		conn:      nil,
		mu:        &sync.RWMutex{},
		uri:       uri,
		sessionID: sessionID,
	}
	err := c.connect()
	return c, err
}

// connect will create and store a connection to the WebSocket server.
func (c *Client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	header := http.Header{}
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: c.sessionID,
	}
	header.Add("Cookie", cookie.String())

	conn, _, err := websocket.DefaultDialer.Dial(c.uri, header)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// TODO: make use of this
func (c *Client) Reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	baseDelay := time.Second
	maxDelay := 10 * time.Second

	for attempt := 1; attempt <= 3; attempt++ {
		c.conn, _, err = websocket.DefaultDialer.Dial(c.uri, nil)
		if err == nil {
			return nil
		}

		// Calculate exponential backoff with jitter
		delay := baseDelay * (1 << (attempt - 1))
		if delay > maxDelay {
			delay = maxDelay
		}
		delay = time.Duration(rand.Int63n(int64(delay)))

		time.Sleep(delay)
	}
	return fmt.Errorf("failed to reconnect after 10 attempts: %w", err)
}

// Close will gracefully close the WebSocket connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return closeConnection(c.conn)
}

func (c *Client) ReadMessage() (*Msg, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, msgData, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var msg Msg
	err = json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, err
}

func (c *Client) SendChatMessage(message entity.Message, conversationMD entity.ConversationMetadata, recipients []string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.conn.WriteJSON(Msg{
		Type: MsgTypeSendChatMessage,
		Payload: PayloadSendChatMessage{
			ConversationMD: conversationMD,
			Message:        message,
			Recipients:     recipients,
		},
	})
}
