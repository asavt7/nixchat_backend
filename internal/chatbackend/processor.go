package chatbackend

import (
	"context"
	"encoding/json"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/asavt7/nixchat_backend/internal/repos"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type MessageProcessor struct {
}

type ChatRequestsProcessor struct {
	redisClient *redis.Client

	repo repos.Repositories
}

func NewChatRequestsProcessor(redisClient *redis.Client, repo repos.Repositories) *ChatRequestsProcessor {
	return &ChatRequestsProcessor{redisClient: redisClient, repo: repo}
}

func (r *ChatRequestsProcessor) Run() {

	go r.readFromRedis()
}

const (
	group                      = "clientMessagesGroup"
	ClientMessagesStream       = "clientMessages"
	toUserChannelPublishPrefix = "toUser-"
)

func (r *ChatRequestsProcessor) readFromRedis() {

	a, err := r.redisClient.XInfoConsumers(context.Background(), ClientMessagesStream, group).Result()
	if a == nil || err != nil {
		log.Debugf("cannot get redis consumer group info %s", err)

		// create group
		err = r.redisClient.XGroupCreateMkStream(context.Background(), ClientMessagesStream, group, "0").Err()
		if err != nil {
			log.Errorf("cannot create group %s", err)
		}
	}

	consumerID := uuid.NewString()
	for {
		entries, err := r.redisClient.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumerID,
			Streams:  []string{ClientMessagesStream, ">"},
			Count:    2,
			Block:    0,
			NoAck:    false,
		}).Result()
		if err != nil {
			log.Fatal(err)
		}

		// todo process messages (save) and push results to pub sub
		for i := 0; i < len(entries[0].Messages); i++ {
			messageID := entries[0].Messages[i].ID
			values := entries[0].Messages[i].Values

			r.processMessage(values)

			r.redisClient.XAck(context.Background(), ClientMessagesStream, group, messageID)
		}
	}
}

