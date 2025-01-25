package usecase

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

	// mock
	successRepo := mock.NewMockEventRepository(ctrl)
	errorRepo := mock.NewMockEventRepository(ctrl)

	dummyEvent := domain.Event{
		Id:          "1",
		Name:        "Event 1",
		Description: "homines dum docent discunt.",
		ImageUrl:    "https://picsum.photos/seed/example1/150",
	}

	successRepo.EXPECT().ListEvents().Return([]domain.Event{dummyEvent}, nil).AnyTimes()
	errorRepo.EXPECT().ListEvents().Return(nil, fmt.Errorf("error")).AnyTimes()

	tests := []struct {
		name    string
		repo    *mock.MockEventRepository
		want    []domain.Event
		wantErr bool
	}{
		{
			name: "success",
			repo: successRepo,
			want: []domain.Event{
				{
					Id:          "1",
					Name:        "Event 1",
					Description: "homines dum docent discunt.",
					ImageUrl:    "https://picsum.photos/seed/example1/150",
				},
			},
			wantErr: false,
		},
		{
			name:    "error",
			repo:    errorRepo,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewEventUsecase(tt.repo)
			got, err := uc.ListEvents()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}
