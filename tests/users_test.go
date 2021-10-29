package tests

import (
	"fmt"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func (m *MainTestSuite) TestUsers_GetAll() {
	client := http.Client{}
	baseURL := fmt.Sprintf("http://%s:%s", m.cfg.HTTP.Host, m.cfg.HTTP.Port)
	getAllUsersURL := baseURL + "/api/v1/users"

	m.T().Run("ok", func(t *testing.T) {
		req, _ := http.NewRequest(echo.GET, getAllUsersURL, nil)
		req.Header.Set("Authorization", "Bearer "+m.credentials.accessToken)
		got, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusOK, got.StatusCode)
	})

	m.T().Run("not authorized - no token provided", func(t *testing.T) {
		req, _ := http.NewRequest(echo.GET, getAllUsersURL, nil)
		got, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusUnauthorized, got.StatusCode)
	})

	m.T().Run("not authorized - invalid provided", func(t *testing.T) {
		req, _ := http.NewRequest(echo.GET, getAllUsersURL, nil)
		req.Header.Set("Authorization", "Bearer someBadToken")
		got, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, http.StatusUnauthorized, got.StatusCode)
	})
}

func (m *MainTestSuite) TestUsers_UpdateUser() {
	client := http.Client{}
	baseURL := fmt.Sprintf("http://%s:%s", m.cfg.HTTP.Host, m.cfg.HTTP.Port)
	userResourceURL := baseURL + "/api/v1/users/" + testUserID

	m.T().Run("ok", func(t *testing.T) {
		avatarUrl := "http://new_url"

		req, _ := http.NewRequest(echo.PUT, userResourceURL, strings.NewReader(`{"avatar_url":"`+avatarUrl+`"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+m.credentials.accessToken)
		response, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, http.StatusOK, response.StatusCode)
		jsonassert.New(t).Assertf(string(body), `{"username":"`+testUsername+`","avatar_url":"`+avatarUrl+`","email":"`+testEmail+`","id":"`+testUserID+`"}`)
	})

}
