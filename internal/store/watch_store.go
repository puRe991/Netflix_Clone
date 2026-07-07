package store

import (
	"context"

	"github.com/pure991/streamflix/internal/models"
)

const progressColumns = `id, user_id, profile_id, media_id, episode_id, progress_seconds, completed, updated_at`

func scanProgress(row interface{ Scan(dest ...any) error }) (models.WatchProgress, error) {
	var p models.WatchProgress
	err := row.Scan(&p.ID, &p.UserID, &p.ProfileID, &p.MediaID, &p.EpisodeID, &p.ProgressSeconds, &p.Completed, &p.UpdatedAt)
	return p, err
}

// UpsertProgress writes watch progress for a profile+media(+episode) tuple.
// Race-free thanks to the partial unique indexes on watch_progress (see
// migration 0001): movies (episodeID == nil) and episodes each have their
// own conflict target.
func (s *Store) UpsertProgress(ctx context.Context, userID, profileID, mediaID string, episodeID *string, progressSeconds int, completed bool) (models.WatchProgress, error) {
	var row interface{ Scan(dest ...any) error }

	if episodeID == nil {
		row = s.Pool.QueryRow(ctx, `
			insert into watch_progress (user_id, profile_id, media_id, episode_id, progress_seconds, completed)
			values ($1, $2, $3, null, $4, $5)
			on conflict (profile_id, media_id) where episode_id is null
			do update set progress_seconds = excluded.progress_seconds, completed = excluded.completed, updated_at = now()
			returning `+progressColumns,
			userID, profileID, mediaID, progressSeconds, completed)
	} else {
		row = s.Pool.QueryRow(ctx, `
			insert into watch_progress (user_id, profile_id, media_id, episode_id, progress_seconds, completed)
			values ($1, $2, $3, $4, $5, $6)
			on conflict (profile_id, media_id, episode_id) where episode_id is not null
			do update set progress_seconds = excluded.progress_seconds, completed = excluded.completed, updated_at = now()
			returning `+progressColumns,
			userID, profileID, mediaID, *episodeID, progressSeconds, completed)
	}

	return scanProgress(row)
}

func (s *Store) GetMovieProgress(ctx context.Context, profileID, mediaID string) (models.WatchProgress, error) {
	row := s.Pool.QueryRow(ctx, `
		select `+progressColumns+` from watch_progress
		where profile_id = $1 and media_id = $2 and episode_id is null
	`, profileID, mediaID)
	p, err := scanProgress(row)
	return p, mapNotFound(err)
}

func (s *Store) GetEpisodeProgress(ctx context.Context, profileID, episodeID string) (models.WatchProgress, error) {
	row := s.Pool.QueryRow(ctx, `
		select `+progressColumns+` from watch_progress
		where profile_id = $1 and episode_id = $2
	`, profileID, episodeID)
	p, err := scanProgress(row)
	return p, mapNotFound(err)
}

func (s *Store) ListProgressForUserMedia(ctx context.Context, userID, mediaID string) ([]models.WatchProgress, error) {
	rows, err := s.Pool.Query(ctx, `
		select `+progressColumns+` from watch_progress where user_id = $1 and media_id = $2
	`, userID, mediaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.WatchProgress
	for rows.Next() {
		p, err := scanProgress(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (s *Store) SumWatchtimeSeconds(ctx context.Context) (int64, error) {
	var total *int64
	err := s.Pool.QueryRow(ctx, `select sum(progress_seconds) from watch_progress`).Scan(&total)
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}
	return *total, nil
}

// VerifyProfileOwnership confirms profileID belongs to userID (used before
// accepting a progress write, replacing the previous throw-on-mismatch
// findFirstOrThrow behavior with a clean checkable error).
func (s *Store) VerifyProfileOwnership(ctx context.Context, profileID, userID string) error {
	var exists bool
	err := s.Pool.QueryRow(ctx, `select exists(select 1 from profiles where id = $1 and user_id = $2)`, profileID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}
