package services

import (
	"errors"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/asavt7/nixchat_backend/internal/repos"
	"github.com/asavt7/nixchat_backend/internal/utils"
)

type UserServiceImpl struct {
	repo repos.UserRepo
}

func NewUserServiceImpl(repo repos.UserRepo) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (u *UserServiceImpl) CreateUser(user model.User, password string) (model.User, error) {
	existingUsers, err := u.repo.FindByUsernameOrEmail(user.Username, user.Email)
	if err != nil {
		return user, err
	}
	if userWithSameEmailAlreadyExists(existingUsers, user.Email) {
		return user, errors.New("ERROR : User with same email already exists")
	}
	if userWithSameUsernameAlreadyExists(existingUsers, user.Username) {
		return user, errors.New("ERROR : User with same username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return model.User{}, err
	}
	user.PasswordHash = hashedPassword
	return u.repo.Create(user)
}

func userWithSameUsernameAlreadyExists(users []model.User, username string) bool {
	for _, user := range users {
		if user.Username == username {
			return true
		}
	}
	return false
}

func userWithSameEmailAlreadyExists(users []model.User, email string) bool {
	for _, user := range users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (u *UserServiceImpl) FindByUsernameOrEmail(username, email string) ([]model.User, error) {
	return u.repo.FindByUsernameOrEmail(username, email)
}
