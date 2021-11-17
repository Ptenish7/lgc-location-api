package eventrepo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/ozonmp/lgc-location-api/internal/metrics"
	"github.com/ozonmp/lgc-location-api/internal/model"
	pb "github.com/ozonmp/lgc-location-api/pkg/lgc-location-api"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

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

	var result []model.LocationEvent
	err = r.db.SelectContext(ctx, &result, s, args...)

	if err != nil {
		metrics.AddEventsInRetranslator(len(result))
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
