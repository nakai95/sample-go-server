//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOFILE  -package=mock -self_package=sample-go-server/mock

package usecase

import (
	"sample-go-server/internal/domain"
)

type ChatRepository interface {
	ListChatRooms() ([]domain.ChatRoom, error)
	GetMessages(roomID string) ([]domain.ChatMessage, error)
	SaveMessage(message domain.ChatMessage) error
}

type ChatUsecase struct {
	repository ChatRepository
}

func NewChatUsecase(repository ChatRepository) domain.ChatUseCase {
	return &ChatUsecase{repository: repository}
}

func (u *ChatUsecase) ListChatRooms() ([]domain.ChatRoom, error) {
	chatRooms, err := u.repository.ListChatRooms()
	if err != nil {
		return nil, err
	}
	return chatRooms, nil
}

func (u *ChatUsecase) GetMessages(roomID string) ([]domain.ChatMessage, error) {
	messages, err := u.repository.GetMessages(roomID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (u *ChatUsecase) SaveMessage(message domain.ChatMessage) error {
	err := u.repository.SaveMessage(message)
	if err != nil {
		return err
	}
	return nil
}
