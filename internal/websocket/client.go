package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/gorilla/websocket"
)

// Client is a struct that represents the websocket client.
type Client struct {
	conn     *websocket.Conn
	endpoint string
	username string
}

// NewClient is a function used to create a new websocket client.
func NewClient(endpoint, username string) (*Client, error) {
	c := &Client{
		conn:     nil,
		endpoint: endpoint,
		username: username,
	}
	err := c.connect()
	return c, err
}

// connect will create and store a connection to the WebSocket server.
func (c *Client) connect() error {
	header := http.Header{}
	header.Add("username", c.username)

	conn, resp, err := websocket.DefaultDialer.Dial(c.endpoint, header)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.conn = conn
	return nil
}

// Close will gracefully close the WebSocket connection.
func (c *Client) Close() error {
	err := c.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Minute),
	)
	if err != nil {
		return err
	}

	err = c.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}

	for {
		_, _, err = c.conn.NextReader()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			break
		}
		if err != nil {
			break
		}
	}

	return c.conn.Close()
}

func (c *Client) ReadMessage() (*Msg, error) {
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

func (c *Client) SendChatMessage(message data.Message, recipientUsernames []string) error {
	return c.conn.WriteJSON(Msg{
		Type: MsgTypeSendChatMessage,
		Payload: PayloadSendChatMessage{
			ChatName:           recipientUsernames[0],
			RecipientUsernames: recipientUsernames,
			Message:            message,
		},
	})
}
