package httpserver

import (
	"net/http"
	"strings"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

type watchData struct {
	PageData
	Title           string
	VideoURL        string
	MediaID         string
	EpisodeID       string
	ProfileID       string
	InitialProgress int
	NextHref        string

	// Set for official YouTube embeds (eSports VODs): the player plays
	// these video IDs in order instead of a <video> element.
	YouTubeIDs   string
	InitialPart  int
	HideSpoilers bool
	BackHref     string
}

func (s *Server) WatchMoviePage(w http.ResponseWriter, r *http.Request) error {
	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	canStream, err := s.canStreamFullContent(r.Context(), &user)
	if err != nil {
		return err
	}
	if !canStream {
		redirect(w, r, "/pricing")
		return nil
	}

	media, err := s.Store.GetMediaByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	if media.VideoURL == nil || *media.VideoURL == "" {
		return store.ErrNotFound
	}
	if err := s.assertRightsCurrent(r.Context(), media.ID); err != nil {
		return err
	}

	profile, err := s.Store.FirstOrCreateProfile(r.Context(), user.ID, displayName(user))
	if err != nil {
		return err
	}

	initial := 0
	if progress, err := s.Store.GetMovieProgress(r.Context(), profile.ID, media.ID); err == nil {
		initial = progress.ProgressSeconds
	} else if err != store.ErrNotFound {
		return err
	}

	return s.render(w, r, "watch.html", watchData{
		PageData:        s.basePageData(w, r),
		Title:           media.Title,
		VideoURL:        *media.VideoURL,
		MediaID:         media.ID,
		ProfileID:       profile.ID,
		InitialProgress: initial,
	})
}

func (s *Server) WatchEpisodePage(w http.ResponseWriter, r *http.Request) error {
	episode, err := s.Store.GetEpisodeByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	season, err := s.Store.GetSeasonByID(r.Context(), episode.SeasonID)
	if err != nil {
		return err
	}
	series, err := s.Store.GetSeriesByID(r.Context(), season.SeriesID)
	if err != nil {
		return err
	}
	media, err := s.Store.GetMediaByID(r.Context(), series.MediaID)
	if err != nil {
		return err
	}

	parts := episode.VideoParts()
	if allYouTube(parts) {
		return s.watchEmbeddedEpisode(w, r, episode, season, media, parts)
	}

	user, err := s.Sessions.RequireUser(r)
	if err != nil {
		return err
	}

	canStream, err := s.canStreamFullContent(r.Context(), &user)
	if err != nil {
		return err
	}
	if !canStream {
		redirect(w, r, "/pricing")
		return nil
	}

	if err := s.assertRightsCurrent(r.Context(), media.ID); err != nil {
		return err
	}

	profile, err := s.Store.FirstOrCreateProfile(r.Context(), user.ID, displayName(user))
	if err != nil {
		return err
	}

	initial := 0
	if progress, err := s.Store.GetEpisodeProgress(r.Context(), profile.ID, episode.ID); err == nil {
		initial = progress.ProgressSeconds
	} else if err != store.ErrNotFound {
		return err
	}

	nextHref, err := s.nextEpisodeHref(r, season, episode)
	if err != nil {
		return err
	}

	return s.render(w, r, "watch.html", watchData{
		PageData:        s.basePageData(w, r),
		Title:           media.Title + " – " + episode.Title,
		VideoURL:        episode.VideoURL,
		MediaID:         media.ID,
		EpisodeID:       episode.ID,
		ProfileID:       profile.ID,
		InitialProgress: initial,
		NextHref:        nextHref,
	})
}

// watchEmbeddedEpisode plays an episode whose VOD parts are all official
// YouTube embeds. YouTube's API policies forbid charging for or otherwise
// gating playback in the embedded player beyond the play click, so unlike
// self-hosted videos this needs neither a login nor a subscription. Being
// public, it does insist on the episode and its title being published.
func (s *Server) watchEmbeddedEpisode(w http.ResponseWriter, r *http.Request, episode models.Episode, season models.Season, media models.Media, parts []string) error {
	user := s.Sessions.OptionalUser(r)
	isAdmin := user != nil && user.Role == models.RoleAdmin
	if !isAdmin && (!episode.IsPublished || !media.IsPublished) {
		return store.ErrNotFound
	}
	if err := s.assertRightsCurrent(r.Context(), media.ID); err != nil {
		return err
	}

	v, err := s.currentViewer(r.Context(), r)
	if err != nil {
		return err
	}

	ids := make([]string, len(parts))
	for i, p := range parts {
		ids[i] = youTubeVideoID(p)
	}

	data := watchData{
		PageData:     s.basePageData(w, r),
		Title:        media.Title + " – " + episode.Title,
		MediaID:      media.ID,
		EpisodeID:    episode.ID,
		YouTubeIDs:   strings.Join(ids, ","),
		HideSpoilers: v.HideSpoilers,
		BackHref:     "/series/" + media.Slug,
	}
	if episode.Match != nil {
		data.Title = media.Title + " – " + episode.Match.Stage + ": " + episode.Match.TeamA.Name + " vs. " + episode.Match.TeamB.Name
	}

	if v.Profile != nil {
		data.ProfileID = v.Profile.ID
		if progress, err := s.Store.GetEpisodeProgress(r.Context(), v.Profile.ID, episode.ID); err == nil {
			// Clamp in case parts were removed since the progress was saved.
			if progress.PartIndex < len(ids) {
				data.InitialPart = progress.PartIndex
				data.InitialProgress = progress.ProgressSeconds
			}
		} else if err != store.ErrNotFound {
			return err
		}
	}

	if data.NextHref, err = s.nextEpisodeHref(r, season, episode); err != nil {
		return err
	}

	return s.render(w, r, "watch.html", data)
}

func (s *Server) nextEpisodeHref(r *http.Request, season models.Season, episode models.Episode) (string, error) {
	next, err := s.Store.GetNextEpisode(r.Context(), season.ID, episode.EpisodeNumber)
	if err == store.ErrNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return "/watch/episode/" + next.ID, nil
}

func displayName(u auth.SessionUser) string {
	if u.Name != nil {
		return *u.Name
	}
	return "Profil"
}
