package retranslator

import (
	"time"

	"github.com/gammazero/workerpool"

	"github.com/ozonmp/lgc-location-api/internal/app/consumer"
	"github.com/ozonmp/lgc-location-api/internal/app/producer"
	"github.com/ozonmp/lgc-location-api/internal/app/repo"
	"github.com/ozonmp/lgc-location-api/internal/app/sender"
	"github.com/ozonmp/lgc-location-api/internal/model"
)

type Retranslator interface {
	Start()
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount   uint64
	ConsumerSize    uint64
	ConsumerTimeout time.Duration

	ProducerCount uint64
	WorkerCount   int
	BatchSize     uint64

	Repo   repo.EventRepo
	Sender sender.EventSender
}

type retranslator struct {
	events     chan model.LocationEvent
	consumer   consumer.Consumer
	producer   producer.Producer
	workerPool *workerpool.WorkerPool
}

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
	r.producer.Start()
	r.consumer.Start()
}

func (r *retranslator) Close() {
	r.consumer.Close()
	r.producer.Close()
	r.workerPool.StopWait()
}
