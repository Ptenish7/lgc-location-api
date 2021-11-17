package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Location structure
type Location struct {
	ID        uint64    `db:"id"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	Title     string    `db:"title"`
	Removed   bool      `db:"removed"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// EventType type alias
type EventType uint8

// EventType constants
const (
	Created EventType = 1 + iota
	Updated
	Removed
)

var eventTypeNames = []string{"Created", "Updated", "Removed"}

// Scan converts EventType from DB enum value to Go value
func (t *EventType) Scan(src interface{}) error {
	for i, v := range eventTypeNames {
		if v == src.(string) {
			*t = EventType(i + 1)
			return nil
		}
	}
	return fmt.Errorf("EventType name not found: %s", src)
}

// Value converts EventType to DB enum value
func (t EventType) Value() (driver.Value, error) {
	return t.String(), nil
}

func (t EventType) String() string {
	return eventTypeNames[t-1]
}

// EventTypeExtra type alias
type EventTypeExtra uint16

// WithLatitude adds latitude bits
func (e EventTypeExtra) WithLatitude() EventTypeExtra {
	return e | EventTypeExtra(1)
}

// WithLongitude adds longitude bits
func (e EventTypeExtra) WithLongitude() EventTypeExtra {
	return e | EventTypeExtra(2)
}

// WithTitle adds title bits
func (e EventTypeExtra) WithTitle() EventTypeExtra {
	return e | EventTypeExtra(4)
}

// HasLatitude returns true if latitude bit is set to 1
func (e EventTypeExtra) HasLatitude() bool {
	return (e & EventTypeExtra(1)) == 1
}

// HasLongitude returns true if longitude bit is set to 1
func (e EventTypeExtra) HasLongitude() bool {
	return (e & EventTypeExtra(2)) == 2
}

// HasTitle returns true if title bit is set to 1
func (e EventTypeExtra) HasTitle() bool {
	return (e & EventTypeExtra(4)) == 4
}

// EventStatus type alias
type EventStatus uint8

// EventStatus constants
const (
	Deferred EventStatus = 1 + iota
	Processed
)

var eventStatusNames = []string{"Deferred", "Processed"}

// Scan converts EventStatus from DB enum value to Go value
func (s *EventStatus) Scan(src interface{}) error {
	for i, v := range eventStatusNames {
		if v == src.(string) {
			*s = EventStatus(i + 1)
			return nil
		}
	}
	return fmt.Errorf("EventStatus name not found: %s", src)
}

// Value converts EventStatus to DB enum value
func (s EventStatus) Value() (driver.Value, error) {
	return eventStatusNames[s-1], nil
}

// LocationEvent structure
type LocationEvent struct {
	ID         uint64         `db:"id"`
	LocationID uint64         `db:"location_id"`
	Type       EventType      `db:"type"`
	TypeExtra  EventTypeExtra `db:"type_extra"`
	Status     EventStatus    `db:"status"`
	Entity     *Location      `db:"payload"`
	UpdatedAt  time.Time      `db:"updated_at"`
}
