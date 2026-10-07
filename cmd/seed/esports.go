package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

// esportsCatalog is the curated Counter-Strike VOD archive: every tournament
// from its first officially uploaded match (qualifiers included) to the
// grand final. Every video is an upload on the organizer's own official
// YouTube channel; see docs/esports-vods.md for sources, how the catalog
// was built and what was left out.
//
//go:embed esports_catalog.json
var esportsCatalogJSON []byte

type catalog struct {
	Games       []catalogGame       `json:"games"`
	Teams       []catalogTeam       `json:"teams"`
	Tournaments []catalogTournament `json:"tournaments"`
}

type catalogGame struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type catalogTeam struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Short string `json:"short,omitempty"`
}

type catalogTournament struct {
	Slug      string         `json:"slug"`
	Title     string         `json:"title"`
	Year      int            `json:"year"`
	Location  string         `json:"location"`
	Game      string         `json:"game"`
	Category  string         `json:"category"` // major, international, dach-liga
	Organizer string         `json:"organizer"`
	Channels  []string       `json:"channels"`
	Language  string         `json:"language"`
	Note      string         `json:"note,omitempty"`
	Stages    []catalogStage `json:"stages"`
}

type catalogStage struct {
	Title    string           `json:"title"`
	Episodes []catalogEpisode `json:"episodes"`
}

// catalogEpisode is a match when both teams are set, otherwise a plain
// episode (a full broadcast day, or a pairing whose teams couldn't be
// identified unambiguously, e.g. WCG national teams).
type catalogEpisode struct {
	Title    string   `json:"title"`
	TeamA    string   `json:"teamA,omitempty"`
	TeamB    string   `json:"teamB,omitempty"`
	Stage    string   `json:"stage,omitempty"`
	BestOf   int      `json:"bestOf,omitempty"`
	Score    []int    `json:"score,omitempty"`
	PlayedOn string   `json:"playedOn,omitempty"`
	Videos   []string `json:"videos"`
	Minutes  int      `json:"minutes"`
}

func loadCatalog() (catalog, error) {
	var c catalog
	err := json.Unmarshal(esportsCatalogJSON, &c)
	return c, err
}

func seedEsports(ctx context.Context, st *store.Store) error {
	cat, err := loadCatalog()
	if err != nil {
		return fmt.Errorf("parse esports catalog: %w", err)
	}

	genre, err := ensureGenre(ctx, st, "eSports", "esports")
	if err != nil {
		return err
	}

	games := map[string]models.Game{}
	labels := map[string]string{}
	for _, g := range cat.Games {
		game, err := st.EnsureGame(ctx, g.Name, g.Slug)
		if err != nil {
			return err
		}
		games[g.Slug] = game
		labels[g.Slug] = g.Label
	}

	teams := map[string]models.Team{}
	for _, t := range cat.Teams {
		in := store.TeamInput{Name: t.Name, Slug: t.Slug}
		if t.Short != "" {
			short := t.Short
			in.ShortName = &short
		}
		team, err := st.EnsureTeam(ctx, in)
		if err != nil {
			return err
		}
		teams[t.Slug] = team
	}

	for _, t := range cat.Tournaments {
		created, added, err := syncTournament(ctx, st, t, genre.ID, games[t.Game], labels[t.Game], teams)
		if err != nil {
			return fmt.Errorf("%s: %w", t.Slug, err)
		}
		switch {
		case created:
			log.Printf("created eSports tournament %s (%d episodes)", t.Title, added)
		case added > 0:
			log.Printf("updated eSports tournament %s (+%d episodes)", t.Title, added)
		}
	}
	return nil
}

// syncTournament creates the tournament if needed and then makes sure every
// stage and episode of the catalog exists in catalog order. It never
// deletes or overwrites anything, so admin edits and watch progress on
// existing episodes survive re-running the seed.
func syncTournament(ctx context.Context, st *store.Store, t catalogTournament, genreID string, game models.Game, gameLabel string, teams map[string]models.Team) (created bool, added int, err error) {
	poster := "/static/img/esports/" + t.Slug + "-poster.svg"

	media, err := st.GetMediaBySlug(ctx, t.Slug, false)
	if err == store.ErrNotFound {
		media, err = createTournamentMedia(ctx, st, t, genreID, gameLabel, poster)
		created = true
	}
	if err != nil {
		return created, 0, err
	}
	// Earlier seeds only imported the grand final and described the title
	// that way; refresh those generated texts, but never an admin's own.
	if !created && strings.HasPrefix(media.Description, "Das ") && strings.Contains(media.Description, "als offizielles VOD von") {
		if err := st.SetMediaDescription(ctx, media.ID, tournamentDescription(t, gameLabel)); err != nil {
			return created, 0, err
		}
	}

	series, err := st.GetSeriesByMediaID(ctx, media.ID)
	if err == store.ErrNotFound {
		if series, err = st.CreateSeries(ctx, media.ID, t.Title, "Alle offiziell veröffentlichten Matches von "+t.Title+"."); err != nil {
			return created, 0, err
		}
	} else if err != nil {
		return created, 0, err
	}

	// Park existing seasons on high numbers so stages can be inserted in
	// front of them (season_number is unique per series).
	existingSeasons := map[string]models.Season{}
	for _, s := range series.Seasons {
		existingSeasons[s.Title] = s
		if err := st.SetSeasonNumber(ctx, s.ID, 1000+s.SeasonNumber); err != nil {
			return created, 0, err
		}
	}

	number := 1
	for _, stage := range t.Stages {
		season, ok := existingSeasons[stage.Title]
		if ok {
			delete(existingSeasons, stage.Title)
			if err := st.SetSeasonNumber(ctx, season.ID, number); err != nil {
				return created, added, err
			}
		} else if season, err = st.CreateSeason(ctx, series.ID, number, stage.Title); err != nil {
			return created, added, err
		}
		number++

		n, err := syncStageEpisodes(ctx, st, season, stage, game, teams, poster)
		added += n
		if err != nil {
			return created, added, err
		}
	}

	// Seasons that aren't in the catalog (e.g. added by an admin) go last.
	for _, s := range series.Seasons {
		if _, leftover := existingSeasons[s.Title]; leftover {
			if err := st.SetSeasonNumber(ctx, s.ID, number); err != nil {
				return created, added, err
			}
			number++
		}
	}
	return created, added, nil
}

