package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

type browseData struct {
	PageData
	Hero     *models.Media
	Continue []models.Media
	Popular  []models.Media
	Movies   []models.Media
	Series   []models.Media
	ForYou   []models.Media
	Esports  []models.Media

	FollowedMatches []models.TeamMatch
	HideSpoilers    bool
}

func hasTag(m models.Media, tag string) bool {
	for _, t := range m.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (s *Server) BrowsePage(w http.ResponseWriter, r *http.Request) error {
	all, err := s.Store.ListPublishedMedia(r.Context())
	if err != nil {
		return err
	}

	data := browseData{PageData: s.basePageData(w, r), Continue: all, Popular: all}
	if len(all) > 0 {
		data.Hero = &all[0]
	}
	for _, m := range all {
		if m.Type == models.MediaMovie {
			data.Movies = append(data.Movies, m)
		} else {
			data.Series = append(data.Series, m)
		}
		if hasTag(m, "featured") {
			data.ForYou = append(data.ForYou, m)
		}
	}

	if data.Esports, err = s.Store.ListEsportsMedia(r.Context()); err != nil {
		return err
	}
	v, err := s.currentViewer(r.Context(), r)
	if err != nil {
		return err
	}
	data.HideSpoilers = v.HideSpoilers
	if v.Profile != nil {
		if data.FollowedMatches, err = s.Store.ListMatchesForFollowedTeams(r.Context(), v.Profile.ID); err != nil {
			return err
		}
	}

	return s.render(w, r, "browse.html", data)
}
