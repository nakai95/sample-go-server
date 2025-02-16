package controller

import (
	"sample-go-server/api"
	"sample-go-server/internal/adapter/presenter"
	"sample-go-server/internal/domain"
	"sample-go-server/internal/usecase"
	"time"
)

type chatCtrl struct {
	usecase domain.ChatUseCase
	pres    presenter.ChatPresenter
}

type ChatController interface {
	ListChatRooms() ([]api.ChatRoom, error)
	GetMessages(roomID string) ([]api.ChatMessage, error)
	SaveMessage(message api.ChatMessage) error
}

func NewChatController(repo usecase.ChatRepository, pres presenter.ChatPresenter) ChatController {
	uc := usecase.NewChatUsecase(repo)
	return &chatCtrl{
		usecase: uc,
		pres:    pres,
	}
}

func (c *chatCtrl) ListChatRooms() ([]api.ChatRoom, error) {
	Chats, err := c.usecase.ListChatRooms()
	if err != nil {
		return nil, err
	}
	return c.pres.PresentChatRooms(Chats), nil
}

func (c *chatCtrl) GetMessages(roomID string) ([]api.ChatMessage, error) {
	messages, err := c.usecase.GetMessages(roomID)
	if err != nil {
		return nil, err
	}
	return c.pres.PresentChatMessages(messages), nil
}

func (c *chatCtrl) SaveMessage(message api.ChatMessage) error {
	draftMessage := domain.ChatMessage{
		RoomId:    message.RoomId,
		UserId:    message.UserId,
		Message:   message.Message,
		CreatedAt: time.Now(),
	}
	err := c.usecase.SaveMessage(draftMessage)
	if err != nil {
		return err
	}
	return nil
}
