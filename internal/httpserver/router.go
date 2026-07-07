package httpserver

import (
	"io/fs"
	"net/http"

	"github.com/pure991/streamflix/web"
)

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	static, err := fs.Sub(web.FS, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))

	// Public pages
	mux.HandleFunc("GET /{$}", s.handler((*Server).HomePage))
	mux.HandleFunc("GET /pricing", s.handler((*Server).PricingPage))
	mux.HandleFunc("GET /legal/impressum", s.handler((*Server).ImpressumPage))
	mux.HandleFunc("GET /legal/datenschutz", s.handler((*Server).DatenschutzPage))
	mux.HandleFunc("GET /login", s.handler((*Server).LoginPage))
	mux.HandleFunc("GET /register", s.handler((*Server).RegisterPage))
	mux.HandleFunc("GET /browse", s.handler((*Server).BrowsePage))
	mux.HandleFunc("GET /movie/{slug}", s.handler((*Server).MoviePage))
	mux.HandleFunc("GET /series/{slug}", s.handler((*Server).SeriesPage))
	mux.HandleFunc("GET /watch/episode/{id}", s.handler((*Server).WatchEpisodePage))
	mux.HandleFunc("GET /watch/{id}", s.handler((*Server).WatchMoviePage))

	// Authenticated pages
	mux.HandleFunc("GET /profiles", s.handler((*Server).ProfilesPage))
	mux.HandleFunc("GET /account", s.handler((*Server).AccountPage))
	mux.HandleFunc("GET /billing", s.handler((*Server).BillingPage))
	mux.HandleFunc("GET /my-list", s.handler((*Server).MyListPage))

	// Admin pages
	mux.HandleFunc("GET /admin", s.handler((*Server).AdminDashboardPage))
	mux.HandleFunc("GET /admin/media", s.handler((*Server).AdminMediaListPage))
	mux.HandleFunc("GET /admin/media/new", s.handler((*Server).AdminMediaNewFormPage))
	mux.HandleFunc("POST /admin/media", s.handler((*Server).AdminMediaCreate))
	mux.HandleFunc("GET /admin/media/{id}/edit", s.handler((*Server).AdminMediaEditFormPage))
	mux.HandleFunc("POST /admin/media/{id}/delete", s.handler((*Server).AdminMediaDelete))
	mux.HandleFunc("POST /admin/media/{id}", s.handler((*Server).AdminMediaUpdate))

	mux.HandleFunc("GET /admin/genres", s.handler((*Server).AdminGenresPage))
	mux.HandleFunc("POST /admin/genres/{id}/delete", s.handler((*Server).AdminGenreDelete))
	mux.HandleFunc("POST /admin/genres", s.handler((*Server).AdminGenreCreate))

	mux.HandleFunc("GET /admin/series", s.handler((*Server).AdminSeriesListPage))
	mux.HandleFunc("GET /admin/series/new", s.handler((*Server).AdminSeriesNewFormPage))
	mux.HandleFunc("POST /admin/series", s.handler((*Server).AdminSeriesCreate))
	mux.HandleFunc("GET /admin/series/{id}", s.handler((*Server).AdminSeriesDetailPage))
	mux.HandleFunc("POST /admin/series/{id}/delete", s.handler((*Server).AdminSeriesDelete))
	mux.HandleFunc("POST /admin/series/{id}/seasons", s.handler((*Server).AdminSeasonCreate))
	mux.HandleFunc("POST /admin/series/{id}", s.handler((*Server).AdminSeriesUpdate))

	mux.HandleFunc("GET /admin/seasons/{id}", s.handler((*Server).AdminSeasonDetailPage))
	mux.HandleFunc("POST /admin/seasons/{id}/delete", s.handler((*Server).AdminSeasonDelete))
	mux.HandleFunc("POST /admin/seasons/{id}/episodes", s.handler((*Server).AdminEpisodeCreate))
	mux.HandleFunc("POST /admin/seasons/{id}", s.handler((*Server).AdminSeasonUpdate))

	mux.HandleFunc("GET /admin/episodes/{id}/edit", s.handler((*Server).AdminEpisodeEditFormPage))
	mux.HandleFunc("POST /admin/episodes/{id}/delete", s.handler((*Server).AdminEpisodeDelete))
	mux.HandleFunc("POST /admin/episodes/{id}", s.handler((*Server).AdminEpisodeUpdate))

	mux.HandleFunc("GET /admin/users", s.handler((*Server).AdminUsersPage))
	mux.HandleFunc("POST /admin/users/{id}/role", s.handler((*Server).AdminUserToggleRole))
	mux.HandleFunc("GET /admin/subscriptions", s.handler((*Server).AdminSubscriptionsPage))
	mux.HandleFunc("GET /admin/analytics", s.handler((*Server).AdminAnalyticsPage))

	// Auth API
	mux.HandleFunc("POST /api/auth/register", s.handler((*Server).Register))
	mux.HandleFunc("POST /api/auth/login", s.handler((*Server).Login))
	mux.HandleFunc("POST /api/auth/logout", s.handler((*Server).Logout))

	// Billing API
	mux.HandleFunc("POST /api/billing/create-checkout-session", s.handler((*Server).CreateCheckoutSession))
	mux.HandleFunc("POST /api/billing/create-portal-session", s.handler((*Server).CreatePortalSession))
	mux.HandleFunc("POST /api/webhooks/stripe", s.handler((*Server).StripeWebhook))

	// Profiles API
	mux.HandleFunc("GET /api/profiles", s.handler((*Server).APIListProfiles))
	mux.HandleFunc("POST /api/profiles", s.handler((*Server).CreateProfile))
	mux.HandleFunc("PATCH /api/profiles/{id}", s.handler((*Server).UpdateProfile))
	mux.HandleFunc("DELETE /api/profiles/{id}", s.handler((*Server).DeleteProfile))

	// Watch API
	mux.HandleFunc("GET /api/watch/progress/{id}", s.handler((*Server).APIGetWatchProgressForMedia))
	mux.HandleFunc("POST /api/watch/progress", s.handler((*Server).APIPostWatchProgress))
	mux.HandleFunc("GET /api/watch/{id}", s.handler((*Server).APIGetWatch))

	// Watchlist API
	mux.HandleFunc("GET /api/watchlist", s.handler((*Server).APIListWatchlist))
	mux.HandleFunc("POST /api/watchlist", s.handler((*Server).AddToWatchlist))
	mux.HandleFunc("DELETE /api/watchlist/{id}", s.handler((*Server).RemoveFromWatchlist))

	// Public catalog API
	mux.HandleFunc("GET /api/media/{slug}", s.handler((*Server).APIGetMediaBySlug))
	mux.HandleFunc("GET /api/media", s.handler((*Server).APIListMedia))
	mux.HandleFunc("GET /api/genres", s.handler((*Server).APIListGenres))

	// Admin API (read-only + genre create, matching the original's API surface)
	mux.HandleFunc("GET /api/admin/stats", s.handler((*Server).APIAdminStats))
	mux.HandleFunc("GET /api/admin/subscriptions", s.handler((*Server).APIAdminSubscriptions))
	mux.HandleFunc("GET /api/admin/users", s.handler((*Server).APIAdminUsers))
	mux.HandleFunc("POST /api/admin/genres", s.handler((*Server).APIAdminGenreCreate))

	return mux
}
