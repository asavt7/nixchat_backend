package repos

import (
	"github.com/asavt7/nixchat_backend/internal/model"
)

// TokenKeeper interface provides methods for storing model.CachedTokens in cache
type TokenKeeper interface {
	Get(userID string) (model.CachedTokens, error)
	Delete(userID string) error
	Save(userID string, tokens model.CachedTokens) (model.CachedTokens, error)
}
