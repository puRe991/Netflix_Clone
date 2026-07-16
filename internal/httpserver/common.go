package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

// PageData is embedded by every page-specific data struct so templates can
// always rely on .User / .CSRFToken / .Flash being present.
type PageData struct {
	User      *auth.SessionUser
	CSRFToken string
	Flash     string
}

func (s *Server) basePageData(w http.ResponseWriter, r *http.Request) PageData {
	return PageData{
		User:      s.Sessions.OptionalUser(r),
		CSRFToken: s.Sessions.CSRFToken(w, r),
		Flash:     r.URL.Query().Get("flash"),
	}
}

func (s *Server) checkCSRF(r *http.Request) error {
	token := r.FormValue("csrf_token")
	if token == "" {
		token = r.Header.Get("X-CSRF-Token")
	}
	return s.Sessions.VerifyCSRF(r, token)
}

// canStreamFullContent decides whether a session user may play protected
// video. Subscription status is always re-read from the database here
// (never trusted from the JWT), so a lapsed/renewed Stripe subscription
// takes effect on the very next request.
func (s *Server) canStreamFullContent(ctx context.Context, u *auth.SessionUser) (bool, error) {
	if u == nil {
		return false, nil
	}
	if u.Role == models.RoleAdmin {
		return true, nil
	}
	status, err := s.Store.GetSubscriptionStatus(ctx, u.ID)
	if err != nil {
		return false, err
	}
	return status.CanStream(), nil
}

// assertRightsCurrent blocks stream delivery once a media title's tracked
// rights have lapsed (ValidUntil in the past), regardless of IsPublished.
// Media without tracked rights info (the common case) streams normally, so
// it returns store.ErrNotFound only when a RightsInfo row exists and has
// expired — the same "looks like it doesn't exist" signal the handlers
// already use for e.g. a missing video URL.
func (s *Server) assertRightsCurrent(ctx context.Context, mediaID string) error {
	rights, err := s.Store.GetRightsInfoByMediaID(ctx, mediaID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil
		}
		return err
	}
	if rights.ValidUntil != nil && rights.ValidUntil.Before(time.Now()) {
		return store.ErrNotFound
	}
	return nil
}
