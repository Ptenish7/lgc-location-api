package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/repo"
)

// Consumer interface
type Consumer interface {
	Start(ctx context.Context)
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.LocationEvent

	repo eventrepo.EventRepo

	batchSize uint64
	timeout   time.Duration

	wg *sync.WaitGroup
}

// NewDbConsumer creates a new db consumer
func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo eventrepo.EventRepo,
	events chan<- model.LocationEvent,
) Consumer {
	wg := &sync.WaitGroup{}

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
	}
}

func (c *consumer) Start(ctx context.Context) {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			c.consume(ctx)
		}()
	}
}

func (c *consumer) Close() {
	c.wg.Wait()
}

func (c *consumer) consume(ctx context.Context) {
	ticker := time.NewTicker(c.timeout)

	for {
		select {
		case <-ticker.C:
			events, err := c.repo.Lock(ctx, c.batchSize)
			if err != nil {
				logger.ErrorKV(ctx, "failed to lock events", "err", err)
				continue
			}

			logger.InfoKV(ctx, fmt.Sprintf("locked %d events", len(events)))
			for _, event := range events {
				c.events <- event
			}

		case <-ctx.Done():
			return
		}
	}
}
