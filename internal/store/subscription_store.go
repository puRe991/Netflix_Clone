package store

import (
	"context"
	"time"

	"github.com/pure991/streamflix/internal/models"
)

func (s *Store) UpsertSubscription(ctx context.Context, userID, stripeSubscriptionID string, plan models.Plan, status models.SubscriptionStatus, currentPeriodEnd *time.Time) error {
	_, err := s.Pool.Exec(ctx, `
		insert into subscriptions (user_id, stripe_subscription_id, plan, status, current_period_end)
		values ($1, $2, $3, $4, $5)
		on conflict (stripe_subscription_id)
		do update set plan = excluded.plan, status = excluded.status, current_period_end = excluded.current_period_end, updated_at = now()
	`, userID, stripeSubscriptionID, plan, status, currentPeriodEnd)
	return err
}

func (s *Store) ListSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	rows, err := s.Pool.Query(ctx, `
		select s.id, s.user_id, s.stripe_subscription_id, s.plan, s.status, s.current_period_end,
			s.created_at, s.updated_at, u.email, u.name
		from subscriptions s
		join users u on u.id = s.user_id
		order by s.created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.StripeSubscriptionID, &sub.Plan, &sub.Status,
			&sub.CurrentPeriodEnd, &sub.CreatedAt, &sub.UpdatedAt, &sub.UserEmail, &sub.UserName); err != nil {
			return nil, err
		}
		list = append(list, sub)
	}
	return list, rows.Err()
}

// InsertPaymentLogIdempotent records a webhook event for audit/idempotency,
// doing nothing if stripeEventID was already recorded.
func (s *Store) InsertPaymentLogIdempotent(ctx context.Context, userID *string, stripeEventID, eventType string, rawData []byte) error {
	_, err := s.Pool.Exec(ctx, `
		insert into payment_logs (user_id, stripe_event_id, type, raw_data)
		values ($1, $2, $3, $4)
		on conflict (stripe_event_id) do nothing
	`, userID, stripeEventID, eventType, rawData)
	return err
}
