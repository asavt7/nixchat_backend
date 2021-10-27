package handlers

import (
	"errors"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/asavt7/nixchat_backend/internal/services"
	mock_services "github.com/asavt7/nixchat_backend/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	username                 = "Bret"
	email                    = "Sincere@april.biz"
	createdUserID            = uuid.New()
	createdUserIDStr         = createdUserID.String()
	password                 = "password"
	createUserRqBody         = `{"password": "` + password + `", "name": "Leanne Graham", "username": "` + username + `","email": "` + email + `"}`
	createUserRsBodyExpected = `{"id": "` + createdUserIDStr + `","username": "` + username + `", "email": "` + email + `", "avatar_url": "` + "" + `"}`
)

func TestAPIHandler_signUp(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockSrv := mock_services.NewMockUserService(controller)
	handler := NewAPIHandler(&services.Services{UserService: mockSrv})
	e := echo.New()

	table := []struct {
		name                 string
		requestBody          string
		expectedResponseCode int
		expectedResponseBody string
		mocksBehaviourFunc   func()
	}{
		{
			name:                 "ok",
			requestBody:          createUserRqBody,
			expectedResponseCode: http.StatusCreated,
			expectedResponseBody: createUserRsBodyExpected,
			mocksBehaviourFunc: func() {
				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "no password",
			requestBody:          `{"name": "Leanne Graham", "username": "` + username + `","email": "` + email + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {

				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Times(0).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "invalid password: too short",
			requestBody:          `{"password": "` + "12345" + `","name": "Leanne Graham", "username": "` + username + `","email": "` + email + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {

				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Times(0).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "no username",
			requestBody:          `{"password": "` + password + `","name": "Leanne Graham", "email": "` + email + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {

				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Times(0).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "invalid username",
			requestBody:          `{"password": "` + password + `","name": "Leanne Graham", "username": "` + "" + `", "email": "` + email + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {

				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Times(0).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "no email",
			requestBody:          `{"password": "` + password + `","name": "Leanne Graham", "username": "` + username + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {

				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Times(0).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "invalid email",
			requestBody:          `{"password": "` + password + `","name": "Leanne Graham", "username": "` + username + `", "email": "` + "email" + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {

				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Times(0).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "invalid request body",
			requestBody:          `}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {

				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Times(0).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, nil)
			},
		},
		{
			name:                 "cannot create user",
			requestBody:          createUserRqBody,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
			mocksBehaviourFunc: func() {
				mockSrv.EXPECT().CreateUser(model.User{
					Username: username,
					Email:    email,
				}, password).Return(model.User{
					ID:           createdUserID,
					Username:     username,
					Email:        email,
					PasswordHash: "",
					AvatarURL:    "",
				}, errors.New("cannot create error"))
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocksBehaviourFunc()

			req := httptest.NewRequest(echo.POST, "/sign-up", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			if assert.NoError(t, handler.signUp(c)) {
				assert.Equal(t, tt.expectedResponseCode, rec.Code)
				jsonassert.New(t).Assertf(rec.Body.String(), tt.expectedResponseBody)
			}
		})

	}

}
