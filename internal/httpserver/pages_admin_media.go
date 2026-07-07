package httpserver

import (
	"net/http"
	"strings"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

type adminMediaListData struct {
	PageData
	Media []models.Media
}

func (s *Server) AdminMediaListPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	media, err := s.Store.ListAllMedia(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_media_list.html", adminMediaListData{PageData: s.basePageData(w, r), Media: media})
}

type adminMediaFormData struct {
	PageData
	Media  *models.Media
	Genres []models.Genre
}

func (s *Server) AdminMediaNewFormPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	genres, err := s.Store.ListGenres(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_media_form.html", adminMediaFormData{PageData: s.basePageData(w, r), Genres: genres})
}

func (s *Server) AdminMediaEditFormPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	media, err := s.Store.GetMediaByID(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}
	genres, err := s.Store.ListGenres(r.Context())
	if err != nil {
		return err
	}
	media.Genres, err = s.Store.GetGenresForMedia(r.Context(), media.ID)
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_media_form.html", adminMediaFormData{PageData: s.basePageData(w, r), Media: &media, Genres: genres})
}

func parseMediaForm(r *http.Request) (store.MediaInput, error) {
	var in store.MediaInput
	var err error

	if in.Title, err = requireMinLen("Titel", r.FormValue("title"), 2); err != nil {
		return in, err
	}
	if in.Slug, err = requireSlug(r.FormValue("slug")); err != nil {
		return in, err
	}
	if in.Description, err = requireMinLen("Beschreibung", r.FormValue("description"), 10); err != nil {
		return in, err
	}
	typeVal, err := requireEnum("type", r.FormValue("type"), string(models.MediaMovie), string(models.MediaSeries))
	if err != nil {
		return in, err
	}
	in.Type = models.MediaType(typeVal)

	if in.ThumbnailURL, err = requireURL("thumbnailUrl", r.FormValue("thumbnailUrl")); err != nil {
		return in, err
	}
	if in.BannerURL, err = requireURL("bannerUrl", r.FormValue("bannerUrl")); err != nil {
		return in, err
	}
	if in.TrailerURL, err = optionalURL("trailerUrl", r.FormValue("trailerUrl")); err != nil {
		return in, err
	}
	if in.VideoURL, err = optionalURL("videoUrl", r.FormValue("videoUrl")); err != nil {
		return in, err
	}
	if in.ReleaseYear, err = coerceInt("releaseYear", r.FormValue("releaseYear")); err != nil {
		return in, err
	}
	if err := requireIntMin("releaseYear", in.ReleaseYear, 1900); err != nil {
		return in, err
	}
	if in.Duration, err = coerceOptionalInt(r.FormValue("duration")); err != nil {
		return in, err
	}
	if in.AgeRating, err = requireMinLen("Altersfreigabe", r.FormValue("ageRating"), 1); err != nil {
		return in, err
	}
	in.Language = r.FormValue("language")
	if in.Language == "" {
		in.Language = "Deutsch"
	}
	in.Tags = []string{}
	if tags := strings.TrimSpace(r.FormValue("tags")); tags != "" {
		for _, t := range strings.Split(tags, ",") {
			if t = strings.TrimSpace(t); t != "" {
				in.Tags = append(in.Tags, t)
			}
		}
	}
	in.IsPublished = r.FormValue("isPublished") == "on"
	in.GenreIDs = r.Form["genreIds"]

	return in, nil
}

func (s *Server) AdminMediaCreate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	in, err := parseMediaForm(r)
	if err != nil {
		return err
	}

	if _, err := s.Store.CreateMedia(r.Context(), in); err != nil {
		return err
	}

	redirect(w, r, "/admin/media")
	return nil
}

func (s *Server) AdminMediaUpdate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	in, err := parseMediaForm(r)
	if err != nil {
		return err
	}

	if _, err := s.Store.UpdateMedia(r.Context(), r.PathValue("id"), in); err != nil {
		return err
	}

	redirect(w, r, "/admin/media")
	return nil
}

func (s *Server) AdminMediaDelete(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}
	if err := s.Store.DeleteMedia(r.Context(), r.PathValue("id")); err != nil {
		return err
	}
	redirect(w, r, "/admin/media")
	return nil
}
