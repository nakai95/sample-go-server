package domain

type DraftEvent struct {
	Name        string
	Description string
	ImageUrl    string
}

type Event struct {
	Id          string
	Name        string
	Description string
	ImageUrl    string
}

type EventUseCase interface {
	ListEvents() ([]Event, error)
}