func syncStageEpisodes(ctx context.Context, st *store.Store, season models.Season, stage catalogStage, game models.Game, teams map[string]models.Team, poster string) (int, error) {
	existing := map[string]models.Episode{}
	for _, e := range season.Episodes {
		existing[e.Title] = e
		if err := st.SetEpisodeNumber(ctx, e.ID, 1000+e.EpisodeNumber); err != nil {
			return 0, err
		}
	}

	added, number := 0, 1
	for _, ce := range stage.Episodes {
		if len(ce.Videos) == 0 {
			continue
		}
		urls := make([]string, len(ce.Videos))
		for i, id := range ce.Videos {
			urls[i] = "https://www.youtube.com/watch?v=" + id
		}

		episode, ok := existing[ce.Title]
		if ok {
			delete(existing, ce.Title)
			if err := st.SetEpisodeNumber(ctx, episode.ID, number); err != nil {
				return added, err
			}
		} else {
			minutes := ce.Minutes
			if minutes <= 0 {
				minutes = 1
			}
			var err error
			episode, err = st.CreateEpisode(ctx, season.ID, ce.Title, matchEpisodeDescription(ce, teams), number, minutes, urls[0], poster, true)
			if err != nil {
				return added, err
			}
			added++
		}
		number++

		teamA, okA := teams[ce.TeamA]
		teamB, okB := teams[ce.TeamB]
		if !okA || !okB || episode.Match != nil {
			continue
		}
		in := store.MatchInput{
			TeamAID:        teamA.ID,
			TeamBID:        teamB.ID,
			Stage:          ce.Stage,
			BestOf:         ce.BestOf,
			ExtraVideoURLs: urls[1:],
		}
		if in.Stage == "" {
			in.Stage = stage.Title
		}
		if game.ID != "" {
			gameID := game.ID
			in.GameID = &gameID
		}
		if ce.PlayedOn != "" {
			d, err := time.Parse("2006-01-02", ce.PlayedOn)
			if err != nil {
				return added, err
			}
			in.PlayedOn = &d
		}
		if len(ce.Score) == 2 {
			a, b := ce.Score[0], ce.Score[1]
			in.ScoreA, in.ScoreB = &a, &b
		}
		if err := st.UpsertMatch(ctx, episode.ID, in); err != nil {
			return added, err
		}
	}

	for _, e := range season.Episodes {
		if _, leftover := existing[e.Title]; leftover {
			if err := st.SetEpisodeNumber(ctx, e.ID, number); err != nil {
				return added, err
			}
			number++
		}
	}
	return added, nil
}

// matchEpisodeDescription must stay spoiler-free: no winner, no score.
func matchEpisodeDescription(ce catalogEpisode, teams map[string]models.Team) string {
	a, okA := teams[ce.TeamA]
	b, okB := teams[ce.TeamB]
	if okA && okB {
		return ce.Title + ": " + a.Name + " gegen " + b.Name + ". Offizielles VOD des Veranstalters."
	}
	return ce.Title + ". Offizielles VOD des Veranstalters."
}

// tournamentDescription must stay spoiler-free: no winner, no score.
func tournamentDescription(t catalogTournament, gameLabel string) string {
	return "Alle offiziell veröffentlichten Matches von " + t.Title + " (" + gameLabel + ", " + t.Location +
		") bis zum Grand Final als VODs von " + t.Organizer + ". Ergebnisse bleiben verborgen, solange der Spoiler-Schutz aktiv ist."
}

func createTournamentMedia(ctx context.Context, st *store.Store, t catalogTournament, genreID, gameLabel, poster string) (models.Media, error) {
	media, err := st.CreateMedia(ctx, store.MediaInput{
		Title:        t.Title,
		Slug:         t.Slug,
		Description:  tournamentDescription(t, gameLabel),
		Type:         models.MediaSeries,
		ThumbnailURL: poster,
		BannerURL:    "/static/img/esports/" + t.Slug + "-banner.svg",
		ReleaseYear:  t.Year,
		AgeRating:    "USK 18",
		Language:     t.Language,
		Tags:         []string{"esports", "counter-strike", t.Game, t.Category},
		IsPublished:  true,
		GenreIDs:     []string{genreID},
	})
	if err != nil {
		return media, err
	}

	holder := t.Organizer + " (offizielle YouTube-Kanäle)"
	var channel *string
	if len(t.Channels) > 0 {
		channel = &t.Channels[0]
	}
	notes := "Offizielle, einbettbare YouTube-Videos des Veranstalters. Wiedergabe ausschließlich über den YouTube-Embed-Player, ohne Abo-Pflicht (YouTube API Developer Policies)."
	for _, ch := range t.Channels {
		notes += " Kanal: " + ch + "."
	}
	if t.Note != "" {
		notes += " " + t.Note + "."
	}
	_, err = st.CreateRightsInfo(ctx, media.ID, store.RightsInfoInput{
		Source:        models.RightsOfficialEmbed,
		RightsHolder:  &holder,
		LicenseDocURL: channel,
		ValidFrom:     time.Date(t.Year, 1, 1, 0, 0, 0, 0, time.UTC),
		Notes:         &notes,
	})
	return media, err
}
