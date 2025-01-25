package controller

import (
	"sample-go-server/api"
	"sample-go-server/internal/adapter/presenter"
	"sample-go-server/internal/domain"
	"sample-go-server/internal/usecase"
)

type chatCtrl struct {
	usecase domain.ChatUseCase
	pres    presenter.ChatPresenter
}

type ChatController interface {
	ListChatRooms() ([]api.ChatRoom, error)
}

func NewChatController(pres presenter.ChatPresenter) ChatController {
	uc := usecase.NewChatUsecase()
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
