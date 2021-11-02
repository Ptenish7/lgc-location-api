package repo

import "github.com/ozonmp/lgc-location-api/internal/model"

// EventRepo
type EventRepo interface {
	Lock(n uint64) ([]model.LocationEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.LocationEvent) error
	Remove(eventIDs []uint64) error
}
