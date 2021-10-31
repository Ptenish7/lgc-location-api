package repo

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/lgc-location-api/internal/model"
)

// Repo is DAO for Location
type Repo interface {
	CreateLocation(ctx context.Context, latitude float64, longitude float64, title string) (uint64, error)
	DescribeLocation(ctx context.Context, locationID uint64) (*model.Location, error)
	ListLocations(ctx context.Context) ([]*model.Location, error)
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
	panic("implement me")
}

// DescribeLocation returns a location by id
func (r *repo) DescribeLocation(ctx context.Context, locationID uint64) (*model.Location, error) {
	panic("implement me")
}

// ListLocations returns all locations
func (r *repo) ListLocations(ctx context.Context) ([]*model.Location, error) {
	panic("implement me")
}

// RemoveLocation removes a location by id
func (r *repo) RemoveLocation(ctx context.Context, locationID uint64) (bool, error) {
	panic("implement me")
}
