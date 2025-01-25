package server

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

type ChatRoom struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
	lastActive time.Time
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		lastActive: time.Now(),
	}
}

func (room *ChatRoom) Run() {
	for {
		select {
		case client := <-room.Register:
			room.mu.Lock()
			room.Clients[client] = true
			room.lastActive = time.Now()
			room.mu.Unlock()
		case client := <-room.Unregister:
			room.mu.Lock()
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				close(client.Send)
			}
			room.lastActive = time.Now()
			room.mu.Unlock()
		case message := <-room.Broadcast:
			room.mu.Lock()
			for client := range room.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(room.Clients, client)
				}
			}
			room.lastActive = time.Now()
			room.mu.Unlock()
		}
	}
}

var chatRooms = make(map[string]*ChatRoom)
var chatRoomsMu sync.Mutex

func getOrCreateChatRoom(chatID string) *ChatRoom {
	chatRoomsMu.Lock()
	defer chatRoomsMu.Unlock()

	room, exists := chatRooms[chatID]
	if !exists {
		room = NewChatRoom()
		chatRooms[chatID] = room
		go room.Run()
		go monitorChatRoom(chatID, room)
	}
	return room
}

func monitorChatRoom(chatID string, room *ChatRoom) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		chatRoomsMu.Lock()
		if time.Since(room.lastActive) > 5*time.Minute {
			delete(chatRooms, chatID)
			chatRoomsMu.Unlock()
			return
		}
		chatRoomsMu.Unlock()
	}
}
