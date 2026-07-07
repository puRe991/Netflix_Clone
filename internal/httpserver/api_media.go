package httpserver

import "net/http"

func (s *Server) APIListMedia(w http.ResponseWriter, r *http.Request) error {
	media, err := s.Store.ListPublishedMedia(r.Context())
	if err != nil {
		return err
	}
	for i := range media {
		genres, err := s.Store.GetGenresForMedia(r.Context(), media[i].ID)
		if err != nil {
			return err
		}
		media[i].Genres = genres
	}
	JSON(w, http.StatusOK, media)
	return nil
}

func (s *Server) APIGetMediaBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := r.PathValue("slug")
	media, err := s.Store.GetMediaBySlug(r.Context(), slug, false)
	if err != nil {
		return err
	}
	JSON(w, http.StatusOK, media)
	return nil
}

func (s *Server) APIListGenres(w http.ResponseWriter, r *http.Request) error {
	genres, err := s.Store.ListGenres(r.Context())
	if err != nil {
		return err
	}
	JSON(w, http.StatusOK, genres)
	return nil
}
