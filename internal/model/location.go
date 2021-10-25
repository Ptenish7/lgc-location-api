package model

type Location struct {
	ID        uint64
	Latitude  float64
	Longitude float64
	Title     string
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type LocationEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Location
}
