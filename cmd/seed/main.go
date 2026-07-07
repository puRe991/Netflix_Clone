// Command seed populates the database with a demo admin account, genres,
// a movie, and a series with one season of episodes — mirroring the
// original prisma/seed.ts demo dataset.
package main

import (
	"context"
	"log"
	"strconv"

	"github.com/pure991/streamflix/internal/auth"
	"github.com/pure991/streamflix/internal/config"
	"github.com/pure991/streamflix/internal/db"
	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

const sampleVideo = "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		return err
	}

	st := store.New(pool)

	drama, err := ensureGenre(ctx, st, "Drama", "drama")
	if err != nil {
		return err
	}
	sciFi, err := ensureGenre(ctx, st, "Sci-Fi", "sci-fi")
	if err != nil {
		return err
	}
	doc, err := ensureGenre(ctx, st, "Dokumentation", "dokumentation")
	if err != nil {
		return err
	}

	passwordHash, err := auth.HashPassword("StreamFlix123!")
	if err != nil {
		return err
	}

	adminEmail := "admin@streamflix.local"
	if _, err := st.GetUserByEmail(ctx, adminEmail); err == store.ErrNotFound {
		adminID, err := st.CreateUserWithProfile(ctx, adminEmail, passwordHash, strPtr("Admin"), "Admin")
		if err != nil {
			return err
		}
		if err := st.SetUserRole(ctx, adminID, models.RoleAdmin); err != nil {
			return err
		}
		if err := st.SetSubscriptionStatus(ctx, adminID, models.SubStatusActive); err != nil {
			return err
		}
		log.Println("created admin user admin@streamflix.local / StreamFlix123!")
	} else if err != nil {
		return err
	} else {
		log.Println("admin user already exists, skipping")
	}

	movie, err := st.GetMediaBySlug(ctx, "ocean-city", false)
	if err == store.ErrNotFound {
		movie, err = st.CreateMedia(ctx, storeMediaInputMovie(drama.ID, doc.ID))
		if err != nil {
			return err
		}
		log.Println("created movie Ocean City")
	} else if err != nil {
		return err
	}
	_ = movie

	seriesMedia, err := st.GetMediaBySlug(ctx, "red-orbit", false)
	if err == store.ErrNotFound {
		seriesMedia, err = st.CreateMedia(ctx, storeMediaInputSeries(sciFi.ID))
		if err != nil {
			return err
		}
		log.Println("created series media Red Orbit")

		series, err := st.CreateSeries(ctx, seriesMedia.ID, "Red Orbit", "Staffel 1 der Forschungsmission.")
		if err != nil {
			return err
		}
		season, err := st.CreateSeason(ctx, series.ID, 1, "Mission Start")
		if err != nil {
			return err
		}
		for n := 1; n <= 3; n++ {
			if _, err := st.CreateEpisode(ctx, season.ID,
				episodeTitle(n), episodeDescription(n), n, 596, sampleVideo,
				"https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?auto=format&fit=crop&w=800&q=80", true); err != nil {
				return err
			}
		}
		log.Println("created series Red Orbit with 3 episodes")
	} else if err != nil {
		return err
	}

	log.Println("seed complete")
	return nil
}

func ensureGenre(ctx context.Context, st *store.Store, name, slug string) (models.Genre, error) {
	genres, err := st.ListGenres(ctx)
	if err != nil {
		return models.Genre{}, err
	}
	for _, g := range genres {
		if g.Slug == slug {
			return g, nil
		}
	}
	return st.CreateGenre(ctx, name, slug)
}

func storeMediaInputMovie(dramaID, docID string) store.MediaInput {
	return store.MediaInput{
		Title:        "Ocean City",
		Slug:         "ocean-city",
		Description:  "Ein lizenzierter Demo-Film über Mut, Gemeinschaft und den Neubeginn an der Küste.",
		Type:         models.MediaMovie,
		ThumbnailURL: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=600&q=80",
		BannerURL:    "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1800&q=80",
		VideoURL:     strPtr(sampleVideo),
		TrailerURL:   strPtr(sampleVideo),
		ReleaseYear:  2026,
		Duration:     intPtr(10),
		AgeRating:    "FSK 6",
		Language:     "Deutsch",
		Tags:         []string{"featured", "coast", "family"},
		IsPublished:  true,
		GenreIDs:     []string{dramaID, docID},
	}
}

func storeMediaInputSeries(sciFiID string) store.MediaInput {
	return store.MediaInput{
		Title:        "Red Orbit",
		Slug:         "red-orbit",
		Description:  "Eine legale Sci-Fi-Serie über eine Forschungscrew auf dem Weg zu einem roten Planeten.",
		Type:         models.MediaSeries,
		ThumbnailURL: "https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?auto=format&fit=crop&w=600&q=80",
		BannerURL:    "https://images.unsplash.com/photo-1462331940025-496dfbfc7564?auto=format&fit=crop&w=1800&q=80",
		TrailerURL:   strPtr(sampleVideo),
		ReleaseYear:  2026,
		AgeRating:    "FSK 12",
		Language:     "Deutsch",
		Tags:         []string{"featured", "space"},
		IsPublished:  true,
		GenreIDs:     []string{sciFiID},
	}
}

func episodeTitle(n int) string { return "Episode " + strconv.Itoa(n) }
func episodeDescription(n int) string {
	return "Kapitel " + strconv.Itoa(n) + " der Red-Orbit-Mission."
}

func strPtr(s string) *string { return &s }
func intPtr(n int) *int       { return &n }
