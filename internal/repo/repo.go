package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/lgc-location-api/internal/model"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// Repo is DAO for Location
type Repo interface {
	CreateLocation(ctx context.Context, latitude float64, longitude float64, title string) (uint64, error)
	DescribeLocation(ctx context.Context, locationID uint64) (*model.Location, error)
	ListLocations(ctx context.Context, limit uint64, cursor uint64) ([]*model.Location, error)
	RemoveLocation(ctx context.Context, locationID uint64) (bool, error)
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
	query := psql.
		Insert("locations").
		Columns("latitude", "longitude", "title").
		Values(latitude, longitude, title).Suffix("RETURNING id").
		RunWith(r.db)

	var insertedID uint64
	err := query.QueryRowContext(ctx).Scan(&insertedID)

	return insertedID, err
}

// DescribeLocation returns a location by id
func (r *repo) DescribeLocation(ctx context.Context, locationID uint64) (*model.Location, error) {
	query := psql.
		Select("*").
		From("locations").
		Where(sq.Eq{"id": locationID, "removed": false})

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var result *model.Location
	err = r.db.GetContext(ctx, result, s, args)

	return result, err
}

// ListLocations returns all locations
func (r *repo) ListLocations(ctx context.Context, limit uint64, cursor uint64) ([]*model.Location, error) {
	query := psql.
		Select("*").
		From("locations").
		Where(sq.Eq{"removed": false}).
		Limit(limit).
		Offset(cursor)

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var result []*model.Location
	err = r.db.SelectContext(ctx, result, s, args)

	return result, err
}

// RemoveLocation removes a location by id
func (r *repo) RemoveLocation(ctx context.Context, locationID uint64) (bool, error) {
	query := psql.
		Update("locations").
		Set("removed", true).
		Set("updated_at", "now()").
		Where(sq.Eq{"id": locationID, "removed": false}).
		RunWith(r.db)

	result, err := query.ExecContext(ctx)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
