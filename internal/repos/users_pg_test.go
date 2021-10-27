package repos_test

import (
	"errors"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/asavt7/nixchat_backend/internal/repos"
	"github.com/google/uuid"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func initUserStorage(t *testing.T) (repos.UserRepo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo := repos.NewUserRepoPg(db)
	return repo, mock, func() {
		_ = db.Close()
	}
}

func TestUserRepoPg_Create(t *testing.T) {
	repo, mock, destroyFunc := initUserStorage(t)
	defer destroyFunc()

	userID, _ := uuid.NewUUID()

	tests := []struct {
		name    string
		s       repos.UserRepo
		user    model.User
		mock    func()
		want    model.User
		wantErr bool
	}{
		{
			name: "OK",
			s:    repo,
			user: model.User{
				Email:        "email@googlle.com",
				Username:     "username",
				PasswordHash: "password",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
				mock.ExpectQuery("INSERT INTO nix.users \\(username,email,password_hash\\) values").WithArgs("username", "email@googlle.com", "password").WillReturnRows(rows)
			},
			want: model.User{
				ID:           userID,
				Username:     "username",
				Email:        "email@googlle.com",
				PasswordHash: "password",
			},
		},
		{
			name: "empty fields",
			s:    repo,
			user: model.User{},
			mock: func() {
				mock.ExpectQuery("INSERT INTO nix.users \\(username,email,password_hash\\) values").WithArgs("", "", "").WillReturnError(errors.New("invalid values"))
			},
			want:    model.User{},
			wantErr: true,
		},
		{
			name: "user already exists",
			s:    repo,
			user: model.User{
				Email:        "email@googlle.com",
				Username:     "username",
				PasswordHash: "password",
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO nix.users \\(username,email,password_hash\\) values").WithArgs("username", "email@googlle.com", "password").WillReturnError(errors.New("user already exists"))
			},
			want:    model.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Create(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf(" error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf(" = %v, want %v", got, tt.want)
			}
		})
	}

}

/*
func TestUserRepoPg_FindByUsernameOrEmail(t *testing.T) {
	repo, mock, destroyFunc := initUserStorage(t)
	defer destroyFunc()
	userID ,_:= uuid.NewUUID()

	tests := []struct {
		name    string
		s       repos.UserRepo
		user    model.User
		mock    func()
		wantErr bool
	}{
		{
			name: "OK",
			s:    repo,
			user: model.User{
				ID:           userID,
				Email:        "email@googlle.com",
				Username:     "username",
				PasswordHash: "password",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash"}).AddRow(userID, "username", "email@googlle.com", "password")
				mock.ExpectQuery("SELECT \\* FROM nix.users WHERE username").WithArgs("username","email@googlle.com").WillReturnRows(rows)
			},
		},
		{
			name: "no user found",
			s:    repo,
			user: model.User{
				ID:           userID,
				Email:        "email@googlle.com",
				Username:     "username",
				PasswordHash: "password",
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash"})
				mock.ExpectQuery("SELECT \\* FROM nix.users WHERE username").WithArgs("username","email@googlle.com").WillReturnRows(rows)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.FindByUsernameOrEmail(tt.user.Username,tt.user.Email)
			if ((err != nil) != tt.wantErr)  {
				t.Errorf(" error new = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !assert.Equal(t, []model.User{tt.user},got)  {
				t.Errorf(" = %v, want %v", got, tt.user)
			}
		})
	}

}

*/
