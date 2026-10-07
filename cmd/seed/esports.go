package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

// Counter-Strike Major grand finals, each as a tournament (SERIES) with a
// "Playoffs" season and the grand final as its match episode.
//
// Every video is an upload on the tournament organizer's own official
// YouTube channel, checked via YouTube oEmbed to exist and allow embedding
// (see docs/esports-vods.md for the full list, the channels and the Majors
// without an official VOD). Fan re-uploads are deliberately not used.
//
// Team order follows the official video titles (bracket order), never
// "winner first", so the listing itself doesn't give results away.

type seedTeam struct {
	Slug, Name, Short string
}

var esportsTeams = []seedTeam{
	{"astralis", "Astralis", "AST"},
	{"avangar", "AVANGAR", "AVG"},
	{"ence", "ENCE", "ENCE"},
	{"envyus", "Team EnVyUs", "nV"},
	{"falcons", "Team Falcons", "FLC"},
	{"faze", "FaZe Clan", "FaZe"},
	{"fnatic", "Fnatic", "FNC"},
	{"furia", "FURIA", "FUR"},
	{"g2", "G2 Esports", "G2"},
	{"gambit", "Gambit Esports", "GMB"},
	{"gamerlegion", "GamerLegion", "GL"},
	{"heroic", "Heroic", "HER"},
	{"immortals", "Immortals", "IMT"},
	{"liquid", "Team Liquid", "TL"},
	{"luminosity", "Luminosity Gaming", "LG"},
	{"natus-vincere", "Natus Vincere", "NAVI"},
	{"nip", "Ninjas in Pyjamas", "NiP"},
	{"outsiders", "Outsiders", "OUT"},
	{"sk-gaming", "SK Gaming", "SK"},
	{"spirit", "Team Spirit", "SPR"},
	{"the-mongolz", "The MongolZ", "MGLZ"},
	{"virtus-pro", "Virtus.pro", "VP"},
	{"vitality", "Team Vitality", "VIT"},
}

type seedMajor struct {
	Slug        string
	Title       string
	Year        int
	Location    string
	Game        string // "csgo" or "cs2"
	Organizer   string
	ChannelURL  string
	PlayedOn    string
	TeamA       string
	TeamB       string
	BestOf      int
	ScoreA      int
	ScoreB      int
	DurationMin int
	VideoIDs    []string
	Note        string
}

