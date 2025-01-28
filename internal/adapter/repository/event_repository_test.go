package repository

import (
	"reflect"
	"sample-go-server/internal/domain"
	"sample-go-server/mock"

	"testing"

	"go.uber.org/mock/gomock"
)

func TestListEvents(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dummyEvent := domain.Event{
		Id:          "6cf15595-cd47-40d9-ab99-89c4527e974f",
		Name:        "Event 0",
		Description: "homines dum docent discunt.",
		ImageUrl:    "https://picsum.photos/seed/example0/150",
	}

	// mock DataStore
	ds := mock.NewMockDataStore(ctrl)
	ds.EXPECT().GetEvents().Return([]domain.Event{
		dummyEvent,
	}, nil).Times(1)

	// create UserRepository
	r := NewEventRepository(ds)

	// when
	got, err := r.ListEvents()

	// then
	want := []domain.Event{
		dummyEvent,
	}

	// compare
	if !reflect.DeepEqual(got, want) || err != nil {
		t.Errorf("ListEvents() = %v, %v, want match for %v, nil", got, err, want)
	}
}
