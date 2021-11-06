package model

// Location structure
type Location struct {
	ID        uint64
	Latitude  float64
	Longitude float64
	Title     string
}

// EventType type alias
type EventType uint8

// EventStatus type alias
type EventStatus uint8

// EventType and EventStatus constants
const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

// LocationEvent structure
type LocationEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Location
}
