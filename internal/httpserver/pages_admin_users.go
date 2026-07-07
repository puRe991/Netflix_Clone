package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

type adminUsersData struct {
	PageData
	Users []models.User
}

func (s *Server) AdminUsersPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	users, err := s.Store.ListUsers(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_users.html", adminUsersData{PageData: s.basePageData(w, r), Users: users})
}

func (s *Server) AdminUserToggleRole(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	user, err := s.Store.GetUserByID(r.Context(), id)
	if err != nil {
		return err
	}

	newRole := models.RoleAdmin
	if user.Role == models.RoleAdmin {
		newRole = models.RoleUser
	}
	if err := s.Store.SetUserRole(r.Context(), id, newRole); err != nil {
		return err
	}

	redirect(w, r, "/admin/users")
	return nil
}

type adminSubscriptionsData struct {
	PageData
	Subscriptions []models.Subscription
}

func (s *Server) AdminSubscriptionsPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	subs, err := s.Store.ListSubscriptions(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_subscriptions.html", adminSubscriptionsData{PageData: s.basePageData(w, r), Subscriptions: subs})
}
