package handlers

import (
	"github.com/asavt7/nixchat_backend/internal/handlers/chathub"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)
import "github.com/gorilla/websocket"

var (
	upgrader = websocket.Upgrader{}
)

func (h *APIHandler) websocketHandler(c echo.Context) error {

	currentUser := c.Get(currentUserID).(uuid.UUID)

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}
	}(conn)

	client := chathub.Client{
		UserID:       currentUser,
		SessionID:    uuid.New(),
		Conn:         conn,
		Hub:          h.hub,
		ToClientChan: make(chan string),
	}

	h.hub.Register(&client)

	go client.Reader()
	client.Writer()

	return nil
}
