package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

type accountData struct {
	PageData
	Email              string
	SubscriptionStatus models.SubscriptionStatus
}

func (s *Server) AccountPage(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	// Always re-read subscription status from the database rather than the
	// (up to 7-day-old) JWT claim, so a lapsed or newly-active Stripe
	// subscription is reflected immediately.
	status, err := s.Store.GetSubscriptionStatus(r.Context(), user.ID)
	if err != nil {
		return err
	}

	return s.render(w, r, "account.html", accountData{
		PageData:           s.basePageData(w, r),
		Email:              user.Email,
		SubscriptionStatus: status,
	})
}

type billingData struct {
	PageData
}

func (s *Server) BillingPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireUser(r); err != nil {
		return err
	}
	return s.render(w, r, "billing.html", billingData{PageData: s.basePageData(w, r)})
}

type myListData struct {
	PageData
	Items []models.WatchlistItem
}

func (s *Server) MyListPage(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	items, err := s.Store.ListWatchlist(r.Context(), user.ID)
	if err != nil {
		return err
	}

	return s.render(w, r, "my-list.html", myListData{PageData: s.basePageData(w, r), Items: items})
}
