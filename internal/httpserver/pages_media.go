package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

type movieData struct {
	PageData
	Media     models.Media
	CanStream bool
}

func (s *Server) MoviePage(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")
	media, err := s.Store.GetMediaBySlug(r.Context(), slug, false)
	if err != nil {
		return err
	}

	user := s.Sessions.OptionalUser(r)
	canStream, err := s.canStreamFullContent(r.Context(), user)
	if err != nil {
		return err
	}

	return s.render(w, r, "movie.html", movieData{
		PageData:  s.basePageData(w, r),
		Media:     media,
		CanStream: canStream,
	})
}

type seriesData struct {
	PageData
	Media        models.Media
	HideSpoilers bool
	HasMatches   bool
	Path         string
}

func (s *Server) SeriesPage(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")
	media, err := s.Store.GetMediaBySlug(r.Context(), slug, false)
	if err != nil {
		return err
	}
	if media.Type != models.MediaSeries || media.Series == nil {
		return store.ErrNotFound
	}

	v, err := s.currentViewer(r.Context(), r)
	if err != nil {
		return err
	}

	hasMatches := false
	for _, season := range media.Series.Seasons {
		for _, e := range season.Episodes {
			if e.Match != nil {
				hasMatches = true
			}
		}
	}

	return s.render(w, r, "series.html", seriesData{
		PageData:     s.basePageData(w, r),
		Media:        media,
		HideSpoilers: v.HideSpoilers,
		HasMatches:   hasMatches,
		Path:         r.URL.Path,
	})
}
