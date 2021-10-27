package consumer

import (
	"log"
	"sync"
	"time"

	"github.com/ozonmp/lgc-location-api/internal/app/repo"
	"github.com/ozonmp/lgc-location-api/internal/model"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.LocationEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	done chan bool
	wg   *sync.WaitGroup
}

type Config struct {
	n         uint64
	events    chan<- model.LocationEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.LocationEvent,
) Consumer {
	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		done:      done,
	}
}

func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			c.consume()
		}()
	}
}

func (c *consumer) Close() {
	close(c.done)
	c.wg.Wait()
}

func (c *consumer) consume() {
	ticker := time.NewTicker(c.timeout)

	for {
		select {
		case <-ticker.C:
			events, err := c.repo.Lock(c.batchSize)
			if err != nil {
				log.Printf("failed to lock events: %v", err)
				continue
			}

			log.Printf("locked %d events", len(events))
			for _, event := range events {
				c.events <- event
			}

		case <-c.done:
			return
		}
	}
}
