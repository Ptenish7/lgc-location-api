package sender

import (
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ozonmp/lgc-location-api/internal/model"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

var topics = []string{"create_location", "update_location", "remove_location"}

// EventSender interface
type EventSender interface {
	Send(event *model.LocationEvent) error
}

type eventSender struct {
	producer sarama.SyncProducer
}

// NewEventSender created new event sender with specified brokers
func NewEventSender(brokers []string, maxRetry uint64) (EventSender, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = int(maxRetry)

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &eventSender{producer: producer}, nil
}

func (s *eventSender) Send(event *model.LocationEvent) error {
	eventBytes, err := proto.Marshal(convertLocationEventToPb(*event))
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: getTopic(event.Type),
		Value: sarama.ByteEncoder(eventBytes),
	}

	_, _, err = s.producer.SendMessage(msg)
	return err
}

func getTopic(eventType model.EventType) string {
	return topics[int(eventType)-1]
}

func convertLocationEventToPb(e model.LocationEvent) *pb.LocationEvent {
	return &pb.LocationEvent{
		Id:         e.ID,
		LocationId: e.LocationID,
		Type:       uint32(e.Type),
		ExtraType:  uint32(e.TypeExtra),
		Status:     uint32(e.Status),
		Entity:     convertLocationToPb(e.Entity),
		UpdatedAt:  timestamppb.New(e.UpdatedAt),
	}
}

func convertLocationToPb(l *model.Location) *pb.Location {
	return &pb.Location{
		Id:        l.ID,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Title:     l.Title,
	}
}
