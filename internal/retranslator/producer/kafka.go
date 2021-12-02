package producer

import (
	"context"
	"sync"

	"github.com/gammazero/workerpool"

	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/repo"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/sender"
)

// Producer interface
type Producer interface {
	Start(ctx context.Context)
	Close()
}

type producer struct {
	n uint64
	//timeout   time.Duration
	batchSize uint64

	repo   eventrepo.EventRepo
	sender sender.EventSender
	events <-chan model.LocationEvent

	workerPool *workerpool.WorkerPool

	wg *sync.WaitGroup
}

// NewKafkaProducer creates new kafka producer
func NewKafkaProducer(
	n uint64,
	batchSize uint64,
	repo eventrepo.EventRepo,
	sender sender.EventSender,
	events <-chan model.LocationEvent,
	workerPool *workerpool.WorkerPool,
) Producer {
	wg := &sync.WaitGroup{}

	return &producer{
		n:          n,
		batchSize:  batchSize,
		repo:       repo,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
	}
}

func (p *producer) Start(ctx context.Context) {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.produce(ctx)
		}()
	}
}

func (p *producer) Close() {
	p.wg.Wait()
}

func (p *producer) produce(ctx context.Context) {
	updateBatch := make([]uint64, 0, p.batchSize)
	cleanBatch := make([]uint64, 0, p.batchSize)

	for {
		select {
		case event := <-p.events:
			if err := p.sender.Send(&event); err != nil {
				logger.ErrorKV(ctx, "failed to send event", "err", err)
				updateBatch = append(updateBatch, event.ID)
				if len(updateBatch) == int(p.batchSize) {
					p.update(ctx, updateBatch)
					updateBatch = updateBatch[:0]
				}
			} else {
				cleanBatch = append(cleanBatch, event.ID)
				if len(cleanBatch) == int(p.batchSize) {
					p.clean(ctx, cleanBatch)
					cleanBatch = cleanBatch[:0]
				}
			}

		case <-ctx.Done():
			if len(updateBatch) > 0 {
				p.update(ctx, updateBatch)
			}
			if len(cleanBatch) > 0 {
				p.clean(ctx, cleanBatch)
			}
			return
		}
	}
}

func (p *producer) update(ctx context.Context, eventIDs []uint64) {
	p.workerPool.Submit(func() {
		if err := p.repo.Unlock(ctx, eventIDs); err != nil {
			logger.ErrorKV(ctx, "failed to unlock events", "err", err)
		}
	})
}

func (p *producer) clean(ctx context.Context, eventIDs []uint64) {
	p.workerPool.Submit(func() {
		if err := p.repo.Remove(ctx, eventIDs); err != nil {
			logger.ErrorKV(ctx, "failed to remove events", "err", err)
		}
	})
}