var csMajors = []seedMajor{
	{"ems-one-katowice-2014", "EMS One Katowice 2014", 2014, "Katowice", "csgo", "ESL", "https://www.youtube.com/@esl",
		"2014-03-16", "virtus-pro", "nip", 3, 2, 0, 150, []string{"do4FUYUHmUM"}, ""},
	{"esl-one-cologne-2014", "ESL One Cologne 2014", 2014, "Köln", "csgo", "ESL", "https://www.youtube.com/@ESLCS",
		"2014-08-17", "nip", "fnatic", 3, 2, 1, 215, []string{"-5laY7bjono"}, "ESL Classics: komplette Übertragung des Finales"},
	{"esl-one-katowice-2015", "ESL One Katowice 2015", 2015, "Katowice", "csgo", "ESL", "https://www.youtube.com/@ESLCS",
		"2015-03-15", "fnatic", "nip", 3, 2, 1, 200, []string{"cqmdRLNRZHA"}, "ESL Classics: komplette Übertragung des Finales"},
	{"esl-one-cologne-2015", "ESL One Cologne 2015", 2015, "Köln", "csgo", "ESL", "https://www.youtube.com/@ESLArchives",
		"2015-08-23", "fnatic", "envyus", 3, 2, 0, 142, []string{"02I5vVxlJhU", "_Ef6cXH1oC8"}, ""},
	{"mlg-major-columbus-2016", "MLG Major Championship: Columbus 2016", 2016, "Columbus", "csgo", "Major League Gaming", "https://www.youtube.com/@MajorLeagueGaming",
		"2016-04-03", "natus-vincere", "luminosity", 3, 0, 2, 203, []string{"7nWxZin959o"}, ""},
	{"esl-one-cologne-2016", "ESL One Cologne 2016", 2016, "Köln", "csgo", "ESL", "https://www.youtube.com/@ESLArchives",
		"2016-07-10", "sk-gaming", "liquid", 3, 2, 0, 145, []string{"AK5RY8W-AcM", "TC6N1zHvXpU"}, ""},
	{"pgl-major-krakow-2017", "PGL Major Kraków 2017", 2017, "Kraków", "csgo", "PGL", "https://www.youtube.com/@pgl",
		"2017-07-23", "gambit", "immortals", 3, 2, 1, 226, []string{"ttbe_A8Ee50", "sIQ1Eh11Quk", "w2MdCEZSBtw"}, ""},
	{"faceit-major-london-2018", "FACEIT Major: London 2018", 2018, "London", "csgo", "FACEIT", "https://www.youtube.com/@FACEIT",
		"2018-09-23", "natus-vincere", "astralis", 3, 0, 2, 171, []string{"z7IDPdmYw0A", "9CcrjvC4SkE"}, ""},
	{"iem-katowice-major-2019", "IEM Katowice Major 2019", 2019, "Katowice", "csgo", "ESL", "https://www.youtube.com/@ESLArchives",
		"2019-03-03", "ence", "astralis", 3, 0, 2, 151, []string{"kgitmggEgrA", "A_wk85oA0Rc"}, ""},
	{"starladder-major-berlin-2019", "StarLadder Major: Berlin 2019", 2019, "Berlin", "csgo", "StarLadder", "https://www.youtube.com/@starladder_cs_highlights",
		"2019-09-08", "avangar", "astralis", 3, 0, 2, 188, []string{"c2uQV59-tbI", "hkVQXKCVSFQ"}, ""},
	{"pgl-major-stockholm-2021", "PGL Major Stockholm 2021", 2021, "Stockholm", "csgo", "PGL", "https://www.youtube.com/@pgl",
		"2021-11-07", "g2", "natus-vincere", 3, 0, 2, 146, []string{"UKYnAujcFac", "pogK6a25LO0"}, ""},
	{"pgl-major-antwerp-2022", "PGL Major Antwerp 2022", 2022, "Antwerpen", "csgo", "PGL", "https://www.youtube.com/@pgl",
		"2022-05-22", "faze", "natus-vincere", 3, 2, 0, 136, []string{"O9u1FQKZHVw", "CMBAVBkC8Zw"}, ""},
	{"iem-rio-major-2022", "IEM Rio Major 2022", 2022, "Rio de Janeiro", "csgo", "ESL", "https://www.youtube.com/@ESLArchives",
		"2022-11-13", "outsiders", "heroic", 3, 2, 0, 113, []string{"9pscaxA1QBs", "oVH782n-32U"}, ""},
	{"blast-paris-major-2023", "BLAST.tv Paris Major 2023", 2023, "Paris", "csgo", "BLAST", "https://www.youtube.com/@BLASTPremier",
		"2023-05-21", "vitality", "gamerlegion", 3, 2, 0, 410, []string{"y6Y0l1Y-PTo"}, "Komplette Übertragung des Finaltags inkl. Vorberichterstattung"},
	{"pgl-major-copenhagen-2024", "PGL CS2 Major Copenhagen 2024", 2024, "Kopenhagen", "cs2", "PGL", "https://www.youtube.com/@pgl",
		"2024-03-31", "faze", "natus-vincere", 3, 1, 2, 112, []string{"09N-bQbLrx8", "L9eBTSwDwaM", "vdXveK6G70g"}, ""},
	{"perfect-world-shanghai-major-2024", "Perfect World Shanghai Major 2024", 2024, "Shanghai", "cs2", "Perfect World / PGL", "https://www.youtube.com/@pgl",
		"2024-12-15", "spirit", "faze", 3, 2, 1, 413, []string{"Vv7eXxl--LA"}, "Komplette Übertragung inkl. Showmatch vor dem Finale"},
	{"blast-austin-major-2025", "BLAST.tv Austin Major 2025", 2025, "Austin", "cs2", "BLAST", "https://www.youtube.com/@BLASTPremier",
		"2025-06-22", "vitality", "the-mongolz", 3, 2, 1, 382, []string{"meijWUaxNYg"}, "Komplette Übertragung des Finaltags"},
	{"starladder-budapest-major-2025", "StarLadder Budapest Major 2025", 2025, "Budapest", "cs2", "StarLadder", "https://www.youtube.com/@starladder_cs",
		"2025-12-14", "vitality", "faze", 5, 3, 1, 169, []string{"9qqiaW7sjD0"}, "Erstes Major-Finale im Best-of-5-Format"},
	{"iem-cologne-major-2026", "IEM Cologne Major 2026", 2026, "Köln", "cs2", "ESL", "https://www.youtube.com/@ESLArchives",
		"2026-06-21", "furia", "falcons", 5, 0, 3, 164, []string{"p6U28DEBn8A"}, ""},
}

