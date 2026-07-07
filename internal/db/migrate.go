package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrate applies any embedded .sql migration files not yet recorded in
// schema_migrations, in filename order. Each migration runs in its own
// transaction.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `create table if not exists schema_migrations (
		version text primary key,
		applied_at timestamptz not null default now()
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		var applied bool
		if err := pool.QueryRow(ctx, `select exists(select 1 from schema_migrations where version = $1)`, name).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		if applied {
			continue
		}

		sqlBytes, err := fs.ReadFile(migrationFiles, "migrations/"+name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", name, err)
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("apply migration %s: %w", name, err)
		}

		if _, err := tx.Exec(ctx, `insert into schema_migrations (version) values ($1)`, name); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("record migration %s: %w", name, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}
	}

	return nil
}
