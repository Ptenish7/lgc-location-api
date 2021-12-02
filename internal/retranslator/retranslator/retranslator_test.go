package retranslator

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/ozonmp/lgc-location-api/internal/mocks"
	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/config"
)

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := config.Retranslator{
		ChannelSize:     512,
		ConsumerCount:   2,
		ConsumerSize:    10,
		ConsumerTimeout: 10 * time.Second,
		ProducerCount:   2,
		WorkerCount:     2,
	}

	retranslator := NewRetranslator(&cfg, repo, sender)
	retranslator.Start()
	retranslator.Close()
}

func TestLockRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	events := []model.LocationEvent{
		{
			ID:     1,
			Type:   model.Created,
			Status: model.Deferred,
			Entity: &model.Location{ID: 1, Title: "L1"},
		},
		{
			ID:     2,
			Type:   model.Created,
			Status: model.Deferred,
			Entity: &model.Location{ID: 2, Title: "L2"},
		},
	}

	lockAll := repo.EXPECT().Lock(gomock.Any(), uint64(2)).Return(events, nil).Times(1)
	repo.EXPECT().Lock(gomock.Any(), gomock.Any()).Return(nil, errors.New("no events to lock")).AnyTimes().After(lockAll)

	sentIDs := make([]uint64, 0)
	sender.EXPECT().Send(gomock.Any()).Times(2).DoAndReturn(func(event *model.LocationEvent) error {
		sentIDs = append(sentIDs, event.ID)
		return nil
	})

	removedIDs := make([]uint64, 0)
	repo.EXPECT().Remove(gomock.Any(), gomock.Any()).Times(2).DoAndReturn(func(eventIDs []uint64) error {
		removedIDs = append(removedIDs, eventIDs...)
		return nil
	})

	if !gomock.InAnyOrder(sentIDs).Matches(removedIDs) {
		t.Errorf("sent and removed IDs not matched")
	}

	cfg := config.Retranslator{
		ChannelSize:     512,
		ConsumerCount:   2,
		ConsumerSize:    10,
		ConsumerTimeout: 10 * time.Second,
		ProducerCount:   2,
		WorkerCount:     2,
	}

	retranslator := NewRetranslator(&cfg, repo, sender)
	retranslator.Start()
	time.Sleep(2 * time.Second)
	retranslator.Close()
}

func TestLockUnlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	events := []model.LocationEvent{
		{
			ID:     1,
			Type:   model.Created,
			Status: model.Deferred,
			Entity: &model.Location{ID: 1, Title: "L1"},
		},
		{
			ID:     2,
			Type:   model.Created,
			Status: model.Deferred,
			Entity: &model.Location{ID: 2, Title: "L2"},
		},
	}

	lockAll := repo.EXPECT().Lock(gomock.Any(), uint64(2)).Return(events, nil).Times(1)
	repo.EXPECT().Lock(gomock.Any(), gomock.Any()).Return(nil, errors.New("no events to lock")).AnyTimes().After(lockAll)

	notSentIDs := make([]uint64, 0)
	sender.EXPECT().Send(gomock.Any()).Times(2).DoAndReturn(func(event *model.LocationEvent) error {
		notSentIDs = append(notSentIDs, event.ID)
		return errors.New("failed to send event")
	})

	unlockedIDs := make([]uint64, 0)
	repo.EXPECT().Unlock(gomock.Any(), gomock.Any()).Times(2).DoAndReturn(func(eventIDs []uint64) error {
		unlockedIDs = append(unlockedIDs, eventIDs...)
		return nil
	})

	if !gomock.InAnyOrder(notSentIDs).Matches(unlockedIDs) {
		t.Errorf("failed to sent and unlocked IDs not matched")
	}

	cfg := config.Retranslator{
		ChannelSize:     512,
		ConsumerCount:   2,
		ConsumerSize:    10,
		ConsumerTimeout: 10 * time.Second,
		ProducerCount:   2,
		WorkerCount:     2,
	}

	retranslator := NewRetranslator(&cfg, repo, sender)
	retranslator.Start()
	time.Sleep(2 * time.Second)
	retranslator.Close()
}
