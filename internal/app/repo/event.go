package eventrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/ozonmp/lgc-location-api/internal/metrics"
	"github.com/ozonmp/lgc-location-api/internal/model"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type locationEvent struct {
	ID         uint64               `db:"id"`
	LocationID uint64               `db:"location_id"`
	Type       model.EventType      `db:"type"`
	TypeExtra  model.EventTypeExtra `db:"type_extra"`
	Status     model.EventStatus    `db:"status"`
	Entity     locationPayload      `db:"payload"`
	UpdatedAt  time.Time            `db:"updated_at"`
}

type locationPayload pb.Location

func (p *locationPayload) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("incompatible type for locationPayload")
	}

	pl := &pb.Location{}

	err := protojson.Unmarshal(source, pl)
	if err != nil {
		return err
	}

	p.Id = pl.Id
	p.Latitude = pl.Latitude
	p.Longitude = pl.Longitude
	p.Title = pl.Title

	return nil
}

func convertProtobufToLocation(pb *locationEvent) model.LocationEvent {
	return model.LocationEvent{
		ID:         pb.ID,
		LocationID: pb.LocationID,
		Type:       pb.Type,
		TypeExtra:  pb.TypeExtra,
		Status:     pb.Status,
		Entity: &model.Location{
			ID:        pb.Entity.Id,
			Latitude:  pb.Entity.Latitude,
			Longitude: pb.Entity.Longitude,
			Title:     pb.Entity.Title,
		},
		UpdatedAt: pb.UpdatedAt,
	}
}

// EventRepo interface
type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]model.LocationEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error

	Add(ctx context.Context, event *model.LocationEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}

type eventRepo struct {
	db *sqlx.DB
}

// NewEventRepo creates a new event repo
func NewEventRepo(db *sqlx.DB) EventRepo {
	return &eventRepo{db: db}
}

// Lock locks n records
func (r *eventRepo) Lock(ctx context.Context, n uint64) ([]model.LocationEvent, error) {
	subquery := psql.
		Select("id").
		From("locations_events").
		Where(sq.Eq{"status": model.Deferred}).
		Limit(n)

	query := psql.
		Update("locations_events").
		Set("status", model.Processed).
		Set("updated_at", "now()").
		Where(subquery.Prefix("id IN (").Suffix(")")).
		Suffix("RETURNING *")

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var resultPb []locationEvent
	err = r.db.SelectContext(ctx, &resultPb, s, args...)
	if err != nil {
		metrics.AddEventsInRetranslator(len(resultPb))
	}

	result := make([]model.LocationEvent, 0, len(resultPb))
	for _, v := range resultPb {
		result = append(result, convertProtobufToLocation(&v))
	}

	return result, err
}

// Unlock unlocks specified events
func (r *eventRepo) Unlock(ctx context.Context, eventIDs []uint64) error {
	query := psql.
		Update("locations_events").
		Set("status", model.Deferred).
		Set("updated_at", "now()").
		Where(sq.Eq{"id": eventIDs}).
		RunWith(r.db)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(eventIDs)) {
		return fmt.Errorf("%d rows were updated; want %d", rowsAffected, len(eventIDs))
	}

	metrics.SubtractEventsInRetranslator(int(rowsAffected))

	return nil
}

// Add adds new event
func (r *eventRepo) Add(ctx context.Context, event *model.LocationEvent) error {
	payload := []byte("{}")
	if event.Entity != nil {
		payloadPb := &pb.Location{
			Id: event.LocationID,
		}

		if event.Type == model.Created {
			payloadPb.Latitude = event.Entity.Latitude
			payloadPb.Longitude = event.Entity.Longitude
			payloadPb.Title = event.Entity.Title
		} else if event.Type == model.Updated {
			if event.TypeExtra.HasLatitude() {
				payloadPb.Latitude = event.Entity.Latitude
			}
			if event.TypeExtra.HasLongitude() {
				payloadPb.Longitude = event.Entity.Longitude
			}
			if event.TypeExtra.HasTitle() {
				payloadPb.Title = event.Entity.Title
			}
		}

		var err error
		payload, err = protojson.Marshal(payloadPb)
		if err != nil {
			return err
		}
	}

	query := psql.
		Insert("locations_events").
		Columns("location_id", "type", "type_extra", "status", "payload").
		Values(event.LocationID, event.Type, event.TypeExtra, event.Status, payload).
		RunWith(r.db)

	_, err := query.ExecContext(ctx)

	return err
}

// Remove removes specified events
func (r *eventRepo) Remove(ctx context.Context, eventIDs []uint64) error {
	query := psql.
		Delete("locations_events").
		Where(sq.Eq{"id": eventIDs}).
		RunWith(r.db)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(eventIDs)) {
		return fmt.Errorf("%d rows were deleted; want %d", rowsAffected, len(eventIDs))
	}

	metrics.SubtractEventsInRetranslator(int(rowsAffected))

	return nil
}
