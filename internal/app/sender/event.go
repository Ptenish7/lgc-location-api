package sender

import "github.com/ozonmp/lgc-location-api/internal/model"

type EventSender interface {
	Send(location *model.LocationEvent) error
}
