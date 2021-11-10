package eventrepo

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/encoding/protojson"

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
	err = r.db.SelectContext(ctx, result, s, args)

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

	return nil
}

// Add adds new event
func (r *eventRepo) Add(ctx context.Context, event *model.LocationEvent) error {
	payloadPb := &pb.Location{
		Id:        event.Entity.ID,
		Latitude:  event.Entity.Latitude,
		Longitude: event.Entity.Longitude,
		Title:     event.Entity.Title,
	}

	payload, err := protojson.Marshal(payloadPb)
	if err != nil {
		return err
	}

	query := psql.
		Insert("locations_events").
		Columns("location_id", "type", "status", "payload").
		Values(event.Entity.ID, event.Type, event.Status, payload).
		RunWith(r.db)

	_, err = query.ExecContext(ctx)

	return nil
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

	return nil
}
