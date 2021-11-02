package consumer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/ozonmp/lgc-location-api/internal/mocks"
	"github.com/ozonmp/lgc-location-api/internal/model"
)

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	eventsChan := make(chan model.LocationEvent, 32)

	ctx, _ := context.WithCancel(context.Background())

	c := NewDbConsumer(1, 1, time.Second, repo, eventsChan)
	c.Start(ctx)
	c.Close()
}

func TestLockAndWriteToChan(t *testing.T) {
	events := []model.LocationEvent{
		{
			ID:     1,
			Type:   model.Created,
			Status: model.Deferred,
			Entity: &model.Location{ID: 1, Title: "L1"},
		},
	}

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)

	gomock.InOrder(
		repo.EXPECT().Lock(gomock.Any()).Return(events, nil).Times(1),
		repo.EXPECT().Lock(gomock.Any()).Return(nil, errors.New("failed to lock events")).AnyTimes(),
	)

	eventsChan := make(chan model.LocationEvent, 32)

	ctx, _ := context.WithCancel(context.Background())

	c := NewDbConsumer(1, 1, time.Second, repo, eventsChan)
	c.Start(ctx)

	e := <-eventsChan
	if e != events[0] {
		t.Errorf("event received from channel is not equal to locked event")
	}

	c.Close()
}
