package repo

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"

	"github.com/ozonmp/lgc-location-api/internal/metrics"
	"github.com/ozonmp/lgc-location-api/internal/model"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	ErrLocationNotFound = errors.New("location not found")
)

// Repo is DAO for Location
type Repo interface {
	CreateLocation(ctx context.Context, latitude float64, longitude float64, title string) (uint64, error)
	DescribeLocation(ctx context.Context, locationID uint64) (*model.Location, error)
	ListLocations(ctx context.Context, limit uint64, cursor uint64) ([]*model.Location, error)
	UpdateLocation(ctx context.Context, location *model.Location) error
	RemoveLocation(ctx context.Context, locationID uint64) error
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

// CreateLocation creates a new location
func (r *repo) CreateLocation(ctx context.Context, latitude float64, longitude float64, title string) (uint64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "repo.CreateLocation")
	defer span.Finish()

	query := psql.
		Insert("locations").
		Columns("latitude", "longitude", "title").
		Values(latitude, longitude, title).Suffix("RETURNING id").
		RunWith(r.db)

	var insertedID uint64
	err := query.QueryRowContext(ctx).Scan(&insertedID)

	if err == nil {
		metrics.IncEventCUDCounter(model.Created)
	}

	return insertedID, err
}

// DescribeLocation returns a location by id
func (r *repo) DescribeLocation(ctx context.Context, locationID uint64) (*model.Location, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "repo.DescribeLocation",
		opentracing.Tag{
			Key:   "locationID",
			Value: locationID,
		},
	)
	defer span.Finish()

	query := psql.
		Select("*").
		From("locations").
		Where(sq.Eq{"id": locationID, "removed": false})

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var result model.Location
	err = r.db.GetContext(ctx, &result, s, args...)

	return &result, err
}

// ListLocations returns all locations
func (r *repo) ListLocations(ctx context.Context, limit uint64, cursor uint64) ([]*model.Location, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "repo.ListLocations",
		opentracing.Tag{
			Key:   "limit",
			Value: limit,
		},
		opentracing.Tag{
			Key:   "cursor",
			Value: cursor,
		},
	)
	defer span.Finish()

	query := psql.
		Select("*").
		From("locations").
		Where(sq.Eq{"removed": false}).
		OrderBy("id").
		Limit(limit).
		Offset(cursor)

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var result []*model.Location
	err = r.db.SelectContext(ctx, &result, s, args...)

	return result, err
}

// UpdateLocation updates specified location
func (r *repo) UpdateLocation(ctx context.Context, location *model.Location) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "repo.RemoveLocation",
		opentracing.Tag{
			Key:   "locationID",
			Value: location.ID,
		},
	)
	defer span.Finish()

	query := psql.
		Update("locations").
		Set("latitude", location.Latitude).
		Set("longitude", location.Longitude).
		Set("title", location.Title).
		Set("updated_at", "now()").
		Where(sq.Eq{"id": location.ID, "removed": false}).
		RunWith(r.db)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrLocationNotFound
	}

	metrics.IncEventCUDCounter(model.Updated)

	return nil
}

// RemoveLocation removes a location by id
func (r *repo) RemoveLocation(ctx context.Context, locationID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "repo.RemoveLocation",
		opentracing.Tag{
			Key:   "locationID",
			Value: locationID,
		},
	)
	defer span.Finish()

	query := psql.
		Update("locations").
		Set("removed", true).
		Set("updated_at", "now()").
		Where(sq.Eq{"id": locationID, "removed": false}).
		RunWith(r.db)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrLocationNotFound
	}

	metrics.IncEventCUDCounter(model.Removed)

	return nil
}
