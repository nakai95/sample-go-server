//go:generate mockgen -source=$GOFILE -destination=../../../mock/mock_$GOFILE -package=mock -self_package=maptalk/mock

package repository

import (
	"fmt"
	"sample-go-server/internal/domain"
	"sample-go-server/internal/usecase"
)

type DataStore interface {
	// Add your data store methods here
}

type repository struct {
	ds DataStore
}

func NewEventRepository(ds DataStore) usecase.EventRepository {
	return &repository{ds: ds}
}

func (r *repository) ListEvents() ([]domain.Event, error) {
	// create dummy data
	events := make([]domain.Event, 10)
	for i := 0; i < 10; i++ {
		events[i] = domain.Event{
			ID:          fmt.Sprintf("%d", i),
			Name:        fmt.Sprintf("Event %d", i),
			Description: "homines dum docent discunt.",
			ImageURL:    fmt.Sprintf("https://picsum.photos/seed/example%d/150", i),
		}
	}
	return events, nil
}
