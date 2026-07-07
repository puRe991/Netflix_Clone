package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/pure991/streamflix/internal/models"
)

// CreateUserWithProfile creates a user and a default profile atomically,
// mirroring the original nested Prisma create. Returns the new user id.
func (s *Store) CreateUserWithProfile(ctx context.Context, email, passwordHash string, name *string, profileName string) (string, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var userID string
	err = tx.QueryRow(ctx, `
		insert into users (email, password_hash, name)
		values (lower($1), $2, $3)
		returning id
	`, email, passwordHash, name).Scan(&userID)
	if err != nil {
		return "", err
	}

	if _, err := tx.Exec(ctx, `
		insert into profiles (id, user_id, name)
		values ($1, $2, $3)
	`, uuid.NewString(), userID, profileName); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return userID, nil
}

func scanUser(row interface {
	Scan(dest ...any) error
}) (models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.StripeCustomerID,
		&u.SubscriptionStatus, &u.LastProfileID, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

const userColumns = `id, email, password_hash, name, role, stripe_customer_id, subscription_status, last_profile_id, created_at, updated_at`

func (s *Store) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	row := s.Pool.QueryRow(ctx, `select `+userColumns+` from users where email = lower($1)`, email)
	u, err := scanUser(row)
	return u, mapNotFound(err)
}

func (s *Store) GetUserByID(ctx context.Context, id string) (models.User, error) {
	row := s.Pool.QueryRow(ctx, `select `+userColumns+` from users where id = $1`, id)
	u, err := scanUser(row)
	return u, mapNotFound(err)
}

func (s *Store) GetSubscriptionStatus(ctx context.Context, userID string) (models.SubscriptionStatus, error) {
	var status models.SubscriptionStatus
	err := s.Pool.QueryRow(ctx, `select subscription_status from users where id = $1`, userID).Scan(&status)
	return status, mapNotFound(err)
}

func (s *Store) SetStripeCustomerID(ctx context.Context, userID, customerID string) error {
	_, err := s.Pool.Exec(ctx, `update users set stripe_customer_id = $1, updated_at = now() where id = $2`, customerID, userID)
	return err
}

func (s *Store) GetUserByStripeCustomerID(ctx context.Context, customerID string) (models.User, error) {
	row := s.Pool.QueryRow(ctx, `select `+userColumns+` from users where stripe_customer_id = $1`, customerID)
	u, err := scanUser(row)
	return u, mapNotFound(err)
}

func (s *Store) SetSubscriptionStatus(ctx context.Context, userID string, status models.SubscriptionStatus) error {
	_, err := s.Pool.Exec(ctx, `update users set subscription_status = $1, updated_at = now() where id = $2`, status, userID)
	return err
}

func (s *Store) SetUserRole(ctx context.Context, userID string, role models.Role) error {
	_, err := s.Pool.Exec(ctx, `update users set role = $1, updated_at = now() where id = $2`, role, userID)
	return err
}

func (s *Store) ListUsers(ctx context.Context) ([]models.User, error) {
	rows, err := s.Pool.Query(ctx, `select `+userColumns+` from users order by created_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *Store) CountUsers(ctx context.Context) (int64, error) {
	var n int64
	err := s.Pool.QueryRow(ctx, `select count(*) from users`).Scan(&n)
	return n, err
}

func (s *Store) CountActiveSubscriptions(ctx context.Context) (int64, error) {
	var n int64
	err := s.Pool.QueryRow(ctx, `select count(*) from users where subscription_status = 'ACTIVE'`).Scan(&n)
	return n, err
}
