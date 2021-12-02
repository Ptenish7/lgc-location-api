package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/lgc-location-api/internal/pkg/logger"
)

// NewPostgres returns DB
func NewPostgres(ctx context.Context, dsn, driver string, maxRetry uint64) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		logger.ErrorKV(ctx, "failed to create database connection", "err", err)

		return nil, err
	}

	var retries uint64 = 0

	if err = db.Ping(); err != nil {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if retries >= maxRetry {
				logger.ErrorKV(ctx, "failed to ping the database", "err", err)

				return nil, err
			}

			retries++

			logger.InfoKV(ctx, "failed to ping the database, retrying...")

			if err = db.Ping(); err == nil {
				return db, nil
			}
		}
	}

	return db, nil
}
