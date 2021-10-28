package server

import (
	"context"
	"github.com/asavt7/nixchat_backend/internal/config"
	"net/http"
)

// APIServer struct
type APIServer struct {
	httpServer *http.Server
}

// NewAPIServer constructs APIServer
func NewAPIServer(cfg *config.Config, handler http.Handler) *APIServer {
	s := &APIServer{
		httpServer: &http.Server{
			Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
	}
	return s
}

// Run method run server or fail app
func (s *APIServer) Run() error {
	return s.httpServer.ListenAndServe()
}

// Stop http server
func (s *APIServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
