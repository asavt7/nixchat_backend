package tests

import (
	"encoding/json"
	"fmt"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	username = "Bret"
	email    = "Sincere@april.biz"
	password = "password"

	testUserID = "CHANGE_ME"
)

const (
	testUsername = "testUsername"
	testEmail    = "test@gmail.com"
	testPassword = "test@gmail.com"
)

func (m *MainTestSuite) registerUser() {
	var once sync.Once
	once.Do(func() {
		baseURL := fmt.Sprintf("http://%s:%s", m.cfg.HTTP.Host, m.cfg.HTTP.Port)
		signUpURL := baseURL + "/sign-up"
		client := http.Client{}

		for i := 0; i < 10; i++ {
			log.Info("Trying register testUser")
			time.Sleep(100 * time.Millisecond)
			req, _ := http.NewRequest(echo.POST, signUpURL, strings.NewReader(fmt.Sprintf(`{"username":"%s","email":"%s","password":"%s"}`, testUsername, testEmail, testPassword)))
			req.Header.Set("Content-Type", "application/json")
			response, err := client.Do(req)
			if err != nil || response.StatusCode != 201 {
				log.Warn(err)
			}
			respBody, err := io.ReadAll(response.Body)
			if err != nil {
				log.Warn(err)
				continue
			}

			var createdTestUser model.User
			if err := json.Unmarshal(respBody, &createdTestUser); err != nil {
				log.Warn(err)
				continue
			}
			testUserID = createdTestUser.ID.String()

			return
		}
		m.FailNow("Cannot register test user!")
	})

}

func (m *MainTestSuite) TestSignUp() {
	baseURL := fmt.Sprintf("http://%s:%s", m.cfg.HTTP.Host, m.cfg.HTTP.Port)
	signUpURL := baseURL + "/sign-up"
	client := http.Client{}

	table := []struct {
		name                 string
		requestBody          string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:                 "ok - new user",
			requestBody:          `{"password": "` + password + `", "name": "Leanne Graham", "username": "` + username + `","email": "` + email + `"}`,
			expectedResponseCode: http.StatusCreated,
			expectedResponseBody: `{"id": "<<PRESENCE>>","username": "` + username + `", "email": "` + email + `", "avatar_url": "` + "" + `"}`,
		},
		{
			name:                 "user already exists",
			requestBody:          `{"password": "` + password + `", "name": "Leanne Graham", "username": "` + username + `","email": "` + email + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "no password",
			requestBody:          `{"name": "Leanne Graham", "username": "` + gofakeit.Username() + `","email": "` + gofakeit.Email() + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "invalid password: too short",
			requestBody:          `{"password": "` + "12345" + `","name": "Leanne Graham", "username": "` + gofakeit.Username() + `","email": "` + gofakeit.Email() + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "no username",
			requestBody:          `{"password": "` + gofakeit.Password(true, true, true, true, false, 10) + `","name": "Leanne Graham", "email": "` + gofakeit.Email() + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "invalid username",
			requestBody:          `{"password": "` + password + `","name": "Leanne Graham", "username": "` + "" + `", "email": "` + gofakeit.Email() + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "no email",
			requestBody:          `{"password": "` + password + `","name": "Leanne Graham", "username": "` + gofakeit.Username() + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "invalid email",
			requestBody:          `{"password": "` + password + `","name": "Leanne Graham", "username": "` + gofakeit.Username() + `", "email": "` + "email" + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "invalid request body",
			requestBody:          `}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "empty body",
			requestBody:          `{}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
	}

	for _, tt := range table {
		m.T().Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(echo.POST, signUpURL, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			response, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error(err)
			}

			jsonassert.New(t).Assertf(string(body), tt.expectedResponseBody)
			assert.Equal(t, tt.expectedResponseCode, response.StatusCode)
		})
	}
}
