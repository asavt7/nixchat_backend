package tests

import (
	"encoding/json"
	"fmt"
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

func (m *MainTestSuite) signIn() {
	var once sync.Once
	once.Do(func() {

		client := http.Client{}
		baseURL := fmt.Sprintf("http://%s:%s", m.cfg.HTTP.Host, m.cfg.HTTP.Port)
		signInURL := baseURL + "/sign-in"

		for i := 0; i < 10; i++ {
			log.Info("Trying signIn testUser")
			time.Sleep(100 * time.Millisecond)
			req, _ := http.NewRequest(echo.POST, signInURL, strings.NewReader(`{"username":"`+testUsername+`", "password":"`+testPassword+`"}`))
			req.Header.Set("Content-Type", "application/json")

			response, err := client.Do(req)
			if err != nil {
				log.Error(err)
				continue
			}
			if response.StatusCode != 200 {
				continue
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Error(err)
				continue
			}

			var res struct {
				AccessToken string `json:"access-token"`
			}
			err = json.Unmarshal(body, &res)
			if err != nil {
				log.Error(err)
				continue
			}

			m.credentials.accessToken = res.AccessToken
			log.Debugln("accessToken=" + res.AccessToken)
			return
		}
		m.FailNow("Cannot sign in as test user!")
	})
}

func (m *MainTestSuite) TestSignIn() {

	baseURL := fmt.Sprintf("http://%s:%s", m.cfg.HTTP.Host, m.cfg.HTTP.Port)
	signInURL := baseURL + "/sign-in"
	client := http.Client{}

	table := []struct {
		name                 string
		requestBody          string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:                 "no password",
			requestBody:          `{"username": "` + gofakeit.Username() + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "no username",
			requestBody:          `{"password": "` + password + `"}`,
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "invalid username",
			requestBody:          `{"password": "` + password + `","username": "` + "another name" + `"}`,
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: `{"message":"<<PRESENCE>>"}`,
		},
		{
			name:                 "invalid password",
			requestBody:          `{"password": "` + "another_pass" + `", "username": "` + username + `"}`,
			expectedResponseCode: http.StatusUnauthorized,
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
			req, _ := http.NewRequest(echo.POST, signInURL, strings.NewReader(tt.requestBody))
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
