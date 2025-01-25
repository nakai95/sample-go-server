package domain

// WebSocketMessage represents the structure of the message sent/received via WebSocket
type WebSocketMessage struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ChatRoom struct {
	ID   string
	Name string
}

type ChatUseCase interface {
	ListChatRooms() ([]ChatRoom, error)
}
