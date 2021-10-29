//go:generate mockgen -destination=./mocks/mocks.go -source=./repositories.go

package repos

import (
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/google/uuid"
)

type UserRepo interface {
	Create(user model.User) (model.User, error)
	FindByUsernameOrEmail(username, email string) ([]model.User, error)
	GetByUsername(username string) (model.User, error)
	GetByID(userID uuid.UUID) (model.User, error)
	GetAll(pagedQuery model.PagedQuery) ([]model.User, error)
	Update(userID uuid.UUID, input model.UpdateUserInfo) (model.User, error)
}

type Repositories struct {
	UserRepo
}
