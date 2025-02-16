//go:generate mockgen -source=$GOFILE -destination=../../../mock/mock_$GOFILE -package=mock -self_package=maptalk/mock

package repository

import (
	"fmt"
	"sample-go-server/internal/domain"
	"sample-go-server/internal/usecase"
)

type EventDataStore interface {
	AddEvent(event domain.DraftEvent) (string, error)
	GetEvents() ([]domain.Event, error)
}

type eventRepository struct {
	ds EventDataStore
}

func NewEventRepository(ds EventDataStore) usecase.EventRepository {
	return &eventRepository{ds: ds}
}

func (r *eventRepository) ListEvents() ([]domain.Event, error) {
	events, err := r.ds.GetEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}
