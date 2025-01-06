package domain

type Event struct {
	ID          string
	Name        string
	Description string
	ImageURL    string
}

type EventUseCase interface {
	ListEvents() ([]Event, error)
}
