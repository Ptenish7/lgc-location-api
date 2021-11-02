package model

// Location
type Location struct {
	ID        uint64
	Latitude  float64
	Longitude float64
	Title     string
}

// EventType
type EventType uint8

// EventStatus
type EventStatus uint8

// EventType and EventStatus constants
const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

// LocationEvent
type LocationEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Location
}
