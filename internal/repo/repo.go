package repo

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/lgc-location-api/internal/model"
)

// ErrNotImplemented is returned if repo method is not yet implemented
var ErrNotImplemented = errors.New("repo method not implemented")

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
	return 0, ErrNotImplemented
}

// DescribeLocation returns a location by id
func (r *repo) DescribeLocation(ctx context.Context, locationID uint64) (*model.Location, error) {
	return nil, ErrNotImplemented
}

// ListLocations returns all locations
func (r *repo) ListLocations(ctx context.Context) ([]*model.Location, error) {
	return nil, ErrNotImplemented
}

// RemoveLocation removes a location by id
func (r *repo) RemoveLocation(ctx context.Context, locationID uint64) (bool, error) {
	return false, ErrNotImplemented
}
