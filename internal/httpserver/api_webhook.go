package httpserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pure991/streamflix/internal/billing"
	"github.com/pure991/streamflix/internal/store"
)

func (s *Server) StripeWebhook(w http.ResponseWriter, r *http.Request) error {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return validationErr("Konnte Request-Body nicht lesen")
	}

	event, err := s.Billing.ConstructEvent(payload, r.Header.Get("Stripe-Signature"))
	if err != nil {
		// Always fail closed: an unverifiable webhook is rejected outright,
		// there is no "trust it anyway" fallback here.
		JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid webhook signature"})
		return nil
	}

	if strings.HasPrefix(string(event.Type), "customer.subscription.") {
		sub, err := billing.ParseSubscriptionEvent(event.Data.Raw)
		if err != nil {
			return err
		}

		user, err := s.Store.GetUserByStripeCustomerID(r.Context(), sub.Customer.ID)
		if err != nil {
			if err == store.ErrNotFound {
				log.Printf("stripe webhook: no user for customer %s, dropping event %s", sub.Customer.ID, event.ID)
				JSON(w, http.StatusOK, map[string]bool{"received": true})
				return nil
			}
			return err
		}

		status := billing.MapStripeStatus(sub.Status)
		plan := billing.PlanFromMetadata(sub.Metadata)

		var periodEndTime *time.Time
		if sub.CurrentPeriodEnd > 0 {
			t := time.Unix(sub.CurrentPeriodEnd, 0)
			periodEndTime = &t
		}

		if err := s.Store.UpsertSubscription(r.Context(), user.ID, sub.ID, plan, status, periodEndTime); err != nil {
			return err
		}
		if err := s.Store.SetSubscriptionStatus(r.Context(), user.ID, status); err != nil {
			return err
		}

		raw, _ := json.Marshal(event)
		if err := s.Store.InsertPaymentLogIdempotent(r.Context(), &user.ID, event.ID, string(event.Type), raw); err != nil {
			return err
		}
	}

	JSON(w, http.StatusOK, map[string]bool{"received": true})
	return nil
}
