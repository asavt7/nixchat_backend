package app

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

type ChatApp struct {
	cfg *config.Config
}

func NewChatApp(cfg *config.Config) *ChatApp {
	return &ChatApp{cfg: cfg}
}

func (a *ChatApp) Run() {

	pgdb, err := db.NewPostgreDb(a.cfg.Postgres)
	if err != nil {
		log.Fatal("ERROR init postgres db", err)
	}
	redisClient := repos.InitRedisClient(repos.InitRedisOpts(&a.cfg.Redis))

	repositories := repos.NewRepositoriesPg(pgdb)

	tokenKeeper := repos.NewRedisTokenStorage(redisClient, &a.cfg.Auth)

	serves := services.NewServices(a.cfg, repositories, tokenKeeper)

	handler := handlers.NewAPIHandler(serves)

	apiServer := server.NewAPIServer(a.cfg, handler.InitRoutes(a.cfg))

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
