package domain

import "time"

type ChatMessage struct {
	Id        string    `json:"id"`
	RoomId    string    `json:"roomId"`
	UserId    string    `json:"userId"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChatRoom struct {
	ID   string
	Name string
}

type ChatUseCase interface {
	ListChatRooms() ([]ChatRoom, error)
	GetMessages(roomId string, limit, offset int) ([]ChatMessage, error)
	SaveMessage(message ChatMessage) error
}
