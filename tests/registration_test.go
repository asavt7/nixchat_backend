package tests

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

var (
	username = "Bret"
	email    = "Sincere@april.biz"
	password = "password"
)

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
