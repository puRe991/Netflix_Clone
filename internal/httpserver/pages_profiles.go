package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

type profilesData struct {
	PageData
	Profiles []models.Profile
}

func (s *Server) ProfilesPage(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	profiles, err := s.Store.ListProfiles(r.Context(), user.ID)
	if err != nil {
		return err
	}

	return s.render(w, r, "profiles.html", profilesData{
		PageData: s.basePageData(w, r),
		Profiles: profiles,
	})
}
