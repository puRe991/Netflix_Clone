package store

import (
	"context"

	"github.com/pure991/streamflix/internal/models"
)

func (s *Store) GetSeriesByMediaID(ctx context.Context, mediaID string) (models.Series, error) {
	var series models.Series
	err := s.Pool.QueryRow(ctx, `select id, media_id, title, description from series where media_id = $1`, mediaID).
		Scan(&series.ID, &series.MediaID, &series.Title, &series.Description)
	if err != nil {
		return series, mapNotFound(err)
	}

	seasons, err := s.ListSeasons(ctx, series.ID)
	if err != nil {
		return series, err
	}
	series.Seasons = seasons
	return series, nil
}

func (s *Store) ListAllSeries(ctx context.Context) ([]models.Series, error) {
	rows, err := s.Pool.Query(ctx, `
		select se.id, se.media_id, se.title, se.description, m.title, m.slug,
			(select count(*) from seasons sn where sn.series_id = se.id)
		from series se
		join media m on m.id = se.media_id
		order by se.title asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Series
	for rows.Next() {
		var se models.Series
		var mediaTitle, mediaSlug string
		var seasonCount int
		if err := rows.Scan(&se.ID, &se.MediaID, &se.Title, &se.Description, &mediaTitle, &mediaSlug, &seasonCount); err != nil {
			return nil, err
		}
		se.Media = &models.Media{ID: se.MediaID, Title: mediaTitle, Slug: mediaSlug}
		se.Seasons = make([]models.Season, seasonCount)
		list = append(list, se)
	}
	return list, rows.Err()
}

// ListSeriesMediaWithoutSeries returns published-or-not Media rows of type
// SERIES that don't have a Series row yet, for the "create series" picker.
func (s *Store) ListSeriesMediaWithoutSeries(ctx context.Context) ([]models.Media, error) {
	rows, err := s.Pool.Query(ctx, `
		select `+mediaColumns+` from media m
		where m.type = 'SERIES' and not exists (select 1 from series se where se.media_id = m.id)
		order by m.created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Media
	for rows.Next() {
		m, err := scanMedia(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (s *Store) GetSeriesByID(ctx context.Context, id string) (models.Series, error) {
	var series models.Series
	err := s.Pool.QueryRow(ctx, `select id, media_id, title, description from series where id = $1`, id).
		Scan(&series.ID, &series.MediaID, &series.Title, &series.Description)
	return series, mapNotFound(err)
}

func (s *Store) CreateSeries(ctx context.Context, mediaID, title, description string) (models.Series, error) {
	var series models.Series
	err := s.Pool.QueryRow(ctx, `
		insert into series (media_id, title, description) values ($1, $2, $3)
		returning id, media_id, title, description
	`, mediaID, title, description).Scan(&series.ID, &series.MediaID, &series.Title, &series.Description)
	return series, err
}

func (s *Store) UpdateSeries(ctx context.Context, id, title, description string) (models.Series, error) {
	var series models.Series
	err := s.Pool.QueryRow(ctx, `
		update series set title = $1, description = $2 where id = $3
		returning id, media_id, title, description
	`, title, description, id).Scan(&series.ID, &series.MediaID, &series.Title, &series.Description)
	return series, mapNotFound(err)
}

func (s *Store) DeleteSeries(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, `delete from series where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) ListSeasons(ctx context.Context, seriesID string) ([]models.Season, error) {
	rows, err := s.Pool.Query(ctx, `
		select id, series_id, season_number, title from seasons where series_id = $1 order by season_number asc
	`, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seasons []models.Season
	for rows.Next() {
		var se models.Season
		if err := rows.Scan(&se.ID, &se.SeriesID, &se.SeasonNumber, &se.Title); err != nil {
			return nil, err
		}
		episodes, err := s.ListEpisodes(ctx, se.ID)
		if err != nil {
			return nil, err
		}
		se.Episodes = episodes
		seasons = append(seasons, se)
	}
	return seasons, rows.Err()
}

func (s *Store) GetSeasonByID(ctx context.Context, id string) (models.Season, error) {
	var se models.Season
	err := s.Pool.QueryRow(ctx, `select id, series_id, season_number, title from seasons where id = $1`, id).
		Scan(&se.ID, &se.SeriesID, &se.SeasonNumber, &se.Title)
	return se, mapNotFound(err)
}

func (s *Store) CreateSeason(ctx context.Context, seriesID string, seasonNumber int, title string) (models.Season, error) {
	var se models.Season
	err := s.Pool.QueryRow(ctx, `
		insert into seasons (series_id, season_number, title) values ($1, $2, $3)
		returning id, series_id, season_number, title
	`, seriesID, seasonNumber, title).Scan(&se.ID, &se.SeriesID, &se.SeasonNumber, &se.Title)
	return se, err
}

func (s *Store) UpdateSeason(ctx context.Context, id string, seasonNumber int, title string) (models.Season, error) {
	var se models.Season
	err := s.Pool.QueryRow(ctx, `
		update seasons set season_number = $1, title = $2 where id = $3
		returning id, series_id, season_number, title
	`, seasonNumber, title, id).Scan(&se.ID, &se.SeriesID, &se.SeasonNumber, &se.Title)
	return se, mapNotFound(err)
}

func (s *Store) DeleteSeason(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, `delete from seasons where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

const episodeColumns = `id, season_id, title, description, episode_number, duration, video_url, thumbnail_url, is_published`

func scanEpisode(row interface{ Scan(dest ...any) error }) (models.Episode, error) {
	var e models.Episode
	err := row.Scan(&e.ID, &e.SeasonID, &e.Title, &e.Description, &e.EpisodeNumber, &e.Duration,
		&e.VideoURL, &e.ThumbnailURL, &e.IsPublished)
	return e, err
}

func (s *Store) ListEpisodes(ctx context.Context, seasonID string) ([]models.Episode, error) {
	rows, err := s.Pool.Query(ctx, `select `+episodeColumns+` from episodes where season_id = $1 order by episode_number asc`, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		e, err := scanEpisode(rows)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, e)
	}
	return episodes, rows.Err()
}

func (s *Store) GetEpisodeByID(ctx context.Context, id string) (models.Episode, error) {
	row := s.Pool.QueryRow(ctx, `select `+episodeColumns+` from episodes where id = $1`, id)
	e, err := scanEpisode(row)
	return e, mapNotFound(err)
}

// GetNextEpisode returns the next published episode in the same season
// after episodeNumber, or ErrNotFound if this was the last one.
func (s *Store) GetNextEpisode(ctx context.Context, seasonID string, episodeNumber int) (models.Episode, error) {
	row := s.Pool.QueryRow(ctx, `
		select `+episodeColumns+` from episodes
		where season_id = $1 and episode_number > $2 and is_published = true
		order by episode_number asc
		limit 1
	`, seasonID, episodeNumber)
	e, err := scanEpisode(row)
	return e, mapNotFound(err)
}

func (s *Store) CreateEpisode(ctx context.Context, seasonID, title, description string, episodeNumber, duration int, videoURL, thumbnailURL string, isPublished bool) (models.Episode, error) {
	row := s.Pool.QueryRow(ctx, `
		insert into episodes (season_id, title, description, episode_number, duration, video_url, thumbnail_url, is_published)
		values ($1,$2,$3,$4,$5,$6,$7,$8)
		returning `+episodeColumns,
		seasonID, title, description, episodeNumber, duration, videoURL, thumbnailURL, isPublished)
	return scanEpisode(row)
}

func (s *Store) UpdateEpisode(ctx context.Context, id, title, description string, episodeNumber, duration int, videoURL, thumbnailURL string, isPublished bool) (models.Episode, error) {
	row := s.Pool.QueryRow(ctx, `
		update episodes set title=$1, description=$2, episode_number=$3, duration=$4, video_url=$5,
			thumbnail_url=$6, is_published=$7
		where id = $8
		returning `+episodeColumns,
		title, description, episodeNumber, duration, videoURL, thumbnailURL, isPublished, id)
	e, err := scanEpisode(row)
	return e, mapNotFound(err)
}

func (s *Store) DeleteEpisode(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, `delete from episodes where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
