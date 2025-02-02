package datastore

import (
	"sample-go-server/internal/domain"
	"sample-go-server/test"
	"testing"
)

func TestAddEvent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup postgres container
	db, err := test.SetupPostgresContainer(t)
	if err != nil {
		t.Fatal("failed to setup postgres container:", err)
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

	// setup postgres container
	db, err := test.SetupPostgresContainer(t)
	if err != nil {
		t.Fatal("failed to setup postgres container:", err)
	}

	// insert dummy data
	dummyData := []string{
		`INSERT INTO events (id, name, description, image_url) VALUES ('6cf15595-cd47-40d9-ab99-89c4527e974f', 'Event 0', 'homines dum docent discunt.', 'https://picsum.photos/seed/example0/150')`,
	}

	for _, query := range dummyData {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("failed to insert dummy data: %v", err)
		}
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
