package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type FrontMessageRq struct {
	UserID          uuid.UUID   `json:"userId"`
	MessageType     string      `json:"type"`
	UserGeneratedID string      `json:"ugId"`
	Timestamp       time.Time   `json:"timestamp"`
	Payload         interface{} `json:"payload"`
}

func (d *FrontMessageRq) UnmarshalJSON(data []byte) error {
	// Parse only the "type" field first.
	var meta struct {
		Payload         json.RawMessage `json:"payload"`
		UserID          uuid.UUID       `json:"userId"`
		MessageType     string          `json:"type"`
		UserGeneratedID string          `json:"ugId"`
		Timestamp       time.Time       `json:"timestamp"`
	}

	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}
	// Determine which struct to unmarshal into according to "type".
	switch meta.MessageType {
	case ChatMessageRqType:
		d.Payload = &ChatMessage{}
	case LoadChatsRqType:
		d.Payload = &LoadChatsRq{}
	case LoadPrivateChatsRqType:
		d.Payload = &LoadPrivateChatsRq{}
	case LoadUsersRqType:
		d.Payload = &LoadUsersRq{}
	case LoadChatMessagesRqType:
		d.Payload = &LoadChatMessagesRq{}
	case NewPrivateChatRqType:
		d.Payload = &NewPrivateChatRq{}
	case NewGroupChatRqType:
		d.Payload = &NewGroupChatRq{}
	default:
		return fmt.Errorf("%q is an invalid item type", meta.MessageType)
	}

	d.UserID = meta.UserID
	d.MessageType = meta.MessageType
	d.Timestamp = meta.Timestamp
	d.UserGeneratedID = meta.UserGeneratedID

	return json.Unmarshal(meta.Payload, d.Payload)
}

type FrontMessageRs struct {
	Ok              bool        `json:"ok"`
	MessageType     string      `json:"type"`
	UserGeneratedID string      `json:"ugId"`
	Payload         interface{} `json:"payload"`
	ErrorMessage    string      `json:"error"`
}
