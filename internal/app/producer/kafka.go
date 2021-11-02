package producer

import (
	"context"
	"log"
	"sync"

	"github.com/gammazero/workerpool"

	"github.com/ozonmp/lgc-location-api/internal/app/repo"
	"github.com/ozonmp/lgc-location-api/internal/app/sender"
	"github.com/ozonmp/lgc-location-api/internal/model"
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

	repo   repo.EventRepo
	sender sender.EventSender
	events <-chan model.LocationEvent

	workerPool *workerpool.WorkerPool

	wg *sync.WaitGroup
}

func NewKafkaProducer(
	n uint64,
	batchSize uint64,
	repo repo.EventRepo,
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
			if event.Type == model.Created {
				if err := p.sender.Send(&event); err != nil {
					log.Printf("failed to send event: %v", err)
					updateBatch = append(updateBatch, event.ID)
					if len(updateBatch) == int(p.batchSize) {
						p.update(updateBatch)
						updateBatch = updateBatch[:0]
					}
				} else {
					cleanBatch = append(cleanBatch, event.ID)
					if len(cleanBatch) == int(p.batchSize) {
						p.clean(cleanBatch)
						cleanBatch = cleanBatch[:0]
					}
				}
			}

		case <-ctx.Done():
			if len(updateBatch) > 0 {
				p.update(updateBatch)
			}
			if len(cleanBatch) > 0 {
				p.clean(cleanBatch)
			}
			return
		}
	}
}

func (p *producer) update(eventIDs []uint64) {
	p.workerPool.Submit(func() {
		if err := p.repo.Unlock(eventIDs); err != nil {
			log.Printf("failed to unlock events: %v", err)
		}
	})
}

func (p *producer) clean(eventIDs []uint64) {
	p.workerPool.Submit(func() {
		if err := p.repo.Remove(eventIDs); err != nil {
			log.Printf("failed to remove events: %v", err)
		}
	})
}
