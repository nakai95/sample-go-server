//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOFILE  -package=mock -self_package=sample-go-server/mock

package usecase

import (
	"sample-go-server/internal/domain"
)

type ChatRepository interface {
	ListChatRooms() ([]domain.ChatRoom, error)
}

type ChatUsecase struct {
}

func NewChatUsecase() domain.ChatUseCase {
	return &ChatUsecase{}
}

func (u *ChatUsecase) ListChatRooms() ([]domain.ChatRoom, error) {
	// dummy
	chatRooms := []domain.ChatRoom{
		{
			ID:   "1",
			Name: "Room 1",
		},
		{
			ID:   "2",
			Name: "Room 2",
		},
	}
	return chatRooms, nil

}
