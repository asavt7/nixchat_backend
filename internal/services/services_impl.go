package services

import "github.com/asavt7/nixchat_backend/internal/repos"

func NewServices(repos *repos.Repositories) *Services {
	return &Services{UserService: NewUserServiceImpl(repos.UserRepo)}
}
