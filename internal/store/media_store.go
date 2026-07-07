package store

import (
	"context"

	"github.com/pure991/streamflix/internal/models"
)

const mediaColumns = `id, title, slug, description, type, thumbnail_url, banner_url, trailer_url, video_url,
	release_year, duration, age_rating, language, tags, is_published, created_at, updated_at`

// mediaColumnsQualified is mediaColumns with an "m." table-alias prefix, for
// use in joins where an unqualified "id" would otherwise be ambiguous.
const mediaColumnsQualified = `m.id, m.title, m.slug, m.description, m.type, m.thumbnail_url, m.banner_url, m.trailer_url, m.video_url,
	m.release_year, m.duration, m.age_rating, m.language, m.tags, m.is_published, m.created_at, m.updated_at`

func scanMedia(row interface{ Scan(dest ...any) error }) (models.Media, error) {
	var m models.Media
	err := row.Scan(&m.ID, &m.Title, &m.Slug, &m.Description, &m.Type, &m.ThumbnailURL, &m.BannerURL,
		&m.TrailerURL, &m.VideoURL, &m.ReleaseYear, &m.Duration, &m.AgeRating, &m.Language, &m.Tags,
		&m.IsPublished, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (s *Store) ListPublishedMedia(ctx context.Context) ([]models.Media, error) {
	rows, err := s.Pool.Query(ctx, `select `+mediaColumns+` from media where is_published = true order by created_at desc`)
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

func (s *Store) ListAllMedia(ctx context.Context) ([]models.Media, error) {
	rows, err := s.Pool.Query(ctx, `select `+mediaColumns+` from media order by created_at desc`)
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

func (s *Store) GetMediaByID(ctx context.Context, id string) (models.Media, error) {
	row := s.Pool.QueryRow(ctx, `select `+mediaColumns+` from media where id = $1`, id)
	m, err := scanMedia(row)
	return m, mapNotFound(err)
}

// GetMediaBySlug loads a media row plus its genres and, for series, the full
// season/episode tree. onlyPublished restricts the lookup to published
// media (used by the public listing API; the detail page/API intentionally
// allows unpublished lookups by slug, matching prior behavior).
func (s *Store) GetMediaBySlug(ctx context.Context, slug string, onlyPublished bool) (models.Media, error) {
	query := `select ` + mediaColumns + ` from media where slug = $1`
	if onlyPublished {
		query += ` and is_published = true`
	}
	row := s.Pool.QueryRow(ctx, query, slug)
	m, err := scanMedia(row)
	if err != nil {
		return m, mapNotFound(err)
	}

	genres, err := s.GetGenresForMedia(ctx, m.ID)
	if err != nil {
		return m, err
	}
	m.Genres = genres

	if m.Type == models.MediaSeries {
		series, err := s.GetSeriesByMediaID(ctx, m.ID)
		if err == nil {
			m.Series = &series
		} else if err != ErrNotFound {
			return m, err
		}
	}

	return m, nil
}

func (s *Store) GetGenresForMedia(ctx context.Context, mediaID string) ([]models.Genre, error) {
	rows, err := s.Pool.Query(ctx, `
		select g.id, g.name, g.slug
		from genres g
		join media_genres mg on mg.genre_id = g.id
		where mg.media_id = $1
		order by g.name asc
	`, mediaID)
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

type MediaInput struct {
	Title        string
	Slug         string
	Description  string
	Type         models.MediaType
	ThumbnailURL string
	BannerURL    string
	TrailerURL   *string
	VideoURL     *string
	ReleaseYear  int
	Duration     *int
	AgeRating    string
	Language     string
	Tags         []string
	IsPublished  bool
	GenreIDs     []string
}

func (s *Store) CreateMedia(ctx context.Context, in MediaInput) (models.Media, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return models.Media{}, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
		insert into media (title, slug, description, type, thumbnail_url, banner_url, trailer_url, video_url,
			release_year, duration, age_rating, language, tags, is_published)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		returning `+mediaColumns,
		in.Title, in.Slug, in.Description, in.Type, in.ThumbnailURL, in.BannerURL, in.TrailerURL, in.VideoURL,
		in.ReleaseYear, in.Duration, in.AgeRating, in.Language, in.Tags, in.IsPublished)
	m, err := scanMedia(row)
	if err != nil {
		return m, err
	}

	for _, genreID := range in.GenreIDs {
		if _, err := tx.Exec(ctx, `insert into media_genres (media_id, genre_id) values ($1, $2)`, m.ID, genreID); err != nil {
			return m, err
		}
	}

	return m, tx.Commit(ctx)
}

func (s *Store) UpdateMedia(ctx context.Context, id string, in MediaInput) (models.Media, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return models.Media{}, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
		update media set title=$1, slug=$2, description=$3, type=$4, thumbnail_url=$5, banner_url=$6,
			trailer_url=$7, video_url=$8, release_year=$9, duration=$10, age_rating=$11, language=$12,
			tags=$13, is_published=$14, updated_at=now()
		where id = $15
		returning `+mediaColumns,
		in.Title, in.Slug, in.Description, in.Type, in.ThumbnailURL, in.BannerURL, in.TrailerURL, in.VideoURL,
		in.ReleaseYear, in.Duration, in.AgeRating, in.Language, in.Tags, in.IsPublished, id)
	m, err := scanMedia(row)
	if err != nil {
		return m, mapNotFound(err)
	}

	if _, err := tx.Exec(ctx, `delete from media_genres where media_id = $1`, id); err != nil {
		return m, err
	}
	for _, genreID := range in.GenreIDs {
		if _, err := tx.Exec(ctx, `insert into media_genres (media_id, genre_id) values ($1, $2)`, id, genreID); err != nil {
			return m, err
		}
	}

	return m, tx.Commit(ctx)
}

func (s *Store) DeleteMedia(ctx context.Context, id string) error {
	tag, err := s.Pool.Exec(ctx, `delete from media where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) CountMedia(ctx context.Context) (int64, error) {
	var n int64
	err := s.Pool.QueryRow(ctx, `select count(*) from media`).Scan(&n)
	return n, err
}
