package main

import (
	"github.com/asavt7/nixchat_backend/internal/app"
	"github.com/asavt7/nixchat_backend/internal/config"
	log "github.com/sirupsen/logrus"
)

// @title nixchat_backed
// @version 1.0
// @description This is a backend for project https://github.com/users/asavt7/projects/2

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	cfg, err := config.Init("configs")
	if err != nil {
		log.Fatal("ERROR init configs", err)
	}
	chatApp := app.NewChatApp(cfg)
	chatApp.Run()
}
