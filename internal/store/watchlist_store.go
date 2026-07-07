package store

import (
	"context"

	"github.com/pure991/streamflix/internal/models"
)

func (s *Store) ListWatchlist(ctx context.Context, userID string) ([]models.WatchlistItem, error) {
	rows, err := s.Pool.Query(ctx, `
		select w.id, w.user_id, w.profile_id, w.media_id, w.created_at, `+mediaColumnsQualified+`
		from watchlist w
		join media m on m.id = w.media_id
		where w.user_id = $1
		order by w.created_at desc
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.WatchlistItem
	for rows.Next() {
		var it models.WatchlistItem
		var m models.Media
		if err := rows.Scan(&it.ID, &it.UserID, &it.ProfileID, &it.MediaID, &it.CreatedAt,
			&m.ID, &m.Title, &m.Slug, &m.Description, &m.Type, &m.ThumbnailURL, &m.BannerURL,
			&m.TrailerURL, &m.VideoURL, &m.ReleaseYear, &m.Duration, &m.AgeRating, &m.Language,
			&m.Tags, &m.IsPublished, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		it.Media = &m
		items = append(items, it)
	}
	return items, rows.Err()
}

// AddToWatchlist is idempotent (matches the previous upsert-on-conflict-do-nothing behavior).
func (s *Store) AddToWatchlist(ctx context.Context, userID, profileID, mediaID string) error {
	_, err := s.Pool.Exec(ctx, `
		insert into watchlist (user_id, profile_id, media_id)
		values ($1, $2, $3)
		on conflict (profile_id, media_id) do nothing
	`, userID, profileID, mediaID)
	return err
}

func (s *Store) RemoveFromWatchlist(ctx context.Context, id, userID string) error {
	_, err := s.Pool.Exec(ctx, `delete from watchlist where id = $1 and user_id = $2`, id, userID)
	return err
}
