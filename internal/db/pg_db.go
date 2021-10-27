package db

import (
	"fmt"
	"github.com/asavt7/nixchat_backend/internal/config"
	"github.com/jmoiron/sqlx"
)

// NewPostgreDb create  *sqlx.DB instance and ping connection. If failed - fail app
func NewPostgreDb(cfg config.PostgresConfig) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}
