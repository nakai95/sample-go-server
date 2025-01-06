//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_$GOFILE  -package=mock -self_package=sample-go-server/mock

package usecase

import (
	"sample-go-server/internal/domain"
)

type EventRepository interface {
	ListEvents() ([]domain.Event, error)
}

type EventUsecase struct {
	repo EventRepository
}

func NewEventUsecase(repo EventRepository) domain.EventUseCase {
	return &EventUsecase{
		repo: repo,
	}
}

func (u *EventUsecase) ListEvents() ([]domain.Event, error) {
	events, err := u.repo.ListEvents()
	if err != nil {
		return nil, err
	}
	return events, nil

}
