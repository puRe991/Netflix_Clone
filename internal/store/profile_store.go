package store

import (
	"context"

	"github.com/pure991/streamflix/internal/models"
)

const profileColumns = `id, user_id, name, avatar_url, is_kids_profile, created_at`

func scanProfile(row interface{ Scan(dest ...any) error }) (models.Profile, error) {
	var p models.Profile
	err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.AvatarURL, &p.IsKidsProfile, &p.CreatedAt)
	return p, err
}

func (s *Store) ListProfiles(ctx context.Context, userID string) ([]models.Profile, error) {
	rows, err := s.Pool.Query(ctx, `select `+profileColumns+` from profiles where user_id = $1 order by created_at asc`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		p, err := scanProfile(rows)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

// FirstOrCreateProfile returns the user's first profile, creating a default
// one if none exists yet (used by the watch page and watchlist add, which
// need a profile to attach progress/watchlist rows to).
func (s *Store) FirstOrCreateProfile(ctx context.Context, userID string, fallbackName string) (models.Profile, error) {
	row := s.Pool.QueryRow(ctx, `select `+profileColumns+` from profiles where user_id = $1 order by created_at asc limit 1`, userID)
	p, err := scanProfile(row)
	if err == nil {
		return p, nil
	}
	if mapNotFound(err) != ErrNotFound {
		return models.Profile{}, err
	}

	row = s.Pool.QueryRow(ctx, `
		insert into profiles (user_id, name)
		values ($1, $2)
		returning `+profileColumns, userID, fallbackName)
	return scanProfile(row)
}

func (s *Store) CreateProfile(ctx context.Context, userID, name string, avatarURL *string, isKids bool) (models.Profile, error) {
	row := s.Pool.QueryRow(ctx, `
		insert into profiles (user_id, name, avatar_url, is_kids_profile)
		values ($1, $2, $3, $4)
		returning `+profileColumns, userID, name, avatarURL, isKids)
	return scanProfile(row)
}

// UpdateProfile updates only the caller-owned profile matching id+userID;
// returns ErrNotFound if it doesn't exist or isn't owned by userID.
func (s *Store) UpdateProfile(ctx context.Context, id, userID, name string, avatarURL *string, isKids bool) (models.Profile, error) {
	row := s.Pool.QueryRow(ctx, `
		update profiles set name = $1, avatar_url = $2, is_kids_profile = $3
		where id = $4 and user_id = $5
		returning `+profileColumns, name, avatarURL, isKids, id, userID)
	p, err := scanProfile(row)
	return p, mapNotFound(err)
}

func (s *Store) DeleteProfile(ctx context.Context, id, userID string) error {
	tag, err := s.Pool.Exec(ctx, `delete from profiles where id = $1 and user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
