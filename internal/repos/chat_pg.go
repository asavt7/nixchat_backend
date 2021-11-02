package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/asavt7/nixchat_backend/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type ChatPgRepo struct {
	db *sqlx.DB
}

func (c *ChatPgRepo) GetChatMessagesByQuery(chatID uuid.UUID, query model.Query) ([]model.ChatMessage, error) {

	chatMessages := make([]model.ChatMessage, 0)

	queryStr := "SELECT * FROM nix.messages WHERE chatid=$1 AND timestamp<$2 ORDER BY timestamp DESC LIMIT $3 OFFSET $4"
	err := c.db.Select(&chatMessages, queryStr, chatID, query.BeforeTimestamp, query.Size, query.Offset)
	if err != nil {
		return chatMessages, err
	}

	return chatMessages, nil
}

func (c *ChatPgRepo) GetUserPrivateChats(userID uuid.UUID) ([]model.UserToUserChat, error) {
	chats := make([]model.UserToUserChat, 0)

	query := "SELECT id, types[1] as type, userids FROM ( SELECT uc.chatid as id, ARRAY_AGG(c.type) AS types, ARRAY_AGG(uc.userid) AS userids FROM nix.user_chats uc INNER JOIN nix.chats c ON c.id=uc.chatid WHERE c.type='userToUser'  GROUP BY uc.chatid) ucs WHERE $1 = ANY(userids);"
	rows, err := c.db.Query(query, userID)
	if err != nil {
		return chats, err
	}
	for rows.Next() {
		var chat model.UserToUserChat
		if err := rows.Scan(&chat.ID, &chat.Type, pq.Array(&chat.UserIds)); err != nil {
			return chats, err
		}
		chats = append(chats, chat)
	}
	return chats, rows.Err()
}

func (c *ChatPgRepo) GetChatMessages(chatID uuid.UUID) ([]model.ChatMessage, error) {
	chatMessages := make([]model.ChatMessage, 0)

	query := "SELECT * FROM nix.messages WHERE chatid=$1 "
	err := c.db.Select(&chatMessages, query, chatID)
	if err != nil {
		return chatMessages, err
	}

	return chatMessages, nil
}

func (c *ChatPgRepo) GetUserChats(userID uuid.UUID) ([]model.Chat, error) {
	chats := make([]model.Chat, 0)

	query := "SELECT * FROM nix.chats WHERE id IN (SELECT chatid FROM nix.user_chats WHERE userid=$1 ) AND type!='userToUser'"
	err := c.db.Select(&chats, query, userID)
	if err != nil {
		return chats, err
	}

	return chats, nil
}

func (c *ChatPgRepo) SaveMessage(mes model.ChatMessage) (model.ChatMessage, error) {
	var createdMessageID uuid.UUID

	query := "INSERT INTO nix.messages (userid, chatid, text, timestamp) VALUES ($1,$2,$3,$4) RETURNING id"
	err := c.db.Get(&createdMessageID, query, mes.UserID, mes.ChatID, mes.Text, mes.Timestamp)
	if err != nil {
		return model.ChatMessage{}, err
	}

	mes.ID = createdMessageID
	return mes, nil
}

func (c *ChatPgRepo) GetChatMembers(chatID uuid.UUID) ([]uuid.UUID, error) {
	var members []uuid.UUID = make([]uuid.UUID, 0)

	query := "SELECT (userid) FROM nix.user_chats WHERE chatid=$1"
	err := c.db.Select(&members, query, chatID)
	if err != nil {
		return members, err
	}

	return members, nil
}

func (c *ChatPgRepo) CreatePrivateChat(ownerID uuid.UUID, chat model.Chat, toUserID uuid.UUID) (model.Chat, error) {
	var createdChatID uuid.UUID

	//check private chat not exists
	var commonChats []model.Chat
	checkChatExistsQuery := fmt.Sprintf("SELECT * FROM nix.chats WHERE id in ( (SELECT chatid FROM nix.user_chats WHERE userid=$1 AND relation='userToUser') INTERSECT (SELECT chatid FROM nix.user_chats WHERE userid=$2 AND relation='userToUser') );")
	err := c.db.Select(&commonChats, checkChatExistsQuery, ownerID, toUserID)
	if err != nil {
		if err != sql.ErrNoRows {
			return model.Chat{}, err
		}
	}
	if len(commonChats) > 0 {
		return commonChats[0], errors.New("Chat already exists ")
	}

	ctx := context.Background()
	query := fmt.Sprintf("INSERT INTO nix.chats (name,title,description,type) VALUES ($1,$2,$3,$4) RETURNING ID;")

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Chat{}, err
	}

	err = tx.QueryRowContext(ctx, query, chat.Name, chat.Title, chat.Description, chat.Type).Scan(&createdChatID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Error(err)
		}
		return model.Chat{}, err
	}

	queryOwner := fmt.Sprintf("INSERT INTO nix.user_chats (userid,chatid,relation) VALUES ($1,$2,$3);")
	result, err := tx.ExecContext(ctx, queryOwner, ownerID, createdChatID, "userToUser")
	log.Info(result)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Error(err)
		}
		return model.Chat{}, err
	}

	queryJoinedUser := fmt.Sprintf("INSERT INTO nix.user_chats (userid,chatid,relation) VALUES ($1,$2,$3);")
	result, err = tx.ExecContext(ctx, queryJoinedUser, toUserID, createdChatID, "userToUser")
	log.Info(result)

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Error(err)
		}
		return model.Chat{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Chat{}, err
	}

	chat.ID = createdChatID
	return chat, nil
}

func NewChatPgRepo(db *sqlx.DB) *ChatPgRepo {
	return &ChatPgRepo{db: db}
}
