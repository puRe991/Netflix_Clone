package httpserver

import (
	"net/http"

	"github.com/pure991/streamflix/internal/models"
)

type adminStatsData struct {
	PageData
	Stats models.Stats
}

func (s *Server) AdminDashboardPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	stats, err := s.Store.GetStats(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_dashboard.html", adminStatsData{PageData: s.basePageData(w, r), Stats: stats})
}

func (s *Server) AdminAnalyticsPage(w http.ResponseWriter, r *http.Request) error {
	if _, err := s.Sessions.RequireAdmin(r); err != nil {
		return err
	}
	stats, err := s.Store.GetStats(r.Context())
	if err != nil {
		return err
	}
	return s.render(w, r, "admin_analytics.html", adminStatsData{PageData: s.basePageData(w, r), Stats: stats})
}
