package firestore

import (
	"context"
	"sample-go-server/internal/adapter/repository"

	"cloud.google.com/go/firestore"
)

type firestoreClient struct {
	client         *firestore.Client
	EventDataStore repository.EventDataStore
}

func NewFirestoreClient(projectID string) (*firestoreClient, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	event := NewEventDatastore(client)
	return &firestoreClient{client: client, EventDataStore: event}, nil
}

func (fc *firestoreClient) Close() error {
	return fc.client.Close()
}
