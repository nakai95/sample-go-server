package datastore

import (
	"database/sql"
	"sample-go-server/internal/adapter/repository"
	"sample-go-server/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const chatTableName = "chat_messages"

type chatDataStore struct {
	DB *sql.DB
}

func NewChatDatastore(db *sql.DB) repository.ChatDataStore {
	return &chatDataStore{DB: db}
}

func (ds *chatDataStore) GetMessages(roomID string) ([]domain.ChatMessage, error) {
	query := `
		SELECT id, room_id, user_id, message, created_at
		FROM ` + chatTableName + `
		WHERE room_id = $1
	`
	rows, err := ds.DB.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var message domain.ChatMessage
		if err := rows.Scan(&message.Id, &message.RoomId, &message.UserId, &message.Message, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (ds *chatDataStore) AddMessage(message domain.ChatMessage) (string, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO ` + chatTableName + ` (id, room_id, user_id, message, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := ds.DB.Exec(query, id, message.RoomId, message.UserId, message.Message, message.CreatedAt)
	if err != nil {
		return "", err
	}

	return id, nil
}
