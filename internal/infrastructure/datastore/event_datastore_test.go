package datastore

import (
	"database/sql"
	"sample-go-server/internal/domain"
	"testing"
)

func TestAddEvent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// create db connection
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=user password=password dbname=sample sslmode=disable")
	if err != nil {
		t.Fatal("failed to open database:", err)
	}

	dummyDraftEvent := domain.DraftEvent{
		Name:        "Event 1",
		Description: "homines dum docent discunt.",
		ImageUrl:    "https://picsum.photos/seed/example1/150",
	}

	tests := []struct {
		name    string
		event   domain.DraftEvent
		wantErr bool
	}{
		{
			name:    "success case: add event and return id",
			event:   dummyDraftEvent,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewEventDatastore(db)
			id, err := ds.AddEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id == "" {
				t.Errorf("AddEvent() = %v, want %v", id, "not empty")
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// create db connection
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=user password=password dbname=sample sslmode=disable")
	if err != nil {
		t.Fatal("failed to open database:", err)
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success case: get events",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewEventDatastore(db)
			events, err := ds.GetEvents()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(events) == 0 {
				t.Errorf("GetEvents() = %v, want %v", len(events), "not empty")
			}
		})
	}
}
