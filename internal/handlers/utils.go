package handlers

import (
	"github.com/labstack/echo/v4"
	"strings"
)

// Message - common message
type Message struct {
	Message string `json:"message"`
}

func responseMessage(status int, message string, c echo.Context) error {
	return response(status, Message{Message: message}, c)
}

func response(status int, body interface{}, c echo.Context) error {
	ctype := c.Request().Header.Get(echo.HeaderContentType)
	acceptType := c.Request().Header.Get(echo.HeaderAccept)
	if len(acceptType) == 0 {
		acceptType = ctype
	}
	switch {
	case strings.Contains(acceptType, echo.MIMEApplicationJSON):
		return c.JSON(status, body)
	default:
		return c.JSON(status, body)
	}
}