func seedEsports(ctx context.Context, st *store.Store) error {
	genre, err := ensureGenre(ctx, st, "eSports", "esports")
	if err != nil {
		return err
	}

	games := map[string]models.Game{}
	for slug, name := range map[string]string{"csgo": "Counter-Strike: Global Offensive", "cs2": "Counter-Strike 2"} {
		g, err := st.EnsureGame(ctx, name, slug)
		if err != nil {
			return err
		}
		games[slug] = g
	}

	teams := map[string]models.Team{}
	for _, t := range esportsTeams {
		short := t.Short
		team, err := st.EnsureTeam(ctx, store.TeamInput{Name: t.Name, Slug: t.Slug, ShortName: &short})
		if err != nil {
			return err
		}
		teams[t.Slug] = team
	}

	for _, mj := range csMajors {
		if _, err := st.GetMediaBySlug(ctx, mj.Slug, false); err == nil {
			continue
		} else if err != store.ErrNotFound {
			return err
		}
		if err := seedMajorFinal(ctx, st, mj, genre.ID, games[mj.Game], teams); err != nil {
			return err
		}
		log.Printf("created eSports tournament %s", mj.Title)
	}
	return nil
}

func seedMajorFinal(ctx context.Context, st *store.Store, mj seedMajor, genreID string, game models.Game, teams map[string]models.Team) error {
	gameLabel := "CS:GO"
	if mj.Game == "cs2" {
		gameLabel = "CS2"
	}
	poster := "/static/img/esports/" + mj.Slug + "-poster.svg"

	// Descriptions must stay spoiler-free: no winner, no score.
	description := "Das " + gameLabel + "-Major " + mj.Title + " in " + mj.Location + ": das Grand Final als offizielles VOD von " +
		mj.Organizer + ". Ergebnisse bleiben verborgen, solange der Spoiler-Schutz aktiv ist."

	media, err := st.CreateMedia(ctx, store.MediaInput{
		Title:        mj.Title,
		Slug:         mj.Slug,
		Description:  description,
		Type:         models.MediaSeries,
		ThumbnailURL: poster,
		BannerURL:    "/static/img/esports/" + mj.Slug + "-banner.svg",
		ReleaseYear:  mj.Year,
		AgeRating:    "USK 18",
		Language:     "Englisch",
		Tags:         []string{"esports", "counter-strike", mj.Game, "major"},
		IsPublished:  true,
		GenreIDs:     []string{genreID},
	})
	if err != nil {
		return err
	}

	holder := mj.Organizer + " (offizieller YouTube-Kanal)"
	channel := mj.ChannelURL
	notes := "Offizielles, einbettbares YouTube-Video des Veranstalters. Wiedergabe ausschließlich über den YouTube-Embed-Player, ohne Abo-Pflicht (YouTube API Developer Policies)."
	if mj.Note != "" {
		notes += " " + mj.Note + "."
	}
	playedOn, err := time.Parse("2006-01-02", mj.PlayedOn)
	if err != nil {
		return err
	}
	if _, err := st.CreateRightsInfo(ctx, media.ID, store.RightsInfoInput{
		Source:        models.RightsOfficialEmbed,
		RightsHolder:  &holder,
		LicenseDocURL: &channel,
		ValidFrom:     playedOn,
		Notes:         &notes,
	}); err != nil {
		return err
	}

	series, err := st.CreateSeries(ctx, media.ID, mj.Title, "Playoffs des "+mj.Title+".")
	if err != nil {
		return err
	}
	season, err := st.CreateSeason(ctx, series.ID, 1, "Playoffs")
	if err != nil {
		return err
	}

	urls := make([]string, len(mj.VideoIDs))
	for i, id := range mj.VideoIDs {
		urls[i] = "https://www.youtube.com/watch?v=" + id
	}

	teamA, teamB := teams[mj.TeamA], teams[mj.TeamB]
	episode, err := st.CreateEpisode(ctx, season.ID, "Grand Final",
		"Grand Final: "+teamA.Name+" gegen "+teamB.Name+" (Best-of-"+strconv.Itoa(mj.BestOf)+").",
		1, mj.DurationMin, urls[0], poster, true)
	if err != nil {
		return err
	}

	scoreA, scoreB := mj.ScoreA, mj.ScoreB
	gameID := game.ID
	return st.UpsertMatch(ctx, episode.ID, store.MatchInput{
		GameID:         &gameID,
		TeamAID:        teamA.ID,
		TeamBID:        teamB.ID,
		Stage:          "Grand Final",
		BestOf:         mj.BestOf,
		ScoreA:         &scoreA,
		ScoreB:         &scoreB,
		PlayedOn:       &playedOn,
		ExtraVideoURLs: urls[1:],
	})
}