func (r *ChatRequestsProcessor) processMessage(values map[string]interface{}) {
	var frontRqDto model.FrontMessageRq
	frontRqDtoStr := values["message"].(string)
	log.Infof("chat back processor got %s", frontRqDtoStr)

	err := json.Unmarshal([]byte(frontRqDtoStr), &frontRqDto)
	if err != nil {
		log.Error(err.Error())
	}

	switch frontRqDto.Payload.(type) {
	case *model.ChatMessage:
		var cm model.ChatMessage = *frontRqDto.Payload.(*model.ChatMessage)
		cm.UserID = frontRqDto.UserID
		cm, err = r.repo.ChatRepo.SaveMessage(cm)
		if err != nil {
			//todo process error no such chat, and other error
			log.Errorf("cannot save message %s", err)
		}

		r.sendRsToSender(frontRqDto.UserID, model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.ChatMessageRsType,
			UserGeneratedID: frontRqDto.UserGeneratedID,
			Payload:         cm,
			ErrorMessage:    "",
		})

		r.notifyChatMembersNewMessage(cm)
	case *model.LoadUsersRq:
		users, err := r.repo.UserRepo.GetAll(model.PagedQuery{
			Size:   100,
			Offset: 0,
		})
		if err != nil {
			log.Errorf("cannot load users %s", err)
		}
		r.sendRsToSender(frontRqDto.UserID, model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.LoadUsersRsType,
			UserGeneratedID: frontRqDto.UserGeneratedID,
			Payload: model.LoadUsersRs{
				Users: users,
			},
			ErrorMessage: "",
		})
	case *model.LoadPrivateChatsRq:
		var cm model.LoadPrivateChatsRq = *frontRqDto.Payload.(*model.LoadPrivateChatsRq)
		// todo check user has access to read chat messages
		cm.UserID = frontRqDto.UserID
		//todo pagination
		chats, err := r.repo.ChatRepo.GetUserPrivateChats(cm.UserID)
		if err != nil {
			//todo process error no such chat, and other error
			log.Errorf("cannot save message %s", err)
		}

		r.sendRsToSender(frontRqDto.UserID, model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.LoadPrivateChatsRsType,
			UserGeneratedID: frontRqDto.UserGeneratedID,
			Payload: model.LoadPrivateChatsRs{
				UserID: frontRqDto.UserID,
				Chats:  chats,
			},
			ErrorMessage: "",
		})

	case *model.LoadChatsRq:
		var cm model.LoadChatsRq = *frontRqDto.Payload.(*model.LoadChatsRq)
		// todo check user has access to read chat messages
		cm.UserID = frontRqDto.UserID
		//todo pagination
		chats, err := r.repo.ChatRepo.GetUserChats(cm.UserID)
		if err != nil {
			//todo process error no such chat, and other error
			log.Errorf("cannot save message %s", err)
		}

		r.sendRsToSender(frontRqDto.UserID, model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.LoadChatsRsType,
			UserGeneratedID: frontRqDto.UserGeneratedID,
			Payload: model.LoadChatsRs{
				UserID: frontRqDto.UserID,
				Chats:  chats,
			},
			ErrorMessage: "",
		})
	case *model.LoadChatMessagesRq:
		var cm model.LoadChatMessagesRq = *frontRqDto.Payload.(*model.LoadChatMessagesRq)

		// todo check user has access to read chat messages

		//todo pagination
		messages, err := r.repo.ChatRepo.GetChatMessagesByQuery(cm.ChatID, model.Query{
			Size:            cm.Size,
			Offset:          cm.Offset,
			BeforeTimestamp: cm.BeforeTimestamp,
		})
		if err != nil {
			//todo process error no such chat, and other error
			log.Errorf("cannot save message %s", err)
		}

		r.sendRsToSender(frontRqDto.UserID, model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.LoadChatMessagesRsType,
			UserGeneratedID: frontRqDto.UserGeneratedID,
			Payload: model.LoadChatMessagesRs{
				ChatID:   cm.ChatID,
				Messages: messages,
			},
			ErrorMessage: "",
		})
	case *model.NewPrivateChatRq:
		var cm model.NewPrivateChatRq = *frontRqDto.Payload.(*model.NewPrivateChatRq)

		// TODO refactor data model for case : private user 2 user chat
		chat, err := r.repo.ChatRepo.CreatePrivateChat(frontRqDto.UserID, model.Chat{
			Name:        "",
			Title:       "",
			Description: "",
			Type:        "userToUser",
		}, cm.ToUserID)
		if err != nil {
			//todo process error no such chat, and other error
			log.Errorf("cannot create chat %s", err)
		}

		userToUserChat := model.UserToUserChat{
			ID:      chat.ID,
			Type:    chat.Type,
			UserIds: []uuid.UUID{frontRqDto.UserID, cm.ToUserID},
		}
		r.sendRsToSender(frontRqDto.UserID, model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.NewPrivateChatRsType,
			UserGeneratedID: frontRqDto.UserGeneratedID,
			Payload: model.NewPrivateChatRs{
				Chat: userToUserChat,
			},
			ErrorMessage: "",
		})

		//notify chat member new chat created
		r.sendRsToSender(cm.ToUserID, model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.NewPrivateChatRsType,
			UserGeneratedID: "",
			Payload: model.NewPrivateChatRs{
				Chat: userToUserChat,
			},
			ErrorMessage: "",
		})
	default:
		log.Warningf("unknown message type %s of payload %s", frontRqDto.MessageType, frontRqDto.Payload)
	}

}

func (r *ChatRequestsProcessor) sendRsToSender(userID uuid.UUID, fr model.FrontMessageRs) {
	jsonMes, err := json.Marshal(fr)
	if err != nil {
		log.Error(err.Error())
	}

	r.redisClient.Publish(context.Background(), toUserChannelPublishPrefix+userID.String(), jsonMes)
}

func (r *ChatRequestsProcessor) notifyChatMembersNewMessage(cm model.ChatMessage) {

	members, err := r.repo.ChatRepo.GetChatMembers(cm.ChatID)
	if err != nil {
		log.Errorf("cannot get chat members %s err : %v", cm.ChatID.String(), err)
	}

	for _, userID := range members {

		mes := model.FrontMessageRs{
			Ok:              true,
			MessageType:     model.ChatMessageNotificationType,
			UserGeneratedID: "",
			Payload:         cm,
			ErrorMessage:    "",
		}
		jsonMes, err := json.Marshal(mes)
		if err != nil {
			log.Error(err.Error())
		}
		r.redisClient.Publish(context.Background(), toUserChannelPublishPrefix+userID.String(), jsonMes)

	}
}
