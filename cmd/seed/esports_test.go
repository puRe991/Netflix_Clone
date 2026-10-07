package main

import (
	"os"
	"regexp"
	"testing"
)

// TestEsportsCatalogConsistent guards the generated VOD catalog against the
// mistakes that are easy to make when (re)building it: duplicate slugs or
// videos, dangling team/game references, malformed video IDs, results that
// don't fit the format, duplicate episode titles (the seed matches existing
// episodes by title) and missing poster files.
func TestEsportsCatalogConsistent(t *testing.T) {
	cat, err := loadCatalog()
	if err != nil {
		t.Fatalf("loadCatalog: %v", err)
	}
	if len(cat.Tournaments) == 0 {
		t.Fatal("catalog is empty")
	}

	teams := map[string]bool{}
	for _, tm := range cat.Teams {
		if teams[tm.Slug] {
			t.Errorf("duplicate team slug %q", tm.Slug)
		}
		teams[tm.Slug] = true
	}
	games := map[string]bool{}
	for _, g := range cat.Games {
		games[g.Slug] = true
	}

	idPattern := regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)
	allowedStages := map[string]bool{
		"Qualifikation": true, "Online-Cups": true, "Opening Stage": true, "Elimination Stage": true, "Challengers Stage": true, "Gruppenphase": true,
		"Legends Stage": true, "Playoffs": true, "Relegation": true, "Weitere Matches": true,
	}
	// Every source must be an organizer's official channel (see docs/esports-vods.md).
	officialChannels := map[string]bool{}
	for _, h := range []string{"@esl", "@ESLCS", "@ESLArchives", "@ESLDeutschland", "@pgl", "@PGL_VODs", "@BLASTPremier",
		"@starladder_cs", "@starladder_cs_highlights", "@FACEIT", "@MajorLeagueGaming", "@WorldCyberGamesOfficial", "@ESEA", "@99DMGCSGO"} {
		officialChannels["https://www.youtube.com/"+h] = true
	}
	slugs := map[string]bool{}
	videos := map[string]string{}
	finals := 0

	for _, tr := range cat.Tournaments {
		if slugs[tr.Slug] {
			t.Errorf("duplicate tournament slug %q", tr.Slug)
		}
		slugs[tr.Slug] = true
		if !games[tr.Game] {
			t.Errorf("%s: unknown game %q", tr.Slug, tr.Game)
		}
		switch tr.Category {
		case "major", "international", "dach-liga":
		default:
			t.Errorf("%s: unknown category %q", tr.Slug, tr.Category)
		}
		if len(tr.Channels) == 0 {
			t.Errorf("%s: no source channel", tr.Slug)
		}
		for _, ch := range tr.Channels {
			if !officialChannels[ch] {
				t.Errorf("%s: channel %s is not on the official list", tr.Slug, ch)
			}
		}
		for _, kind := range []string{"poster", "banner"} {
			if _, err := os.Stat("../../web/static/img/esports/" + tr.Slug + "-" + kind + ".svg"); err != nil {
				t.Errorf("%s: missing %s image", tr.Slug, kind)
			}
		}

		seenStage := map[string]bool{}
		for _, st := range tr.Stages {
			if !allowedStages[st.Title] || seenStage[st.Title] {
				t.Errorf("%s: unexpected or duplicate stage %q", tr.Slug, st.Title)
			}
			seenStage[st.Title] = true

			titles := map[string]bool{}
			for _, ep := range st.Episodes {
				where := tr.Slug + "/" + st.Title + "/" + ep.Title
				if ep.Title == "" || titles[ep.Title] {
					t.Errorf("%s: empty or duplicate episode title", where)
				}
				titles[ep.Title] = true
				if ep.Title == "Grand Final" {
					finals++
				}
				if len(ep.Videos) == 0 {
					t.Errorf("%s: no videos", where)
				}
				for _, id := range ep.Videos {
					if !idPattern.MatchString(id) {
						t.Errorf("%s: malformed YouTube ID %q", where, id)
					}
					if other, dup := videos[id]; dup {
						t.Errorf("video %s used twice (%s and %s)", id, other, where)
					}
					videos[id] = where
				}

				if (ep.TeamA == "") != (ep.TeamB == "") || (ep.TeamA != "" && (!teams[ep.TeamA] || !teams[ep.TeamB] || ep.TeamA == ep.TeamB)) {
					t.Errorf("%s: bad teams %q / %q", where, ep.TeamA, ep.TeamB)
				}
				switch ep.BestOf {
				case 0, 1, 3, 5:
				default:
					t.Errorf("%s: invalid best-of %d", where, ep.BestOf)
				}
				if ep.Score != nil {
					if len(ep.Score) != 2 || ep.BestOf == 0 {
						t.Fatalf("%s: a score needs two values and a known format", where)
					}
					a, b := ep.Score[0], ep.Score[1]
					need := ep.BestOf/2 + 1
					if (a != need) == (b != need) || a+b > ep.BestOf {
						t.Errorf("%s: %d:%d is not a finished Bo%d", where, a, b, ep.BestOf)
					}
					if len(ep.Videos) > 1 && a+b != len(ep.Videos) {
						t.Errorf("%s: %d map uploads but score %d:%d", where, len(ep.Videos), a, b)
					}
				}
			}
		}
	}
	if finals < 46 {
		t.Errorf("expected at least the 46 verified grand finals, got %d", finals)
	}
}
