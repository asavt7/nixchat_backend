package repos

import "github.com/jmoiron/sqlx"

const (
	usersTable = "nix.users"
)

func NewRepositoriesPg(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepo: NewUserRepoPg(db),
		ChatRepo: NewChatPgRepo(db),
	}
}
