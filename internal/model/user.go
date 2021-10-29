package model

import "github.com/google/uuid"

// User model
type User struct {
	ID uuid.UUID `json:"id"`

	Username     string `json:"username" binding:"required" validate:"required,min=2,max=255"`
	Email        string `json:"email" binding:"required" validate:"required,email,max=255"`
	PasswordHash string `json:"-" db:"password_hash"`

	AvatarURL string `json:"avatar_url" db:"avatar_url"`
}

type UpdateUserInfo struct {
	AvatarURL *string `json:"avatar_url" db:"avatar_url"`
}
