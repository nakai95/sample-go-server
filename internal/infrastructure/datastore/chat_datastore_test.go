package datastore

import (
	"sample-go-server/internal/domain"
	"sample-go-server/test"
	"testing"
	"time"
)

func TestAddMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup postgres container
	db, err := test.SetupPostgresContainer(t)
	if err != nil {
		t.Fatal("failed to setup postgres container:", err)
	}

	dummyChatMessage := domain.ChatMessage{
		Id:        "1",
		RoomId:    "a7e6dd14-1c50-a5e7-f003-951d63059f41",
		UserId:    "a7f877a3-1875-5d0c-39d5-daf42b94dcff",
		Message:   "hello",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		message domain.ChatMessage
		wantErr bool
	}{
		{
			name:    "success case: add event and return id",
			message: dummyChatMessage,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewChatDatastore(db)
			id, err := ds.AddMessage(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id == "" {
				t.Errorf("AddMessage() = %v, want %v", id, "not empty")
			}
		})
	}
}

func TestGetMessages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup postgres container
	db, err := test.SetupPostgresContainer(t)
	if err != nil {
		t.Fatal("failed to setup postgres container:", err)
	}

	// insert dummy data
	dummyData := []string{
		`INSERT INTO chat_messages (id, room_id, user_id, message, created_at) VALUES ('02564d88-1456-54e5-8a7d-9dc70a5a9f2c', 'a7e6dd14-1c50-a5e7-f003-951d63059f41', 'a7f877a3-1875-5d0c-39d5-daf42b94dcff', 'hello', '2021-09-01 00:00:00')`,
	}

	for _, data := range dummyData {
		_, err := db.Exec(data)
		if err != nil {
			t.Fatal("failed to insert dummy data:", err)
		}
	}

	tests := []struct {
		name    string
		roomId  string
		wantErr bool
	}{
		{
			name:    "success case: get messages",
			roomId:  "a7e6dd14-1c50-a5e7-f003-951d63059f41",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewChatDatastore(db)
			messages, err := ds.GetMessages(tt.roomId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(messages) == 0 {
				t.Errorf("GetMessages() = %v, want %v", len(messages), "not empty")
			}
		})
	}
}
