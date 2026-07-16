// Package db manages the Postgres connection pool and schema migrations.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// connectTimeout bounds how long Connect waits for the initial ping so a
// missing/unreachable database fails fast with a clear error instead of
// hanging indefinitely (e.g. when PostgreSQL was never installed/started).
const connectTimeout = 10 * time.Second

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database (is PostgreSQL running and reachable at the DATABASE_URL in .env?): %w", err)
	}

	return pool, nil
}
