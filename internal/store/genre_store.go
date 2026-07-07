package store

import (
	"context"

	"github.com/pure991/streamflix/internal/models"
)

func (s *Store) ListGenres(ctx context.Context) ([]models.Genre, error) {
	rows, err := s.Pool.Query(ctx, `select id, name, slug from genres order by name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name, &g.Slug); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, rows.Err()
}

func (s *Store) CreateGenre(ctx context.Context, name, slug string) (models.Genre, error) {
	var g models.Genre
	err := s.Pool.QueryRow(ctx, `
		insert into genres (name, slug) values ($1, $2)
		returning id, name, slug
	`, name, slug).Scan(&g.ID, &g.Name, &g.Slug)
	return g, err
}

func (s *Store) DeleteGenre(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, `delete from genres where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
