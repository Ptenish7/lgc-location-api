package producer

import (
	"log"
	"sync"
	"time"

	"github.com/gammazero/workerpool"

	"github.com/ozonmp/lgc-location-api/internal/app/repo"
	"github.com/ozonmp/lgc-location-api/internal/app/sender"
	"github.com/ozonmp/lgc-location-api/internal/model"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n         uint64
	timeout   time.Duration
	batchSize uint64

	repo   repo.EventRepo
	sender sender.EventSender
	events <-chan model.LocationEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool
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
	done := make(chan bool)

	return &producer{
		n:          n,
		batchSize:  batchSize,
		repo:       repo,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
	}
}

func (p *producer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			p.produce()
		}()
	}
}

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
}

func (p *producer) produce() {
	updateBatch := make([]uint64, 0, p.batchSize)
	cleanBatch := make([]uint64, 0, p.batchSize)

	for {
		select {
		case event := <-p.events:
			if event.Type == model.Created {
				if err := p.sender.Send(&event); err != nil {
					log.Printf("failed to send event: %v", err)
					updateBatch = append(updateBatch, event.ID)
					if len(updateBatch) == cap(updateBatch) {
						p.update(updateBatch)
						updateBatch = updateBatch[:0]
					}
				} else {
					cleanBatch = append(cleanBatch, event.ID)
					if len(cleanBatch) == cap(cleanBatch) {
						p.clean(cleanBatch)
						cleanBatch = cleanBatch[:0]
					}
				}
			}

		case <-p.done:
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
