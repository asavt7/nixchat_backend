package main

import (
	"context"
	"github.com/asavt7/nixchat_backend/internal/config"
	"github.com/asavt7/nixchat_backend/internal/db"
	"github.com/asavt7/nixchat_backend/internal/handlers"
	"github.com/asavt7/nixchat_backend/internal/repos"
	"github.com/asavt7/nixchat_backend/internal/server"
	"github.com/asavt7/nixchat_backend/internal/services"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	pgdb, err := db.NewPostgreDb(cfg.Postgres)
	if err != nil {
		log.Fatal("ERROR init postgres db", err)
	}

	repositorories := repos.NewRepositoriesPg(pgdb)

	serves := services.NewServices(repositorories)

	handler := handlers.NewAPIHandler(serves)

	apiServer := server.NewAPIServer(cfg, handler.InitRoutes(cfg))

	go func() {
		log.Info("Server started")
		err := apiServer.Run()
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Info("Server stopping")

	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	// todo https://github.com/asavt7/nixchat_backend/issues/11 switch off readiness probe to prevent new requests
	if err := apiServer.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: %v", err)
	}
	if err := pgdb.Close(); err != nil {
		log.Errorf("failed to stop db: %v", err)
	}
	log.Info("Server stopped")

}
