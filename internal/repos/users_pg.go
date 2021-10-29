package repos

import (
	"database/sql"
	"errors"
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

func (u *UserRepoPg) GetByID(userID uuid.UUID) (model.User, error) {
	var user model.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)
	err := u.db.Get(&user, query, userID)
	if err == sql.ErrNoRows {
		return user, errors.New("user not found")
	}
	if err != nil {
		log.Error(err.Error())
	}
	return user, err

}

func (u *UserRepoPg) GetAll(pagedQuery model.PagedQuery) ([]model.User, error) {
	users := make([]model.User, 0, pagedQuery.Size)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY username LIMIT $1 OFFSET $2", usersTable)
	err := u.db.Select(&users, query, pagedQuery.Size, pagedQuery.Offset)
	if err != nil {
		log.Error(err.Error())
	}
	return users, err
}

func (u *UserRepoPg) Update(userID uuid.UUID, input model.UpdateUserInfo) (model.User, error) {
	var result model.User
	if input.AvatarURL == nil {
		return result, errors.New("empty fields to update")
	}

	argNum := 0
	updateArgs := make([]string, 0)
	updateVals := make([]interface{}, 0)

	if input.AvatarURL != nil {
		updateArgs = append(updateArgs, "avatar_url")
		updateVals = append(updateVals, *input.AvatarURL)
		argNum++
	}

	updateVals = append(updateVals, userID)

	setExpression := convertToSetStrs(updateArgs)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING *", usersTable, setExpression, argNum+1)

	err := u.db.Get(&result, query, updateVals...)
	if err != nil {
		log.Error(err.Error())
	}
	return result, err

}

func (u *UserRepoPg) GetByUsername(username string) (model.User, error) {
	var user model.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1", usersTable)
	err := u.db.Get(&user, query, username)
	if err == sql.ErrNoRows {
		return user, errors.New("user not found")
	}
	if err != nil {
		log.Error(err.Error())
	}
	return user, err
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
