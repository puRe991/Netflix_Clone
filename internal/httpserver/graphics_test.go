package httpserver

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

// wellFormed fails the test if svg isn't parseable XML.
func wellFormed(t *testing.T, name, svg string) {
	t.Helper()
	d := xml.NewDecoder(strings.NewReader(svg))
	for {
		if _, err := d.Token(); err == io.EOF {
			return
		} else if err != nil {
			t.Fatalf("%s: invalid SVG: %v\n%s", name, err, svg)
		}
	}
}

func TestMonogram(t *testing.T) {
	short := "NAVI"
	cases := []struct {
		name  string
		short *string
		want  string
	}{
		{"Natus Vincere", &short, "NAVI"},
		{"Ninjas in Pyjamas", nil, "NIP"},
		{"Team Vitality", nil, "VIT"},
		{"FaZe Clan", nil, "FAZ"},
		{"Astralis", nil, "AST"},
		{"G2 Esports", nil, "G2"},
		{"Meet Your Makers", nil, "MYM"},
		{"", nil, "?"},
	}
	for _, c := range cases {
		if got := monogram(c.name, c.short); got != c.want {
			t.Errorf("monogram(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestHueMatchesPosterGenerator(t *testing.T) {
	// The poster script uses int(md5(slug).hexdigest()[:4], 16) % 360.
	if h := hueFor("pgl-major-krakow-2017"); h < 0 || h >= 360 {
		t.Fatalf("hue out of range: %d", h)
	}
	if hueFor("a") != hueFor("a") || hueFor("a") == hueFor("b") {
		t.Error("hue must be stable per slug and differ between slugs")
	}
}

func TestGeneratedArtwork(t *testing.T) {
	evil := "<script>alert(1)</script> & \"Co\""
	team := models.Team{ID: "t1", Name: evil, Slug: "evil-team"}
	badge := teamBadgeSVG(team)
	wellFormed(t, "badge", badge)
	if strings.Contains(badge, "<script>") {
		t.Error("team name not escaped in badge")
	}

	two, one := 2, 1
	match := &models.Match{
		Stage: "Runde 3", BestOf: 3, ScoreA: &two, ScoreB: &one,
		TeamA: models.Team{Name: "Gambit Esports", Slug: "gambit"},
		TeamB: team,
	}
	art := store.EpisodeArtwork{
		Episode:     models.Episode{ID: "e1", Title: "Runde 3: Gambit vs. X", Match: match},
		SeasonTitle: "Legends Stage", MediaTitle: "PGL Major Kraków 2017", MediaSlug: "pgl-major-krakow-2017", ReleaseYear: 2017,
	}
	svg := episodeArtSVG(art)
	wellFormed(t, "match art", svg)
	if strings.Contains(svg, "<script>") {
		t.Error("team name not escaped in match art")
	}
	if strings.Contains(svg, "2:1") || strings.Contains(svg, "2 : 1") {
		t.Error("match art must never show the result")
	}
	if !strings.Contains(svg, "Legends Stage · Runde 3") || !strings.Contains(svg, ">PGL MAJOR KRAKÓW 2017<") {
		t.Error("match art misses stage or tournament header")
	}

	broadcast := art
	broadcast.Episode = models.Episode{ID: "e2", Title: "Übertragung: Challengers Stage – Tag 1, Stream A mit einem sehr langen Titel der umbrechen muss"}
	wellFormed(t, "broadcast art", episodeArtSVG(broadcast))
}

func TestTeamLogo(t *testing.T) {
	custom := "https://example.com/logo.png"
	if got := teamLogo(models.Team{Slug: "fnatic"}); got != "/img/team/fnatic.svg" {
		t.Errorf("default logo = %q", got)
	}
	if got := teamLogo(models.Team{Slug: "fnatic", LogoURL: &custom}); got != custom {
		t.Errorf("custom logo = %q", got)
	}
}
