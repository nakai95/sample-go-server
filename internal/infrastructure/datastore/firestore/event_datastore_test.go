package firestore

import (
	"context"
	"sample-go-server/internal/domain"
	"sample-go-server/test"
	"testing"

	"github.com/google/uuid"
)

const projectID = "test-projects"

func TestAddEvent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup postgres container
	client := test.SetupFirestoreContainer(t, projectID)

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
			ds := NewEventDatastore(client)
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
	client := test.SetupFirestoreContainer(t, projectID)

	ctx := context.Background()

	// insert dummy data
	dummyDraftEvent := domain.DraftEvent{
		Name:        "Event 1",
		Description: "homines dum docent discunt.",
		ImageUrl:    "https://picsum.photos/seed/example1/150",
	}

	id := uuid.New().String()

	_, err := client.Collection("events").Doc(id).Set(ctx, dummyDraftEvent)
	if err != nil {
		t.Fatal("failed to add dummy data:", err)
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
			ds := NewEventDatastore(client)
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
