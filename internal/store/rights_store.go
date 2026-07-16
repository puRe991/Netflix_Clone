package store

import (
	"context"
	"time"

	"github.com/pure991/streamflix/internal/models"
)

const rightsColumns = `id, media_id, source, rights_holder, license_doc_url, revenue_share_pct, territory_limit,
	valid_from, valid_until, notes, created_at, updated_at`

func scanRights(row interface{ Scan(dest ...any) error }) (models.RightsInfo, error) {
	var ri models.RightsInfo
	err := row.Scan(&ri.ID, &ri.MediaID, &ri.Source, &ri.RightsHolder, &ri.LicenseDocURL, &ri.RevenueSharePct,
		&ri.TerritoryLimit, &ri.ValidFrom, &ri.ValidUntil, &ri.Notes, &ri.CreatedAt, &ri.UpdatedAt)
	return ri, err
}

type RightsInfoInput struct {
	Source          models.RightsSource
	RightsHolder    *string
	LicenseDocURL   *string
	RevenueSharePct *float64
	TerritoryLimit  *string
	ValidFrom       time.Time
	ValidUntil      *time.Time
	Notes           *string
}

func (s *Store) CreateRightsInfo(ctx context.Context, mediaID string, in RightsInfoInput) (models.RightsInfo, error) {
	row := s.Pool.QueryRow(ctx, `
		insert into rights_info (media_id, source, rights_holder, license_doc_url, revenue_share_pct,
			territory_limit, valid_from, valid_until, notes)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		returning `+rightsColumns,
		mediaID, in.Source, in.RightsHolder, in.LicenseDocURL, in.RevenueSharePct, in.TerritoryLimit,
		in.ValidFrom, in.ValidUntil, in.Notes)
	return scanRights(row)
}

func (s *Store) GetRightsInfoByMediaID(ctx context.Context, mediaID string) (models.RightsInfo, error) {
	row := s.Pool.QueryRow(ctx, `select `+rightsColumns+` from rights_info where media_id = $1`, mediaID)
	ri, err := scanRights(row)
	return ri, mapNotFound(err)
}

func (s *Store) UpdateRightsInfo(ctx context.Context, mediaID string, in RightsInfoInput) (models.RightsInfo, error) {
	row := s.Pool.QueryRow(ctx, `
		update rights_info set source=$1, rights_holder=$2, license_doc_url=$3, revenue_share_pct=$4,
			territory_limit=$5, valid_from=$6, valid_until=$7, notes=$8, updated_at=now()
		where media_id = $9
		returning `+rightsColumns,
		in.Source, in.RightsHolder, in.LicenseDocURL, in.RevenueSharePct, in.TerritoryLimit,
		in.ValidFrom, in.ValidUntil, in.Notes, mediaID)
	ri, err := scanRights(row)
	return ri, mapNotFound(err)
}

func (s *Store) DeleteRightsInfo(ctx context.Context, mediaID string) error {
	tag, err := s.Pool.Exec(ctx, `delete from rights_info where media_id = $1`, mediaID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ListExpiringRights returns rights info whose ValidUntil falls within the
// next withinDays days (already-expired grants included), soonest first.
// Media title/slug are attached via ri.Media, matching the ListAllSeries
// join pattern.
func (s *Store) ListExpiringRights(ctx context.Context, withinDays int) ([]models.RightsInfo, error) {
	rows, err := s.Pool.Query(ctx, `
		select ri.id, ri.media_id, ri.source, ri.rights_holder, ri.license_doc_url, ri.revenue_share_pct,
			ri.territory_limit, ri.valid_from, ri.valid_until, ri.notes, ri.created_at, ri.updated_at,
			m.title, m.slug
		from rights_info ri
		join media m on m.id = ri.media_id
		where ri.valid_until is not null and ri.valid_until <= now() + make_interval(days => $1)
		order by ri.valid_until asc
	`, withinDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.RightsInfo
	for rows.Next() {
		var ri models.RightsInfo
		var title, slug string
		if err := rows.Scan(&ri.ID, &ri.MediaID, &ri.Source, &ri.RightsHolder, &ri.LicenseDocURL, &ri.RevenueSharePct,
			&ri.TerritoryLimit, &ri.ValidFrom, &ri.ValidUntil, &ri.Notes, &ri.CreatedAt, &ri.UpdatedAt, &title, &slug); err != nil {
			return nil, err
		}
		ri.Media = &models.Media{ID: ri.MediaID, Title: title, Slug: slug}
		list = append(list, ri)
	}
	return list, rows.Err()
}
