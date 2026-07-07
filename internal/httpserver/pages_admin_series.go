package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

type adminSeriesListData struct {
	PageData
	Series []models.Series
}

func (s *Server) AdminSeriesListPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	list, err := s.Store.ListAllSeries(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_series_list.html", adminSeriesListData{PageData: s.basePageData(w, r), Series: list})
}

type adminSeriesNewData struct {
	PageData
	AvailableMedia []models.Media
}

func (s *Server) AdminSeriesNewFormPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	media, err := s.Store.ListSeriesMediaWithoutSeries(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_series_new.html", adminSeriesNewData{PageData: s.basePageData(w, r), AvailableMedia: media})
}

func (s *Server) AdminSeriesCreate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
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
	title, err := requireMinLen("Titel", r.FormValue("title"), 2)
	if err != nil {
		return err
	}
	description, err := requireMinLen("Beschreibung", r.FormValue("description"), 10)
	if err != nil {
		return err
	}

	series, err := s.Store.CreateSeries(r.Context(), mediaID, title, description)
	if err != nil {
		return err
	}

	redirect(w, r, "/admin/series/"+series.ID)
	return nil
}

type adminSeriesDetailData struct {
	PageData
	Series models.Series
	Media  models.Media
}

func (s *Server) AdminSeriesDetailPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	series, err := s.Store.GetSeriesByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	seasons, err := s.Store.ListSeasons(r.Context(), series.ID)
	if err != nil {
		return err
	}
	series.Seasons = seasons
	media, err := s.Store.GetMediaByID(r.Context(), series.MediaID)
	if err != nil {
		return err
	}

	return s.render(w, r, "admin_series_detail.html", adminSeriesDetailData{
		PageData: s.basePageData(w, r),
		Series:   series,
		Media:    media,
	})
}

func (s *Server) AdminSeriesUpdate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	title, err := requireMinLen("Titel", r.FormValue("title"), 2)
	if err != nil {
		return err
	}
	description, err := requireMinLen("Beschreibung", r.FormValue("description"), 10)
	if err != nil {
		return err
	}

	id := r.PathValue("id")
	if _, err := s.Store.UpdateSeries(r.Context(), id, title, description); err != nil {
		return err
	}

	redirect(w, r, "/admin/series/"+id)
	return nil
}

func (s *Server) AdminSeriesDelete(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}
	if err := s.Store.DeleteSeries(r.Context(), r.PathValue("id")); err != nil {
		return err
	}
	redirect(w, r, "/admin/series")
	return nil
}

func (s *Server) AdminSeasonCreate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	seriesID := r.PathValue("id")
	seasonNumber, err := coerceInt("seasonNumber", r.FormValue("seasonNumber"))
	if err != nil {
		return err
	}
	title, err := requireMinLen("Titel", r.FormValue("title"), 1)
	if err != nil {
		return err
	}

	if _, err := s.Store.CreateSeason(r.Context(), seriesID, seasonNumber, title); err != nil {
		return err
	}

	redirect(w, r, "/admin/series/"+seriesID)
	return nil
}

type adminSeasonDetailData struct {
	PageData
	Season models.Season
	Series models.Series
}

func (s *Server) AdminSeasonDetailPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	season, err := s.Store.GetSeasonByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	episodes, err := s.Store.ListEpisodes(r.Context(), season.ID)
	if err != nil {
		return err
	}
	season.Episodes = episodes
	series, err := s.Store.GetSeriesByID(r.Context(), season.SeriesID)
	if err != nil {
		return err
	}

	return s.render(w, r, "admin_season_detail.html", adminSeasonDetailData{
		PageData: s.basePageData(w, r),
		Season:   season,
		Series:   series,
	})
}

func (s *Server) AdminSeasonUpdate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	seasonNumber, err := coerceInt("seasonNumber", r.FormValue("seasonNumber"))
	if err != nil {
		return err
	}
	title, err := requireMinLen("Titel", r.FormValue("title"), 1)
	if err != nil {
		return err
	}

	if _, err := s.Store.UpdateSeason(r.Context(), id, seasonNumber, title); err != nil {
		return err
	}

	redirect(w, r, "/admin/seasons/"+id)
	return nil
}

func (s *Server) AdminSeasonDelete(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	season, err := s.Store.GetSeasonByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	if err := s.Store.DeleteSeason(r.Context(), season.ID); err != nil {
		return err
	}

	redirect(w, r, "/admin/series/"+season.SeriesID)
	return nil
}

func parseEpisodeForm(r *http.Request) (title, description string, episodeNumber, duration int, videoURL, thumbnailURL string, isPublished bool, err error) {
	if title, err = requireMinLen("Titel", r.FormValue("title"), 1); err != nil {
		return
	}
	if description, err = requireMinLen("Beschreibung", r.FormValue("description"), 10); err != nil {
		return
	}
	if episodeNumber, err = coerceInt("episodeNumber", r.FormValue("episodeNumber")); err != nil {
		return
	}
	if duration, err = coerceInt("duration", r.FormValue("duration")); err != nil {
		return
	}
	if videoURL, err = requireURL("videoUrl", r.FormValue("videoUrl")); err != nil {
		return
	}
	if thumbnailURL, err = requireURL("thumbnailUrl", r.FormValue("thumbnailUrl")); err != nil {
		return
	}
	isPublished = r.FormValue("isPublished") == "on"
	return
}

func (s *Server) AdminEpisodeCreate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	seasonID := r.PathValue("id")
	title, description, episodeNumber, duration, videoURL, thumbnailURL, isPublished, err := parseEpisodeForm(r)
	if err != nil {
		return err
	}

	if _, err := s.Store.CreateEpisode(r.Context(), seasonID, title, description, episodeNumber, duration, videoURL, thumbnailURL, isPublished); err != nil {
		return err
	}

	redirect(w, r, "/admin/seasons/"+seasonID)
	return nil
}

type adminEpisodeFormData struct {
	PageData
	Episode models.Episode
}

func (s *Server) AdminEpisodeEditFormPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	episode, err := s.Store.GetEpisodeByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_episode_form.html", adminEpisodeFormData{PageData: s.basePageData(w, r), Episode: episode})
}

func (s *Server) AdminEpisodeUpdate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	id := r.PathValue("id")
	title, description, episodeNumber, duration, videoURL, thumbnailURL, isPublished, err := parseEpisodeForm(r)
	if err != nil {
		return err
	}

	episode, err := s.Store.UpdateEpisode(r.Context(), id, title, description, episodeNumber, duration, videoURL, thumbnailURL, isPublished)
	if err != nil {
		return err
	}

	redirect(w, r, "/admin/seasons/"+episode.SeasonID)
	return nil
}

func (s *Server) AdminEpisodeDelete(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	episode, err := s.Store.GetEpisodeByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	if err := s.Store.DeleteEpisode(r.Context(), episode.ID); err != nil {
		return err
	}

	redirect(w, r, "/admin/seasons/"+episode.SeasonID)
	return nil
}
