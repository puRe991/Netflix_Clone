package httpserver

import "net/http"

func (s *Server) APIListWatchlist(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}
	items, err := s.Store.ListWatchlist(r.Context(), user.ID)
	if err != nil {
		return err
	}
	JSON(w, http.StatusOK, items)
	return nil
}

func (s *Server) AddToWatchlist(w http.ResponseWriter, r *http.Request) error {
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

	mediaID := r.FormValue("mediaId")
	if mediaID == "" {
		return validationErr("mediaId ist erforderlich")
	}

	profile, err := s.Store.FirstOrCreateProfile(r.Context(), user.ID, displayName(user))
	if err != nil {
		return err
	}

	if err := s.Store.AddToWatchlist(r.Context(), user.ID, profile.ID, mediaID); err != nil {
		return err
	}

	redirect(w, r, "/my-list")
	return nil
}

func (s *Server) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}
	if err := s.Store.RemoveFromWatchlist(r.Context(), r.PathValue("id"), user.ID); err != nil {
		return err
	}
	JSON(w, http.StatusOK, map[string]bool{"ok": true})
	return nil
}
