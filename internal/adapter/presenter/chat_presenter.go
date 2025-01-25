//go:generate mockgen -source=$GOFILE -destination=../../../mock/mock_$GOFILE -package=mock -self_package=maptalk/mock

package presenter

import (
	"sample-go-server/api"
	"sample-go-server/internal/domain"
)

type ChatPresenter interface {
	PresentChatRooms(rooms []domain.ChatRoom) []api.ChatRoom
}

type chatPres struct {
}

func NewChatPresenter() ChatPresenter {
	return &chatPres{}
}

func (p *chatPres) PresentChatRooms(rooms []domain.ChatRoom) []api.ChatRoom {
	// Convert domain.ChatRoom to api.ChatRoom
	var chatRooms []api.ChatRoom
	for _, room := range rooms {
		chatRoom := api.ChatRoom{
			Id:   room.ID,
			Name: room.Name,
		}
		chatRooms = append(chatRooms, chatRoom)
	}
	return chatRooms
}
