package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/golang/mock/gomock"

	"github.com/ozonmp/lgc-location-api/internal/mocks"
	"github.com/ozonmp/lgc-location-api/internal/model"
)

func TestProducer(t *testing.T) {
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
	sender := mocks.NewMockEventSender(ctrl)

	eventsChan := make(chan model.LocationEvent, 32)
	workerPool := workerpool.New(2)

	ctx, cancel := context.WithCancel(context.Background())

	p := NewKafkaProducer(1, 1, repo, sender, eventsChan, workerPool)
	p.Start(ctx)

	t.Run("fail send and unlock", func(t *testing.T) {
		eventsChan <- events[0]

		gomock.InOrder(
			sender.EXPECT().Send(gomock.Any()).Return(errors.New("failed to send event")).Times(1),
			repo.EXPECT().Unlock(gomock.Any(), gomock.Any()).Times(1),
		)
	})

	t.Run("send and remove", func(t *testing.T) {
		eventsChan <- events[0]

		gomock.InOrder(
			sender.EXPECT().Send(gomock.Any()).Return(nil).Times(1),
			repo.EXPECT().Remove(gomock.Any(), gomock.Any()).Times(1),
		)
	})

	time.Sleep(time.Second)
	cancel()
	p.Close()
}
