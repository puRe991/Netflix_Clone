package httpserver

import (
	"encoding/json"
	"net/http"
)

func (s *Server) APIListProfiles(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}
	profiles, err := s.Store.ListProfiles(r.Context(), user.ID)
	if err != nil {
		return err
	}
	JSON(w, http.StatusOK, profiles)
	return nil
}

func (s *Server) CreateProfile(w http.ResponseWriter, r *http.Request) error {
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

	name, err := requireMinMaxLen("Name", r.FormValue("name"), 2, 40)
	if err != nil {
		return err
	}
	avatarURL, err := optionalURL("avatarUrl", r.FormValue("avatarUrl"))
	if err != nil {
		return err
	}
	isKids := r.FormValue("isKidsProfile") == "on"

	if _, err := s.Store.CreateProfile(r.Context(), user.ID, name, avatarURL, isKids); err != nil {
		return err
	}

	redirect(w, r, "/profiles")
	return nil
}

type updateProfileRequest struct {
	Name          string  `json:"name"`
	AvatarURL     *string `json:"avatarUrl"`
	IsKidsProfile bool    `json:"isKidsProfile"`
}

func (s *Server) UpdateProfile(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	var body updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return validationErr("Ungültiger Anfrage-Body")
	}
	name, err := requireMinMaxLen("Name", body.Name, 2, 40)
	if err != nil {
		return err
	}

	profile, err := s.Store.UpdateProfile(r.Context(), r.PathValue("id"), user.ID, name, body.AvatarURL, body.IsKidsProfile)
	if err != nil {
		return err
	}

	JSON(w, http.StatusOK, profile)
	return nil
}

func (s *Server) DeleteProfile(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	if err := s.Store.DeleteProfile(r.Context(), r.PathValue("id"), user.ID); err != nil {
		return err
	}

	JSON(w, http.StatusOK, map[string]bool{"ok": true})
	return nil
}
