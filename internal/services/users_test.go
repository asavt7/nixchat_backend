package services_test

import (
	"errors"
	"fmt"
	"github.com/asavt7/nixchat_backend/internal/model"
	mock_repos "github.com/asavt7/nixchat_backend/internal/repos/mocks"
	"github.com/asavt7/nixchat_backend/internal/services"
	"github.com/asavt7/nixchat_backend/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserMatcher struct {
	User     model.User
	Password string
}

func (u UserMatcher) Matches(x interface{}) bool {
	switch x.(type) {
	case model.User:
		return x.(model.User).Username == u.User.Username && x.(model.User).Email == u.User.Email && utils.CheckPassword(x.(model.User).PasswordHash, u.Password) == nil
	default:
		return false
	}
}

func (u UserMatcher) String() string {
	return fmt.Sprintf("username=%v password=%s", u.User, u.Password)
}

func TestUserServiceImpl_CreateUser(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	userRepo := mock_repos.NewMockUserRepo(controller)
	userService := services.NewUserServiceImpl(userRepo)

	userID, _ := uuid.NewUUID()

	t.Run("ok", func(t *testing.T) {
		expectedCreatedUser := model.User{Email: "email", Username: "username", ID: userID}
		userRepo.EXPECT().Create(UserMatcher{
			User:     expectedCreatedUser,
			Password: "password",
		}).Times(1).Return(expectedCreatedUser, nil)
		userRepo.EXPECT().FindByUsernameOrEmail(expectedCreatedUser.Username, expectedCreatedUser.Email).Times(1).Return([]model.User{}, nil)

		user, err := userService.CreateUser(model.User{Email: "email", Username: "username"}, "password")
		if err != nil {
			t.Errorf("err should be nil")
		}

		assert.Equal(t, expectedCreatedUser, user)
	})

	t.Run("username already exists", func(t *testing.T) {
		alreadyCreatedUser := model.User{Email: "email", Username: "username", ID: userID}
		userRepo.EXPECT().Create(UserMatcher{
			User:     alreadyCreatedUser,
			Password: "password",
		}).Times(0).Return(alreadyCreatedUser, nil)
		userRepo.EXPECT().FindByUsernameOrEmail(alreadyCreatedUser.Username, "other@email.com").Times(1).Return([]model.User{alreadyCreatedUser}, nil)

		_, err := userService.CreateUser(model.User{Email: "other@email.com", Username: "username"}, "password")
		if err == nil {
			t.Errorf("err should be nil")
		}
	})
	t.Run("email already exists", func(t *testing.T) {
		alreadyCreatedUser := model.User{Email: "email", Username: "username", ID: userID}
		userRepo.EXPECT().Create(UserMatcher{
			User:     alreadyCreatedUser,
			Password: "password",
		}).Times(0).Return(alreadyCreatedUser, nil)
		userRepo.EXPECT().FindByUsernameOrEmail("anotherUsername", alreadyCreatedUser.Email).Times(1).Return([]model.User{alreadyCreatedUser}, nil)

		_, err := userService.CreateUser(model.User{Email: "email", Username: "anotherUsername"}, "password")
		if err == nil {
			t.Errorf("err should be nil")
		}
	})

	t.Run("error", func(t *testing.T) {
		expectedCreatedUser := model.User{Email: "email", Username: "username", ID: userID}
		userRepo.EXPECT().Create(UserMatcher{
			User:     expectedCreatedUser,
			Password: "password",
		}).Return(expectedCreatedUser, errors.New("cannot create user"))
		userRepo.EXPECT().FindByUsernameOrEmail(expectedCreatedUser.Username, expectedCreatedUser.Email).Times(1).Return([]model.User{}, nil)

		_, err := userService.CreateUser(model.User{Email: "email", Username: "username"}, "password")
		if err == nil {
			t.Errorf("err should not nil")
		}
	})
}

func TestUserServiceImpl_FindByUsernameOrEmail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	userRepo := mock_repos.NewMockUserRepo(controller)
	userService := services.NewUserServiceImpl(userRepo)

	userID, _ := uuid.NewUUID()

	t.Run("ok", func(t *testing.T) {
		expectedUser := model.User{Email: "email", Username: "username", ID: userID}
		expectedUserList := []model.User{expectedUser}
		userRepo.EXPECT().FindByUsernameOrEmail(expectedUser.Username, expectedUser.Email).Return(expectedUserList, nil)

		users, err := userService.FindByUsernameOrEmail(expectedUser.Username, expectedUser.Email)
		if err != nil {
			t.Errorf("err should be nil")
		}

		assert.Equal(t, expectedUserList, users)
	})

	t.Run("error", func(t *testing.T) {
		expectedUser := model.User{Email: "email", Username: "username", ID: userID}
		expectedUserList := []model.User{expectedUser}
		userRepo.EXPECT().FindByUsernameOrEmail(expectedUser.Username, expectedUser.Email).Return(expectedUserList, errors.New("some error"))

		_, err := userService.FindByUsernameOrEmail(expectedUser.Username, expectedUser.Email)
		if err == nil {
			t.Errorf("err should not be nil")
		}

		assert.Equal(t, errors.New("some error"), err)
	})
}
