package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

func (s *Server) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	plan, err := requireEnum("plan", r.FormValue("plan"), string(models.PlanBasic), string(models.PlanPremium))
	if err != nil {
		return err
	}

	dbUser, err := s.Store.GetUserByID(r.Context(), user.ID)
	if err != nil {
		return err
	}

	customerID := ""
	if dbUser.StripeCustomerID != nil {
		customerID = *dbUser.StripeCustomerID
	} else {
		customerID, err = s.Billing.CreateCustomer(dbUser.Email, dbUser.ID)
		if err != nil {
			return err
		}
		if err := s.Store.SetStripeCustomerID(r.Context(), dbUser.ID, customerID); err != nil {
			return err
		}
	}

	checkoutURL, err := s.Billing.CreateCheckoutSession(customerID, dbUser.ID, models.Plan(plan))
	if err != nil {
		return err
	}

	http.Redirect(w, r, checkoutURL, http.StatusSeeOther)
	return nil
}

func (s *Server) CreatePortalSession(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	dbUser, err := s.Store.GetUserByID(r.Context(), user.ID)
	if err != nil {
		return err
	}
	if dbUser.StripeCustomerID == nil {
		redirect(w, r, "/pricing")
		return nil
	}

	portalURL, err := s.Billing.CreatePortalSession(*dbUser.StripeCustomerID)
	if err != nil {
		return err
	}

	http.Redirect(w, r, portalURL, http.StatusSeeOther)
	return nil
}
