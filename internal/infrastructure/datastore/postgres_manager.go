package datastore

import (
	"database/sql"
	"fmt"
	"sample-go-server/internal/adapter/repository"

	_ "github.com/lib/pq"
)

type postgresManager struct {
	db             *sql.DB
	EventDataStore repository.EventDataStore
	ChatDataStore  repository.ChatDataStore
}

func NewPostgresManager() (*postgresManager, error) {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=user password=password dbname=sample sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	event := NewEventDatastore(db)
	chat := NewChatDatastore(db)

	return &postgresManager{db: db, EventDataStore: event, ChatDataStore: chat}, nil
}

func (manager *postgresManager) Close() error {
	return manager.db.Close()
}
