package repos

import (
	"fmt"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type UserRepoPg struct {
	db *sqlx.DB
}

func NewUserRepoPg(db *sqlx.DB) *UserRepoPg {
	return &UserRepoPg{db: db}
}

func (u *UserRepoPg) Create(user model.User) (model.User, error) {
	var id uuid.UUID

	query := fmt.Sprintf("INSERT INTO %s (username,email,password_hash) values( $1,$2, $3) RETURNING id", usersTable)
	err := u.db.Get(&id, query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		log.Error(err.Error())
		return user, fmt.Errorf("cannot create user username=%s email=%s", user.Username, user.Email)
	}
	user.ID = id
	return user, nil
}

func (u *UserRepoPg) FindByUsernameOrEmail(username, email string) ([]model.User, error) {
	var users []model.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 or email=$2", usersTable)
	err := u.db.Select(&users, query, username, email)
	if err != nil {
		log.Error(err.Error())
	}
	return users, err
}
