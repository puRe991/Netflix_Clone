// Package billing wraps the Stripe SDK for checkout, the customer portal,
// and webhook verification.
package billing

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pure991/streamflix/internal/models"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/webhook"
)

type Billing struct {
	secretKey     string
	basicPriceID  string
	premiumPrice  string
	webhookSecret string
	appURL        string
}

func New(secretKey, basicPriceID, premiumPriceID, webhookSecret, appURL string) *Billing {
	stripe.Key = secretKey
	return &Billing{
		secretKey:     secretKey,
		basicPriceID:  basicPriceID,
		premiumPrice:  premiumPriceID,
		webhookSecret: webhookSecret,
		appURL:        appURL,
	}
}

func (b *Billing) HasWebhookSecret() bool { return b.webhookSecret != "" }

func (b *Billing) priceIDFor(plan models.Plan) string {
	if plan == models.PlanPremium {
		return b.premiumPrice
	}
	return b.basicPriceID
}

func (b *Billing) CreateCustomer(email, userID string) (string, error) {
	c, err := customer.New(&stripe.CustomerParams{
		Email:  stripe.String(email),
		Params: stripe.Params{Metadata: map[string]string{"userId": userID}},
	})
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (b *Billing) CreateCheckoutSession(customerID, userID string, plan models.Plan) (string, error) {
	params := &stripe.CheckoutSessionParams{
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Customer: stripe.String(customerID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{Price: stripe.String(b.priceIDFor(plan)), Quantity: stripe.Int64(1)},
		},
		SuccessURL: stripe.String(b.appURL + "/billing?success=1"),
		CancelURL:  stripe.String(b.appURL + "/pricing"),
		Params:     stripe.Params{Metadata: map[string]string{"userId": userID, "plan": string(plan)}},
	}
	sess, err := checkoutsession.New(params)
	if err != nil {
		return "", err
	}
	return sess.URL, nil
}

func (b *Billing) CreatePortalSession(customerID string) (string, error) {
	sess, err := session.New(&stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(b.appURL + "/account"),
	})
	if err != nil {
		return "", err
	}
	return sess.URL, nil
}

var ErrWebhookNotConfigured = errors.New("stripe webhook secret is not configured")

// ConstructEvent verifies the Stripe-Signature header against the raw body.
// Unlike the previous implementation, there is no unverified fallback: if
// the webhook secret isn't configured, verification always fails closed.
func (b *Billing) ConstructEvent(payload []byte, signatureHeader string) (stripe.Event, error) {
	if b.webhookSecret == "" {
		return stripe.Event{}, ErrWebhookNotConfigured
	}
	return webhook.ConstructEvent(payload, signatureHeader, b.webhookSecret)
}

// MapStripeStatus mirrors the original mapStatus() mapping.
func MapStripeStatus(s stripe.SubscriptionStatus) models.SubscriptionStatus {
	switch s {
	case stripe.SubscriptionStatusActive:
		return models.SubStatusActive
	case stripe.SubscriptionStatusTrialing:
		return models.SubStatusTrialing
	case stripe.SubscriptionStatusPastDue:
		return models.SubStatusPastDue
	case stripe.SubscriptionStatusCanceled:
		return models.SubStatusCanceled
	default:
		return models.SubStatusIncomplete
	}
}

func PlanFromMetadata(metadata map[string]string) models.Plan {
	if metadata["plan"] == string(models.PlanPremium) {
		return models.PlanPremium
	}
	return models.PlanBasic
}

// ParseSubscriptionEvent extracts the stripe.Subscription object embedded in
// a customer.subscription.* event.
func ParseSubscriptionEvent(event []byte) (stripe.Subscription, error) {
	var sub stripe.Subscription
	if err := json.Unmarshal(event, &sub); err != nil {
		return sub, fmt.Errorf("parse subscription event: %w", err)
	}
	return sub, nil
}
