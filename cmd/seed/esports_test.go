package main

import (
	"os"
	"regexp"
	"testing"
)

// TestSeedFinalsConsistent guards the hand-maintained VOD list against the
// mistakes that are easy to make when adding an event: duplicate slugs,
// typos in team/game slugs, malformed video IDs, scores that don't fit the
// format and missing poster files.
func TestSeedFinalsConsistent(t *testing.T) {
	teams := map[string]bool{}
	for _, tm := range esportsTeams {
		if teams[tm.Slug] {
			t.Errorf("duplicate team slug %q", tm.Slug)
		}
		teams[tm.Slug] = true
	}
	games := map[string]bool{}
	for _, g := range esportsGames {
		games[g.Slug] = true
	}

	idPattern := regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)
	slugs := map[string]bool{}
	videos := map[string]string{}
	for _, f := range csFinals {
		if slugs[f.Slug] {
			t.Errorf("duplicate final slug %q", f.Slug)
		}
		slugs[f.Slug] = true

		if !teams[f.TeamA] || !teams[f.TeamB] || f.TeamA == f.TeamB {
			t.Errorf("%s: bad teams %q / %q", f.Slug, f.TeamA, f.TeamB)
		}
		if !games[f.Game] {
			t.Errorf("%s: unknown game %q", f.Slug, f.Game)
		}
		switch f.Category {
		case catMajor, catClassic, catDACH:
		default:
			t.Errorf("%s: unknown category %q", f.Slug, f.Category)
		}
		if f.BestOf != 1 && f.BestOf != 3 && f.BestOf != 5 {
			t.Errorf("%s: invalid best-of %d", f.Slug, f.BestOf)
		}
		if len(f.VideoIDs) == 0 || f.DurationMin <= 0 {
			t.Errorf("%s: needs videos and a duration", f.Slug)
		}
		for _, id := range f.VideoIDs {
			if !idPattern.MatchString(id) {
				t.Errorf("%s: malformed YouTube ID %q", f.Slug, id)
			}
			if other, dup := videos[id]; dup {
				t.Errorf("video %s used by both %s and %s", id, other, f.Slug)
			}
			videos[id] = f.Slug
		}

		if f.Score != nil {
			if len(f.Score) != 2 {
				t.Fatalf("%s: score must have two values", f.Slug)
			}
			a, b := f.Score[0], f.Score[1]
			need := f.BestOf/2 + 1
			if (a != need) == (b != need) || a+b > f.BestOf {
				t.Errorf("%s: %d:%d is not a finished Bo%d", f.Slug, a, b, f.BestOf)
			}
			// One upload per map: the number of parts must match the maps played.
			if len(f.VideoIDs) > 1 && a+b != len(f.VideoIDs) {
				t.Errorf("%s: %d map uploads but score %d:%d", f.Slug, len(f.VideoIDs), a, b)
			}
		}

		for _, kind := range []string{"poster", "banner"} {
			path := "../../web/static/img/esports/" + f.Slug + "-" + kind + ".svg"
			if _, err := os.Stat(path); err != nil {
				t.Errorf("%s: missing %s", f.Slug, path)
			}
		}
	}
}
