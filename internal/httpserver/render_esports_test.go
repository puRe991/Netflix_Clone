package httpserver

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/pure991/streamflix/internal/models"
)

// TestRenderEsportsTemplates renders every eSports template branch (match
// vs. regular episode, spoilers hidden vs. shown, YouTube vs. self-hosted
// player, with and without a profile) and checks the spoiler guarantee:
// with spoilers hidden the score only appears inside the click-to-reveal
// element and no winner is highlighted.
func TestRenderEsportsTemplates(t *testing.T) {
	tmpls, err := loadTemplates()
	if err != nil {
		t.Fatalf("loadTemplates: %v", err)
	}

	var buf bytes.Buffer
	render := func(name string, data any) string {
		t.Helper()
		buf.Reset()
		if err := tmpls[name].ExecuteTemplate(&buf, "layout", data); err != nil {
			t.Fatalf("%s render failed: %v", name, err)
		}
		return buf.String()
	}

	two, one := 2, 1
	played := time.Date(2017, 7, 23, 0, 0, 0, 0, time.UTC)
	gameID := "g1"
	match := &models.Match{
		EpisodeID: "e1", GameID: &gameID, TeamAID: "ta", TeamBID: "tb", Stage: "Grand Final", BestOf: 3,
		ScoreA: &two, ScoreB: &one, PlayedOn: &played,
		ExtraVideoURLs: []string{"https://youtu.be/sIQ1Eh11Quk"},
		TeamA:          models.Team{ID: "ta", Name: "Gambit Esports", Slug: "gambit"},
		TeamB:          models.Team{ID: "tb", Name: "Immortals", Slug: "immortals"},
	}
	media := models.Media{
		ID: "m1", Title: "PGL Major Kraków 2017", Slug: "pgl-major-krakow-2017", Type: models.MediaSeries,
		Series: &models.Series{Seasons: []models.Season{{SeasonNumber: 1, Title: "Playoffs", Episodes: []models.Episode{
			{ID: "e1", Title: "Grand Final", Duration: 226, IsPublished: true, Match: match},
			{ID: "e2", Title: "Pre-Show", Description: "Analyse vor dem Finale", IsPublished: true},
			{ID: "e3", Title: "Unveröffentlicht", IsPublished: false},
		}}}},
	}

	hidden := render("series.html", seriesData{Media: media, HideSpoilers: true, HasMatches: true, Path: "/series/x"})
	if strings.Contains(hidden, "match-winner") {
		t.Error("winner highlighted although spoilers are hidden")
	}
	if !strings.Contains(hidden, `<details class="spoiler"><summary>Ergebnis anzeigen</summary><span class="match-score">2 : 1</span></details>`) {
		t.Error("score must only appear inside the reveal element while spoilers are hidden")
	}
	if strings.Contains(hidden, "226 Min.") {
		t.Error("match duration leaks the number of maps played")
	}
	if strings.Contains(hidden, "Unveröffentlicht") {
		t.Error("unpublished episode listed on the public series page")
	}

	shown := render("series.html", seriesData{Media: media, HideSpoilers: false, HasMatches: true, Path: "/series/x"})
	if !strings.Contains(shown, `class="match-winner">Gambit Esports`) {
		t.Error("winner not highlighted with spoilers shown")
	}

	// Regular series without matches keeps the classic layout.
	regular := models.Media{ID: "m2", Title: "Red Orbit", Slug: "red-orbit", Type: models.MediaSeries,
		Series: &models.Series{Seasons: []models.Season{{SeasonNumber: 1, Title: "Mission Start", Episodes: []models.Episode{
			{ID: "r1", Title: "Episode 1", Duration: 10, IsPublished: true}}}}}}
	if out := render("series.html", seriesData{Media: regular, HideSpoilers: true}); !strings.Contains(out, "10 Min.") {
		t.Error("regular series lost its duration column")
	}

	tm := []models.TeamMatch{{Match: *match, EpisodeTitle: "Grand Final", MediaTitle: media.Title, MediaSlug: media.Slug, ThumbnailURL: "/static/img/esports/x.svg"}}
	if out := render("team.html", teamData{Team: match.TeamA, Matches: tm, HideSpoilers: true, CanFollow: true, Path: "/team/gambit"}); strings.Contains(out, "match-winner") {
		t.Error("team page highlights winner while spoilers hidden")
	}
	render("team.html", teamData{Team: match.TeamA, HideSpoilers: false})

	short := "GMB"
	teams := []models.Team{match.TeamA, {ID: "tb", Name: "Immortals", Slug: "immortals", ShortName: &short}}
	render("teams.html", teamsData{Teams: teams, Followed: map[string]bool{"ta": true}, CanFollow: true})
	render("teams.html", teamsData{Teams: teams, Followed: map[string]bool{}})
	render("admin_teams.html", adminTeamsData{Teams: teams})

	browseHidden := render("browse.html", browseData{Esports: []models.Media{media}, FollowedMatches: tm, HideSpoilers: true})
	if strings.Contains(browseHidden, "2:1") {
		t.Error("browse row shows score while spoilers hidden")
	}
	if out := render("browse.html", browseData{FollowedMatches: tm, HideSpoilers: false}); !strings.Contains(out, "2:1") {
		t.Error("browse row hides score with spoilers shown")
	}

	yt := render("watch.html", watchData{Title: "x", MediaID: "m1", EpisodeID: "e1", YouTubeIDs: "ttbe_A8Ee50,sIQ1Eh11Quk", HideSpoilers: true})
	if !strings.Contains(yt, `data-video-ids="ttbe_A8Ee50,sIQ1Eh11Quk"`) || strings.Contains(yt, "data-profile-id") {
		t.Error("YouTube watch page: wrong data attributes for anonymous viewer")
	}
	if strings.Contains(yt, "<iframe") {
		t.Error("YouTube must not be loaded before the viewer clicks play")
	}
	if out := render("watch.html", watchData{Title: "x", VideoURL: "https://example.com/a.mp4", MediaID: "m1"}); !strings.Contains(out, "<video") {
		t.Error("self-hosted video lost its <video> player")
	}

	render("admin_episode_form.html", adminEpisodeFormData{Episode: models.Episode{ID: "e2", Title: "Pre-Show"}})
	form := render("admin_episode_form.html", adminEpisodeFormData{
		Episode: models.Episode{ID: "e1", Match: match},
		Teams:   teams,
		Games:   []models.Game{{ID: "g1", Name: "Counter-Strike: Global Offensive"}},
	})
	for _, want := range []string{`value="2017-07-23"`, `<option value="g1" selected>`, "https://youtu.be/sIQ1Eh11Quk"} {
		if !strings.Contains(form, want) {
			t.Errorf("episode form missing %q", want)
		}
	}

	render("profiles.html", profilesData{Profiles: []models.Profile{{ID: "p1", Name: "Admin", HideSpoilers: true}, {ID: "p2", Name: "Kid", IsKidsProfile: true}}})
	render("legal.html", legalData{Title: "Datenschutz"})
}
