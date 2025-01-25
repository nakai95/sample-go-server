package controller

import (
	"reflect"
	"sample-go-server/api"
	"sample-go-server/internal/domain"
	"sample-go-server/mock"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestListEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock
	mockRepo := mock.NewMockEventRepository(ctrl)
	mockPres := mock.NewMockEventPresenter(ctrl)

	dummyEvent := domain.Event{
		Id:          "1",
		Name:        "Event 1",
		Description: "homines dum docent discunt.",
		ImageUrl:    "https://picsum.photos/seed/example1/150",
	}

	dummyEventWithID := api.EventsWithID{
		Id:          dummyEvent.Id,
		Name:        dummyEvent.Name,
		Description: dummyEvent.Description,
		ImageUrl:    dummyEvent.ImageUrl,
	}

	mockRepo.EXPECT().ListEvents().Return([]domain.Event{dummyEvent}, nil)
	mockPres.EXPECT().PresentEvents([]domain.Event{dummyEvent}).Return([]api.EventsWithID{
		dummyEventWithID,
	})

	// controller
	c := NewEventController(mockRepo, mockPres)

	// when
	got, err := c.ListEvents()

	// then
	want := []api.EventsWithID{dummyEventWithID}

	if !reflect.DeepEqual(got, want) || err != nil {
		t.Errorf("ListEvents() = %v, %v; want %v, nil", got, err, want)
	}
}
