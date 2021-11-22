package chathub

import (
	"context"
	"github.com/asavt7/nixchat_backend/internal/chatbackend"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=./mocks/mocks.go -source=./hub.go
type ClientConnectionsHub interface {
	ClientConnectionsRegister
	ClientConnectionsMessageProcessor
}

type ClientConnectionsRegister interface {
	Register(client *Client)
	Unregister(client *Client)
}

type ClientConnectionsMessageProcessor interface {
	SendFromClientMessage(mes string)
}

type Hub struct {
	redisClient *redis.Client

	clients      map[uuid.UUID]*Client
	registerCh   chan *Client
	unregisterCh chan *Client

	fromClientMessageCh chan string
}

func (h *Hub) Register(client *Client) {
	h.registerCh <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregisterCh <- client
}

func (h *Hub) SendFromClientMessage(mes string) {
	h.fromClientMessageCh <- mes
}

func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		redisClient:         redisClient,
		clients:             make(map[uuid.UUID]*Client),
		registerCh:          make(chan *Client, 256),
		unregisterCh:        make(chan *Client, 256),
		fromClientMessageCh: make(chan string, 256),
	}
}

var toUserChannelPublishPrefix = "toUser-"

func (h *Hub) Run() error {

	go func() {

		for {
			select {
			case client := <-h.registerCh:
				log.Infof("Hub register new client %s", client.UserID)
				h.clients[client.UserID] = client

				//todo unsubscribe
				go func() {
					ps := h.redisClient.Subscribe(context.Background(), toUserChannelPublishPrefix+client.UserID.String())
					client.RedisNotificationPubSub = ps
					for {
						select {
						default:
							message, err := ps.ReceiveMessage(context.Background())
							if err != nil {
								log.Errorf("cannot subscribe %s", err)
								return
							}
							log.Infof("read message from SUB %v", message)
							h.clients[client.UserID].ToClientChan <- message.Payload
						}

					}

				}()

			//fixme panic in frequent connect\disconnect cases
			case client := <-h.unregisterCh:
				log.Infof("Hub unregister new client %s", client.UserID)
				err := h.clients[client.UserID].RedisNotificationPubSub.Unsubscribe(context.Background(), toUserChannelPublishPrefix+client.UserID.String())
				if err != nil {
					log.Errorf("Error unsubscribe redis notifications %s", err)
				}
				delete(h.clients, client.UserID)

			case clientMes := <-h.fromClientMessageCh:
				log.Infof("Hub get client message %v", clientMes)

				h.processMessage(clientMes)
			}
		}
	}()

	return nil
}

func (h *Hub) processMessage(mes interface{}) {

	err := h.redisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: chatbackend.ClientMessagesStream,
		Values: map[string]interface{}{
			"message": mes,
		}}).Err()

	if err != nil {
		log.Errorf("cannot send new message to topic %v", err)
	}

	/*
		for userID, client := range h.clients {
			log.Infof("Hub send message to user=%s message=%v", userID, mes)
			client.toClientChan <- mes
		}
	*/
}
