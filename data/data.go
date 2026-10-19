// Package data opens listo's Postgres; queries are generated into
// data/sqlcgen.
package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Open returns a pool for url, or for libpq's PG* variables when it's empty.
//
// NOTE: the pool connects lazily, a database that's down only shows up on the
// first query or Ping.
func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	return pool, nil
}
