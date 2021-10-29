package services

import (
	"github.com/asavt7/nixchat_backend/internal/config"
	"github.com/asavt7/nixchat_backend/internal/repos"
	"github.com/asavt7/nixchat_backend/internal/services/auth"
)

func NewServices(cfg *config.Config, repos *repos.Repositories, keeper repos.TokenKeeper) *Services {

	return &Services{UserService: NewUserServiceImpl(repos.UserRepo),
		AuthorizationService: auth.NewAuthorizationServiceImpl(&cfg.Auth, repos.UserRepo, keeper)}
}
