package presenter

import (
	"reflect"
	"sample-go-server/internal/domain"
	"testing"
)

func TestPresentEvents(t *testing.T) {
	tests := []struct {
		name   string
		events []domain.Event
		want   []domain.Event
	}{
		{
			name: "success",
			events: []domain.Event{
				{
					ID:          "1",
					Name:        "Event 1",
					Description: "homines dum docent discunt.",
					ImageURL:    "https://picsum.photos/seed/example1/150",
				},
			},
			want: []domain.Event{
				{
					ID:          "1",
					Name:        "Event 1",
					Description: "homines dum docent discunt.",
					ImageURL:    "https://picsum.photos/seed/example1/150",
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
