// Package store contains hand-written SQL access for every aggregate,
// using pgx directly (no ORM).
package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pure991/streamflix/internal/models"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	Pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{Pool: pool}
}

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *Store) GetStats(ctx context.Context) (models.Stats, error) {
	var stats models.Stats
	var err error

	if stats.Users, err = s.CountUsers(ctx); err != nil {
		return stats, err
	}
	if stats.ActiveSubscriptions, err = s.CountActiveSubscriptions(ctx); err != nil {
		return stats, err
	}
	if stats.Media, err = s.CountMedia(ctx); err != nil {
		return stats, err
	}
	if stats.WatchtimeSeconds, err = s.SumWatchtimeSeconds(ctx); err != nil {
		return stats, err
	}

	return stats, nil
}
