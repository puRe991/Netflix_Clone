package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/store"
)

func (s *Server) APIGetWatch(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	canStream, err := s.canStreamFullContent(r.Context(), &user)
	if err != nil {
		return err
	}
	if !canStream {
		JSON(w, http.StatusPaymentRequired, map[string]string{"error": "Subscription required"})
		return nil
	}

	media, err := s.Store.GetMediaByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}

	JSON(w, http.StatusOK, map[string]any{
		"id":         media.ID,
		"title":      media.Title,
		"videoUrl":   media.VideoURL,
		"trailerUrl": media.TrailerURL,
	})
	return nil
}

type progressRequest struct {
	MediaID         string  `json:"mediaId"`
	ProfileID       string  `json:"profileId"`
	EpisodeID       *string `json:"episodeId"`
	ProgressSeconds int     `json:"progressSeconds"`
	Completed       bool    `json:"completed"`
}

func (s *Server) APIPostWatchProgress(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	var body progressRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return validationErr("Ungültiger Anfrage-Body")
	}
	if body.MediaID == "" || body.ProfileID == "" {
		return validationErr("mediaId und profileId sind erforderlich")
	}
	if body.ProgressSeconds < 0 {
		return validationErr("progressSeconds darf nicht negativ sein")
	}

	if err := s.Store.VerifyProfileOwnership(r.Context(), body.ProfileID, user.ID); err != nil {
		if err == store.ErrNotFound {
			return auth.ErrForbidden
		}
		return err
	}

	progress, err := s.Store.UpsertProgress(r.Context(), user.ID, body.ProfileID, body.MediaID, body.EpisodeID, body.ProgressSeconds, body.Completed)
	if err != nil {
		return err
	}

	JSON(w, http.StatusOK, progress)
	return nil
}

func (s *Server) APIGetWatchProgressForMedia(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	list, err := s.Store.ListProgressForUserMedia(r.Context(), user.ID, r.PathValue("id"))
	if err != nil {
		return err
	}

	JSON(w, http.StatusOK, list)
	return nil
}
