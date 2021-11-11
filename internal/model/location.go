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
	return eventTypeNames[t-1], nil
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
	ID         uint64      `db:"id"`
	LocationID uint64      `db:"location_id"`
	Type       EventType   `db:"type"`
	Status     EventStatus `db:"status"`
	Entity     *Location   `db:"payload"`
	UpdatedAt  time.Time   `db:"updated_at"`
}
