package model

import "github.com/google/uuid"

type Chat struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
}

type UserToUserChat struct {
	ID      uuid.UUID   `json:"id"`
	Type    string      `json:"type"`
	UserIds []uuid.UUID `json:"userIds"`
}
