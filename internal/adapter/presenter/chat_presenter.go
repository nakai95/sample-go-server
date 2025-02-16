//go:generate mockgen -source=$GOFILE -destination=../../../mock/mock_$GOFILE -package=mock -self_package=maptalk/mock

package presenter

import (
	"sample-go-server/api"
	"sample-go-server/internal/domain"
)

type ChatPresenter interface {
	PresentChatRooms(rooms []domain.ChatRoom) []api.ChatRoom
	PresentChatMessages(messages []domain.ChatMessage) []api.ChatMessage
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

func (p *chatPres) PresentChatMessages(messages []domain.ChatMessage) []api.ChatMessage {
	// Convert domain.ChatMessage to api.ChatMessage
	var chatMessages []api.ChatMessage
	for _, message := range messages {
		chatMessage := api.ChatMessage{
			Id:        message.Id,
			RoomId:    message.RoomId,
			UserId:    message.UserId,
			Message:   message.Message,
			CreatedAt: &message.CreatedAt,
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	return chatMessages
}
