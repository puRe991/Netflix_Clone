package httpserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/models"
)

// spoilerCookie stores the spoiler preference of anonymous visitors.
// Logged-in users keep it on their profile instead.
const spoilerCookie = "streamflix_spoilers"

// viewer bundles who is looking at a page and how results must be shown.
type viewer struct {
	User         *auth.SessionUser
	Profile      *models.Profile
	HideSpoilers bool
}

// currentViewer resolves the viewer's active profile (the first one, the
// same profile the watch pages attach progress to) and spoiler preference.
// Spoilers are hidden unless the viewer explicitly opted in.
func (s *Server) currentViewer(ctx context.Context, r *http.Request) (viewer, error) {
	v := viewer{User: s.Sessions.OptionalUser(r), HideSpoilers: true}
	if v.User == nil {
		if c, err := r.Cookie(spoilerCookie); err == nil && c.Value == "show" {
			v.HideSpoilers = false
		}
		return v, nil
	}

	profile, err := s.Store.FirstOrCreateProfile(ctx, v.User.ID, displayName(*v.User))
	if err != nil {
		return v, err
	}
	v.Profile = &profile
	v.HideSpoilers = profile.HideSpoilers
	return v, nil
}

// safeLocalPath only allows same-site relative redirects, so the "next"
// form value can't be abused as an open redirect.
func safeLocalPath(p, fallback string) string {
	if !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") || strings.HasPrefix(p, "/\\") {
		return fallback
	}
	return p
}

// SetSpoilers switches spoiler-free mode on or off for the current viewer.
func (s *Server) SetSpoilers(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}
	show := r.FormValue("show") == "1"

	v, err := s.currentViewer(r.Context(), r)
	if err != nil {
		return err
	}

	if v.Profile != nil {
		if err := s.Store.SetProfileHideSpoilers(r.Context(), v.Profile.ID, v.User.ID, !show); err != nil {
			return err
		}
	} else {
		value, maxAge := "hide", -1
		if show {
			value, maxAge = "show", 365*24*3600
		}
		http.SetCookie(w, &http.Cookie{
			Name:     spoilerCookie,
			Value:    value,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   maxAge,
		})
	}

	redirect(w, r, safeLocalPath(r.FormValue("next"), "/browse"))
	return nil
}

// SetProfileSpoilers toggles spoiler-free mode for one specific profile
// (from the profiles page).
func (s *Server) SetProfileSpoilers(w http.ResponseWriter, r *http.Request) error {
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
	hide := r.FormValue("hideSpoilers") == "1"
	if err := s.Store.SetProfileHideSpoilers(r.Context(), r.PathValue("id"), user.ID, hide); err != nil {
		return err
	}
	redirect(w, r, "/profiles")
	return nil
}

type teamsData struct {
	PageData
	Teams     []models.Team
	Followed  map[string]bool
	CanFollow bool
}

func (s *Server) TeamsPage(w http.ResponseWriter, r *http.Request) error {
	teams, err := s.Store.ListTeams(r.Context())
	if err != nil {
		return err
	}
	v, err := s.currentViewer(r.Context(), r)
	if err != nil {
		return err
	}

	data := teamsData{PageData: s.basePageData(w, r), Teams: teams, Followed: map[string]bool{}}
	if v.Profile != nil {
		data.CanFollow = true
		if data.Followed, err = s.Store.FollowedTeamIDs(r.Context(), v.Profile.ID); err != nil {
			return err
		}
	}
	return s.render(w, r, "teams.html", data)
}

type teamData struct {
	PageData
	Team         models.Team
	Matches      []models.TeamMatch
	HideSpoilers bool
	IsFollowed   bool
	CanFollow    bool
	Path         string
}

func (s *Server) TeamPage(w http.ResponseWriter, r *http.Request) error {
	team, err := s.Store.GetTeamBySlug(r.Context(), r.PathValue("slug"))
	if err != nil {
		return err
	}
	matches, err := s.Store.ListMatchesForTeam(r.Context(), team.ID)
	if err != nil {
		return err
	}
	v, err := s.currentViewer(r.Context(), r)
	if err != nil {
		return err
	}

	data := teamData{
		PageData:     s.basePageData(w, r),
		Team:         team,
		Matches:      matches,
		HideSpoilers: v.HideSpoilers,
		Path:         r.URL.Path,
	}
	if v.Profile != nil {
		data.CanFollow = true
		followed, err := s.Store.FollowedTeamIDs(r.Context(), v.Profile.ID)
		if err != nil {
			return err
		}
		data.IsFollowed = followed[team.ID]
	}
	return s.render(w, r, "team.html", data)
}

func (s *Server) followTeam(w http.ResponseWriter, r *http.Request, follow bool) error {
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
	team, err := s.Store.GetTeamByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	profile, err := s.Store.FirstOrCreateProfile(r.Context(), user.ID, displayName(user))
	if err != nil {
		return err
	}

	if follow {
		err = s.Store.FollowTeam(r.Context(), profile.ID, team.ID)
	} else {
		err = s.Store.UnfollowTeam(r.Context(), profile.ID, team.ID)
	}
	if err != nil {
		return err
	}

	redirect(w, r, safeLocalPath(r.FormValue("next"), "/team/"+team.Slug))
	return nil
}

func (s *Server) FollowTeam(w http.ResponseWriter, r *http.Request) error {
	return s.followTeam(w, r, true)
}

func (s *Server) UnfollowTeam(w http.ResponseWriter, r *http.Request) error {
	return s.followTeam(w, r, false)
}
