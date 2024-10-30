package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/entity"
)

type Msg struct {
	Type    MsgType    `json:"type"`
	Payload MsgPayload `json:"payload"`
}

type MsgType int

const (
	MsgTypeSendChatMessage MsgType = iota
)

type MsgPayload interface {
	isWebSocketMsgPayload()
}

type PayloadSendChatMessage struct {
	ConversationName string         `json:"chat_name"`
	Recipients       []string       `json:"recipient_usernames"`
	Message          entity.Message `json:"message"`
}

func (PayloadSendChatMessage) isWebSocketMsgPayload() {}

func (m *Msg) UnmarshalJSON(data []byte) error {
	var temp struct {
		Type    MsgType         `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	m.Type = temp.Type

	switch temp.Type {
	case MsgTypeSendChatMessage:
		var payload PayloadSendChatMessage
		if err := json.Unmarshal(temp.Payload, &payload); err != nil {
			return err
		}
		m.Payload = payload
	default:
		return fmt.Errorf("unknown MsgType %v", temp.Type)
	}

	return nil
}
