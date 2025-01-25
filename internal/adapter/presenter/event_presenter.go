//go:generate mockgen -source=$GOFILE -destination=../../../mock/mock_$GOFILE -package=mock -self_package=maptalk/mock

package presenter

import (
	"sample-go-server/api"
	"sample-go-server/internal/domain"
)

type EventPresenter interface {
	PresentEvents(events []domain.Event) []api.EventsWithID
}

type presenter struct {
}

func NewEventPresenter() EventPresenter {
	return &presenter{}
}

func (p *presenter) PresentEvents(events []domain.Event) []api.EventsWithID {
	// Convert domain.Event to api.EventsWithID
	var eventsWithID []api.EventsWithID
	for _, event := range events {
		eventWithID := api.EventsWithID{
			Id:          event.Id,
			Name:        event.Name,
			Description: event.Description,
			ImageUrl:    event.ImageUrl,
		}
		eventsWithID = append(eventsWithID, eventWithID)
	}
	return eventsWithID
}
