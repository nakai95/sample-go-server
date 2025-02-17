package controller

import (
	"sample-go-server/api"
	"sample-go-server/internal/adapter/presenter"
	"sample-go-server/internal/domain"
	"sample-go-server/internal/usecase"
)

type eventContr struct {
	usecase domain.EventUseCase
	pres    presenter.EventPresenter
}

type EventController interface {
	ListEvents() ([]api.EventsWithID, error)
}

func NewEventController(repo usecase.EventRepository, pres presenter.EventPresenter) EventController {
	uc := usecase.NewEventUsecase(repo)
	return &eventContr{
		usecase: uc,
		pres:    pres,
	}
}

func (c *eventContr) ListEvents() ([]api.EventsWithID, error) {
	events, err := c.usecase.ListEvents()
	if err != nil {
		return nil, err
	}
	return c.pres.PresentEvents(events), nil
}
