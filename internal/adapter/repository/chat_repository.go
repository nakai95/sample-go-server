//go:generate mockgen -source=$GOFILE -destination=../../../mock/mock_$GOFILE -package=mock -self_package=maptalk/mock

package repository

import (
	"sample-go-server/internal/domain"
	"sample-go-server/internal/usecase"
)

type ChatDataStore interface {
	GetMessages(roomID string) ([]domain.ChatMessage, error)
	AddMessage(message domain.ChatMessage) (string, error)
}

type chatRepository struct {
	ds ChatDataStore
}

func NewChatRepository(ds ChatDataStore) usecase.ChatRepository {
	return &chatRepository{ds: ds}
}

func (r *chatRepository) ListChatRooms() ([]domain.ChatRoom, error) {
	// dummy
	chatRooms := []domain.ChatRoom{
		{
			ID:   "a7e6dd14-1c50-a5e7-f003-951d63059f41",
			Name: "Room 1",
		},
		{
			ID:   "10b2cedc-4b88-ea16-90aa-5f0bd8a84018",
			Name: "Room 2",
		},
	}
	return chatRooms, nil
}

func (r *chatRepository) GetMessages(roomID string) ([]domain.ChatMessage, error) {
	messages, err := r.ds.GetMessages(roomID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatRepository) SaveMessage(message domain.ChatMessage) error {
	_, err := r.ds.AddMessage(message)
	if err != nil {
		return err
	}
	return nil
}
