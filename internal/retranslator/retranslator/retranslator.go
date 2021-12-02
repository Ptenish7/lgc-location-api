package retranslator

import (
	"context"

	"github.com/gammazero/workerpool"

	"github.com/ozonmp/lgc-location-api/internal/model"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/config"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/consumer"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/producer"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/repo"
	"github.com/ozonmp/lgc-location-api/internal/retranslator/sender"
)

// Retranslator interface
type Retranslator interface {
	Start()
	Close()
}

type retranslator struct {
	events     chan model.LocationEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
	cancelFunc context.CancelFunc
}

// NewRetranslator creates new retranslator
func NewRetranslator(cfg *config.Retranslator, repo eventrepo.EventRepo, sender sender.EventSender) Retranslator {
	events := make(chan model.LocationEvent, cfg.ChannelSize)
	workerPool := workerpool.New(int(cfg.WorkerCount))

	c := consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumerSize,
		cfg.ConsumerTimeout,
		repo,
		events,
	)

	p := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.BatchSize,
		repo,
		sender,
		events,
		workerPool,
	)

	return &retranslator{
		events:     events,
		consumer:   c,
		producer:   p,
		workerPool: workerPool,
	}
}

func (r *retranslator) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFunc = cancel
	r.producer.Start(ctx)
	r.consumer.Start(ctx)
}

func (r *retranslator) Close() {
	r.cancelFunc()
	r.consumer.Close()
	r.producer.Close()
	r.workerPool.StopWait()
}
