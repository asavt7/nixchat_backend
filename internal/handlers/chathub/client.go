package chathub

import (
	"encoding/json"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	UserID       uuid.UUID
	SessionID    uuid.UUID
	Conn         *websocket.Conn
	Hub          ClientConnectionsHub
	ToClientChan chan string
	doneCh       chan struct{}

	RedisNotificationPubSub *redis.PubSub
}

func (c *Client) Reader() {
	defer func() {
		c.Hub.Unregister(c)
	}()

	for {

		_, jsonDto, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warningf("error: %v", err)
			}
			break
		}
		log.Infof("user=%s send message=%s", c.UserID, jsonDto)

		var frontDto model.FrontMessageRq
		err = json.Unmarshal([]byte(jsonDto), &frontDto)
		if err != nil {
			log.Errorf("error unmarchal json %v", err)
		}
		frontDto.UserID = c.UserID

		message, err := json.Marshal(frontDto)
		if err != nil {
			log.Errorf("error unmarchal json %v", err)
		}

		c.Hub.SendFromClientMessage(string(message))
	}
}

func (c *Client) Writer() {
	for {
		select {
		case message := <-c.ToClientChan:

			log.Infof("send message to user=%s message=%s", c.UserID, message)

			err := c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
			//err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Error("error sending message " + err.Error())
			}
		case <-c.doneCh:
			return
		}
	}
}
