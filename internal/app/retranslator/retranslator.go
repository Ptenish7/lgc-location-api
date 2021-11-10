package retranslator

import (
	"context"
	"time"

	"github.com/gammazero/workerpool"

	"github.com/ozonmp/lgc-location-api/internal/app/consumer"
	"github.com/ozonmp/lgc-location-api/internal/app/producer"
	"github.com/ozonmp/lgc-location-api/internal/app/repo"
	"github.com/ozonmp/lgc-location-api/internal/app/sender"
	"github.com/ozonmp/lgc-location-api/internal/model"
)

// Retranslator interface
type Retranslator interface {
	Start()
	Close()
}

// Config represents retranslator config
type Config struct {
	ChannelSize uint64

	ConsumerCount   uint64
	ConsumerSize    uint64
	ConsumerTimeout time.Duration

	ProducerCount uint64
	WorkerCount   int
	BatchSize     uint64

	Repo   eventrepo.EventRepo
	Sender sender.EventSender
}

type retranslator struct {
	events     chan model.LocationEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
	cancelFunc context.CancelFunc
}

// NewRetranslator creates new retranslator
func NewRetranslator(cfg Config) Retranslator {
	events := make(chan model.LocationEvent, cfg.ChannelSize)
	workerPool := workerpool.New(cfg.WorkerCount)

	c := consumer.NewDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumerSize,
		cfg.ConsumerTimeout,
		cfg.Repo,
		events,
	)

	p := producer.NewKafkaProducer(
		cfg.ProducerCount,
		cfg.BatchSize,
		cfg.Repo,
		cfg.Sender,
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
