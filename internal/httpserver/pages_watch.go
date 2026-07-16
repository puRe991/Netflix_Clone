package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/auth"
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

	nextHref := ""
	if next, err := s.Store.GetNextEpisode(r.Context(), season.ID, episode.EpisodeNumber); err == nil {
		nextHref = "/watch/episode/" + next.ID
	} else if err != store.ErrNotFound {
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

func displayName(u auth.SessionUser) string {
	if u.Name != nil {
		return *u.Name
	}
	return "Profil"
}
