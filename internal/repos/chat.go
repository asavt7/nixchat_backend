package repos

import (
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/google/uuid"
)

type ChatRepo interface {
	SaveMessage(mes model.ChatMessage) (model.ChatMessage, error)
	GetChatMembers(chatID uuid.UUID) ([]uuid.UUID, error)
	CreatePrivateChat(ownerID uuid.UUID, chat model.Chat, membersIDs uuid.UUID) (model.Chat, error)
	GetChatMessages(chatID uuid.UUID) ([]model.ChatMessage, error)
	GetChatMessagesByQuery(chatID uuid.UUID, query model.Query) ([]model.ChatMessage, error)

	GetUserChats(userID uuid.UUID) ([]model.Chat, error)
	GetUserPrivateChats(userID uuid.UUID) ([]model.UserToUserChat, error)
}
