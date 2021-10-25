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
	n       uint64
	timeout time.Duration

	repo   repo.EventRepo
	sender sender.EventSender
	events <-chan model.LocationEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool
}

func NewKafkaProducer(
	n uint64,
	repo repo.EventRepo,
	sender sender.EventSender,
	events <-chan model.LocationEvent,
	workerPool *workerpool.WorkerPool,
) Producer {
	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &producer{
		n:          n,
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
	for {
		select {
		case event := <-p.events:
			if event.Type == model.Created {
				if err := p.sender.Send(&event); err != nil {
					log.Printf("failed to send event: %v", err)
					p.workerPool.Submit(func() {
						if err := p.repo.Unlock([]uint64{event.ID}); err != nil {
							log.Printf("failed to unlock events: %v", err)
						}
					})
				} else {
					p.workerPool.Submit(func() {
						if err := p.repo.Remove([]uint64{event.ID}); err != nil {
							log.Printf("failed to remove events: %v", err)
						}
					})
				}
			}

		case <-p.done:
			return
		}
	}
}
