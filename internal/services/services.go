//go:generate mockgen -destination=./mocks/mocks.go -source=./services.go

package services

import (
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	CreateUser(user model.User, password string) (model.User, error)
	FindByUsernameOrEmail(username, email string) ([]model.User, error)
	GetByUsername(username string) (model.User, error)
	GetByID(userID uuid.UUID) (model.User, error)
	GetAll(model.PagedQuery) ([]model.User, error)
	Update(userID uuid.UUID, updateInput model.UpdateUserInfo) (model.User, error)
}

// AuthorizationService interface contains methods for working with tokens
type AuthorizationService interface {
	CheckUserCredentials(username string, password string) (model.User, error)
	GenerateTokens(userID uuid.UUID) (accessToken, refreshToken string, accessExp, refreshExp time.Time, err error)
	ParseAccessTokenToClaims(token string) (*model.Claims, error)
	ParseRefreshTokenToClaims(token string) (*model.Claims, error)
	IsNeedToRefresh(claims *model.Claims) bool
	Logout(accessTokenClaims *model.Claims) error
	ValidateAccessToken(accessTokenClaims *model.Claims) error
	ValidateRefreshToken(accessTokenClaims *model.Claims) error
	GetAccessSigningKey() string
}

type Services struct {
	UserService
	AuthorizationService
}
