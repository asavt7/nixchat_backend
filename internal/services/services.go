//go:generate mockgen -destination=./mocks/mocks.go -source=./services.go

package services

import "github.com/asavt7/nixchat_backend/internal/model"

type UserService interface {
	CreateUser(user model.User, password string) (model.User, error)
	FindByUsernameOrEmail(username, email string) ([]model.User, error)

	/*
			GetUserByID(id int) (model.User, error)
		GetUserByEmail(email string) (model.User, error)
		GetUserByUsername(username string) (model.User, error)
	*/
}

type Services struct {
	UserService
}
