package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

type adminGenresData struct {
	PageData
	Genres []models.Genre
}

func (s *Server) AdminGenresPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	genres, err := s.Store.ListGenres(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_genres.html", adminGenresData{PageData: s.basePageData(w, r), Genres: genres})
}

func (s *Server) AdminGenreCreate(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}

	name, err := requireMinLen("Name", r.FormValue("name"), 2)
	if err != nil {
		return err
	}
	slug, err := requireSlug(r.FormValue("slug"))
	if err != nil {
		return err
	}

	if _, err := s.Store.CreateGenre(r.Context(), name, slug); err != nil {
		return err
	}

	redirect(w, r, "/admin/genres")
	return nil
}

func (s *Server) AdminGenreDelete(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	if err := r.ParseForm(); err != nil {
		return validationErr("Ungültige Formulardaten")
	}
	if err := s.checkCSRF(r); err != nil {
		return err
	}
	if err := s.Store.DeleteGenre(r.Context(), r.PathValue("id")); err != nil {
		return err
	}
	redirect(w, r, "/admin/genres")
	return nil
}
