package model

import (
	"github.com/google/uuid"
	"time"
)

// Types of messages
const (
	ChatMessageRqType = "ChatMessageRq"
	ChatMessageRsType = "ChatMessageRs"

	ChatMessageNotificationType = "ChatMessageNotification"

	LoadChatsRqType = "LoadChatsRq"
	LoadChatsRsType = "LoadChatsRs"

	LoadPrivateChatsRqType = "LoadPrivateChatsRq"
	LoadPrivateChatsRsType = "LoadPrivateChatsRs"

	LoadUsersRqType = "LoadUsersRq"
	LoadUsersRsType = "LoadUsersRs"

	LoadChatMessagesRqType = "LoadChatMessagesRq"
	LoadChatMessagesRsType = "LoadChatMessagesRs"

	NewPrivateChatRqType = "NewPrivateChatRq"
	NewPrivateChatRsType = "NewPrivateChatRs"

	NewGroupChatRqType = "NewGroupChatRq"
	NewGroupChatRsType = "NewGroupChatRs"

	ErrorMessageType = "Error"
)

// ChatMessage rq and rs
type ChatMessage struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	ChatID    uuid.UUID `json:"chatId"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type LoadChatsRq struct {
	UserID uuid.UUID `json:"userId"`
}
type LoadChatsRs struct {
	UserID uuid.UUID `json:"userId"`
	Chats  []Chat    `json:"chats"`
}

type LoadPrivateChatsRq struct {
	UserID uuid.UUID `json:"userId"`
}
type LoadPrivateChatsRs struct {
	UserID uuid.UUID        `json:"userId"`
	Chats  []UserToUserChat `json:"chats"`
}

type LoadUsersRq struct {
	//todo add request params, pagination, etc
}
type LoadUsersRs struct {
	Users []User `json:"users"` //todo add userRsModel
}

type LoadChatMessagesRq struct {
	ChatID          uuid.UUID `json:"chatId"`
	BeforeTimestamp time.Time `json:"before"`

	Size   int `json:"size"`
	Offset int `json:"offset"`
}

type LoadChatMessagesRs struct {
	ChatID   uuid.UUID     `json:"chatId"`
	Messages []ChatMessage `json:"messages"`
}

type NewPrivateChatRq struct {
	FromUserID uuid.UUID `json:"fromUserId"`
	ToUserID   uuid.UUID `json:"toUserId"`
	Timestamp  time.Time `json:"timestamp"`
}

type NewPrivateChatRs struct {
	Chat UserToUserChat `json:"chat"`
}

type NewGroupChatRq struct {
	UserID    uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type NewGroupChatRs struct {
	UserID    uuid.UUID `json:"userId"`
	ChatID    uuid.UUID `json:"chatId"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type GroupChatInvite struct {
	UserID uuid.UUID `json:"userId"`
	ChatID uuid.UUID `json:"chatId"`
}

type JoinGroupChat struct {
	UserID uuid.UUID `json:"userId"`
	ChatID uuid.UUID `json:"chatId"`
}

type LeaveGroupChat struct {
	UserID uuid.UUID `json:"userId"`
	ChatID uuid.UUID `json:"chatId"`
}

type BlockUserInChat struct {
	UserID uuid.UUID `json:"userId"`
	ChatID uuid.UUID `json:"chatId"`
}
