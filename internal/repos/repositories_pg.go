package repos

import "github.com/jmoiron/sqlx"
import _ "github.com/lib/pq" // import postgres driver

const (
	usersTable = "nix.users"
)

func NewRepositoriesPg(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepo: NewUserRepoPg(db),
	}
}
