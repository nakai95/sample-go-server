package firestore

import (
	"context"

	"sample-go-server/constants"
	"sample-go-server/internal/adapter/repository"
	"sample-go-server/internal/domain"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

type eventDataStore struct {
	client *firestore.Client
}

func NewEventDatastore(
	client *firestore.Client) repository.EventDataStore {
	return &eventDataStore{client: client}
}

func (ds *eventDataStore) AddEvent(event domain.DraftEvent) (string, error) {
	id := uuid.New().String()
	ctx := context.Background()
	_, err := ds.client.Collection(constants.EVENTS_COLLECTION_NAME).Doc(id).Set(ctx, event)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (ds *eventDataStore) GetEvents() ([]domain.Event, error) {
	ctx := context.Background()

	iter := ds.client.Collection(constants.EVENTS_COLLECTION_NAME).Documents(ctx)
	defer iter.Stop()

	var events []domain.Event
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []domain.Event{}, err
		}

		var event domain.Event
		err = doc.DataTo(&event)
		if err != nil {
			return []domain.Event{}, err
		}
		events = append(events, event)
	}

	return events, nil
}
