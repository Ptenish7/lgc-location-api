package sender

import "github.com/ozonmp/lgc-location-api/internal/model"

// EventSender interface
type EventSender interface {
	Send(location *model.LocationEvent) error
}
