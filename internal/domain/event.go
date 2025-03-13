package domain

type DraftEvent struct {
	Name        string `json:"name" firestore:"name"`
	Description string `json:"description" firestore:"description"`
	ImageUrl    string `json:"image_url" firestore:"image_url"`
}

type Event struct {
	Id          string `json:"id" firestore:"id"`
	Name        string `json:"name" firestore:"name"`
	Description string `json:"description" firestore:"description"`
	ImageUrl    string `json:"image_url" firestore:"image_url"`
}

type EventUseCase interface {
	ListEvents() ([]Event, error)
}
