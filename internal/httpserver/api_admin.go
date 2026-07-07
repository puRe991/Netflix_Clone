package httpserver

import (
	"encoding/json"
	"net/http"
)

func (s *Server) APIAdminStats(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	stats, err := s.Store.GetStats(r.Context())
	if err != nil {
		return err
	}
	JSON(w, http.StatusOK, stats)
	return nil
}

func (s *Server) APIAdminSubscriptions(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	subs, err := s.Store.ListSubscriptions(r.Context())
	if err != nil {
		return err
	}
	JSON(w, http.StatusOK, subs)
	return nil
}

func (s *Server) APIAdminUsers(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	users, err := s.Store.ListUsers(r.Context())
	if err != nil {
		return err
	}
	// passwordHash intentionally excluded from the wire format
	type publicUser struct {
		ID                 string
		Email              string
		Name               *string
		Role               string
		SubscriptionStatus string
		CreatedAt          string
	}
	out := make([]publicUser, 0, len(users))
	for _, u := range users {
		out = append(out, publicUser{
			ID: u.ID, Email: u.Email, Name: u.Name, Role: string(u.Role),
			SubscriptionStatus: string(u.SubscriptionStatus), CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	JSON(w, http.StatusOK, out)
	return nil
}

type createGenreRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// APIAdminGenreCreate is the JSON API surface for programmatic genre
// creation (distinct from the admin UI's form-based /admin/genres route).
func (s *Server) APIAdminGenreCreate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}

	var body createGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return validationErr("Ungültiger Anfrage-Body")
	}
	name, err := requireMinLen("Name", body.Name, 2)
	if err != nil {
		return err
	}
	slug, err := requireSlug(body.Slug)
	if err != nil {
		return err
	}

	genre, err := s.Store.CreateGenre(r.Context(), name, slug)
	if err != nil {
		return err
	}
	JSON(w, http.StatusOK, genre)
	return nil
}
