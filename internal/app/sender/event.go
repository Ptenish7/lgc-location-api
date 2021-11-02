package sender

import "github.com/ozonmp/lgc-location-api/internal/model"

// EventSender
type EventSender interface {
	Send(location *model.LocationEvent) error
}
