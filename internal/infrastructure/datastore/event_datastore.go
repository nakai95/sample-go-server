package datastore

import (
	"database/sql"
	"sample-go-server/internal/adapter/repository"
	"sample-go-server/internal/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const tableName = "events"

type eventDataStore struct {
	DB *sql.DB
}

func NewEventDatastore(db *sql.DB) repository.EventDataStore {
	return &eventDataStore{DB: db}
}

func (ds *eventDataStore) AddEvent(event domain.DraftEvent) (string, error) {
	id := uuid.New().String()
	query := `
        INSERT INTO ` + tableName + `  (id, name, description, image_url)
        VALUES ($1, $2, $3, $4)
    `
	_, err := ds.DB.Exec(query, id, event.Name, event.Description, event.ImageUrl)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (ds *eventDataStore) GetEvents() ([]domain.Event, error) {
	query := `
		SELECT id, name, description, image_url
		FROM ` + tableName
	rows, err := ds.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		if err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.ImageUrl); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
