package repository

import (
	"fmt"
	"reflect"
	"sample-go-server/internal/domain"
	"sample-go-server/mock"

	"testing"

	"go.uber.org/mock/gomock"
)

func TestListEvents(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock DataStore
	ds := mock.NewMockDataStore(ctrl)

	// create UserRepository
	r := NewEventRepository(ds)

	// when
	got, err := r.ListEvents()

	// then
	want := make([]domain.Event, 10)
	for i := 0; i < 10; i++ {
		want[i] = domain.Event{
			ID:          fmt.Sprintf("%d", i),
			Name:        fmt.Sprintf("Event %d", i),
			Description: "homines dum docent discunt.",
			ImageURL:    fmt.Sprintf("https://picsum.photos/seed/example%d/150", i),
		}
	}

	// compare
	if !reflect.DeepEqual(got, want) || err != nil {
		t.Errorf("ListEvents() = %v, %v, want match for %v, nil", got, err, want)
	}
}
