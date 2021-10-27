//go:generate mockgen -destination=./mocks/mocks.go -source=./repositories.go

package repos

import "github.com/asavt7/nixchat_backend/internal/model"

type UserRepo interface {
	Create(user model.User) (model.User, error)
	FindByUsernameOrEmail(username, email string) ([]model.User, error)

	/*	GetByID(userID int) (model.User, error)
		FindByUsername(username string) (model.User, error)
		FindByEmail(email string) (model.User, error)
	*/
}

type Repositories struct {
	UserRepo
}
