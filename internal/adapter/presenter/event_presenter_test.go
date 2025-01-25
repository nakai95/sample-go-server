package presenter

import (
	"reflect"
	"sample-go-server/api"
	"sample-go-server/internal/domain"
	"testing"
)

func TestPresentEvents(t *testing.T) {
	tests := []struct {
		name   string
		events []domain.Event
		want   []api.EventsWithID
	}{
		{
			name: "success",
			events: []domain.Event{
				{
					Id:          "1",
					Name:        "Event 1",
					Description: "homines dum docent discunt.",
					ImageUrl:    "https://picsum.photos/seed/example1/150",
				},
			},
			want: []api.EventsWithID{
				{
					Id:          "1",
					Name:        "Event 1",
					Description: "homines dum docent discunt.",
					ImageUrl:    "https://picsum.photos/seed/example1/150",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewEventPresenter()
			if got := p.PresentEvents(tt.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PresentEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}
