package httpserver

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

// Generated artwork for eSports content. Real team logos are trademarks and
// YouTube thumbnails frequently show the winner celebrating, so the
// platform draws its own spoiler-free graphics: a monogram badge per team
// and a matchup card per episode, in the tournament's colour (the same hue
// formula as the generated tournament posters).

// hueFor maps a slug to a stable hue (0-359): first 4 hex digits of its MD5.
func hueFor(slug string) int {
	sum := md5.Sum([]byte(slug))
	n, _ := strconv.ParseInt(hex.EncodeToString(sum[:])[:4], 16, 64)
	return int(n % 360)
}

// monogram returns the badge text: the short name if set, otherwise up to
// three initials (or the first letters of a one-word name).
func monogram(name string, short *string) string {
	if short != nil && strings.TrimSpace(*short) != "" {
		return truncateRunes(strings.TrimSpace(*short), 4)
	}
	var words []string
	for _, w := range strings.FieldsFunc(name, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		switch strings.ToLower(w) {
		case "team", "gaming", "esports", "esport", "clan":
			continue
		}
		words = append(words, w)
	}
	if len(words) == 0 {
		return "?"
	}
	if len(words) == 1 {
		return strings.ToUpper(truncateRunes(words[0], 3))
	}
	var b strings.Builder
	for _, w := range words[:min(3, len(words))] {
		b.WriteRune(unicode.ToUpper([]rune(w)[0]))
	}
	return b.String()
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func ellipsize(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

const shieldPath = "M100 8 L182 36 V96 C182 146 146 178 100 194 C54 178 18 146 18 96 V36 Z"

// badgeGroup draws a team shield centred on (cx, cy) at the given scale.
func badgeGroup(cx, cy, scale float64, slug, text string) string {
	hue := hueFor(slug)
	id := "b" + strconv.Itoa(hue) + strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, slug)
	fontSize := 78
	switch n := len([]rune(text)); {
	case n >= 4:
		fontSize = 50
	case n == 3:
		fontSize = 62
	}
	return fmt.Sprintf(`<g transform="translate(%.1f %.1f) scale(%.3f) translate(-100 -100)">`+
		`<defs><linearGradient id="%s" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="hsl(%d,62%%,48%%)"/><stop offset="1" stop-color="hsl(%d,70%%,22%%)"/></linearGradient></defs>`+
		`<path d="%s" fill="url(#%s)" stroke="hsl(%d,80%%,72%%)" stroke-width="6"/>`+
		`<path d="M100 22 L168 46 V96 C168 138 138 166 100 180" fill="none" stroke="#fff" stroke-opacity=".18" stroke-width="4"/>`+
		`<text x="100" y="112" text-anchor="middle" dominant-baseline="middle" font-family="Arial,Helvetica,sans-serif" font-weight="700" font-size="%d" fill="#fff">%s</text>`+
		`</g>`,
		cx, cy, scale, id, hue, (hue+20)%360, shieldPath, id, hue, fontSize, html.EscapeString(text))
}

// teamBadgeSVG is the standalone 200x200 badge served for a team.
func teamBadgeSVG(t models.Team) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" viewBox="0 0 200 200">` +
		badgeGroup(100, 100, 1, t.Slug, monogram(t.Name, t.ShortName)) + `</svg>`
}

// episodeArtSVG is the 16:9 card for an episode: a matchup for matches, a
// titled broadcast card otherwise. It never contains a result.
func episodeArtSVG(a store.EpisodeArtwork) string {
	hue := hueFor(a.MediaSlug)
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="640" height="360" viewBox="0 0 640 360">`+
		`<defs><linearGradient id="bg" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="hsl(%d,65%%,30%%)"/><stop offset="1" stop-color="hsl(%d,70%%,9%%)"/></linearGradient>`+
		`<pattern id="st" width="28" height="28" patternUnits="userSpaceOnUse" patternTransform="rotate(30)"><rect width="2" height="28" fill="#fff" opacity=".05"/></pattern></defs>`+
		`<rect width="640" height="360" fill="url(#bg)"/><rect width="640" height="360" fill="url(#st)"/>`, hue, (hue+40)%360)

	header := ellipsize(a.MediaTitle, 46)
	if year := strconv.Itoa(a.ReleaseYear); !strings.Contains(a.MediaTitle, year) {
		header += " · " + year
	}
	fmt.Fprintf(&b, `<text x="24" y="38" font-family="Arial,Helvetica,sans-serif" font-weight="700" font-size="17" letter-spacing="1.5" fill="hsl(%d,85%%,72%%)">%s</text>`,
		hue, html.EscapeString(strings.ToUpper(header)))

	ep := a.Episode
	if m := ep.Match; m != nil {
		b.WriteString(badgeGroup(160, 160, 0.72, m.TeamA.Slug, monogram(m.TeamA.Name, m.TeamA.ShortName)))
		b.WriteString(badgeGroup(480, 160, 0.72, m.TeamB.Slug, monogram(m.TeamB.Name, m.TeamB.ShortName)))
		b.WriteString(`<circle cx="320" cy="160" r="34" fill="#000" fill-opacity=".35" stroke="#fff" stroke-opacity=".35" stroke-width="2"/>` +
			`<text x="320" y="161" text-anchor="middle" dominant-baseline="middle" font-family="Arial,Helvetica,sans-serif" font-weight="700" font-size="26" fill="#fff">VS</text>`)
		for _, t := range []struct {
			x    int
			name string
		}{{160, m.TeamA.Name}, {480, m.TeamB.Name}} {
			fmt.Fprintf(&b, `<text x="%d" y="262" text-anchor="middle" font-family="Arial,Helvetica,sans-serif" font-weight="700" font-size="22" fill="#fff">%s</text>`,
				t.x, html.EscapeString(ellipsize(t.name, 20)))
		}
		stage := m.Stage
		if stage == "" {
			stage = a.SeasonTitle
		}
		if stage != a.SeasonTitle && !strings.HasPrefix(stage, a.SeasonTitle) {
			stage = a.SeasonTitle + " · " + stage
		}
		fmt.Fprintf(&b, `<text x="320" y="326" text-anchor="middle" font-family="Arial,Helvetica,sans-serif" font-size="17" fill="#fff" fill-opacity=".8">%s</text>`,
			html.EscapeString(ellipsize(stage, 60)))
	} else {
		b.WriteString(`<circle cx="320" cy="140" r="44" fill="#000" fill-opacity=".35" stroke="#fff" stroke-opacity=".4" stroke-width="2"/>` +
			`<path d="M306 116 L344 140 L306 164 Z" fill="#fff"/>`)
		lines := wrapWords(ep.Title, 34, 2)
		for i, l := range lines {
			fmt.Fprintf(&b, `<text x="320" y="%d" text-anchor="middle" font-family="Arial,Helvetica,sans-serif" font-weight="700" font-size="22" fill="#fff">%s</text>`,
				236+i*30, html.EscapeString(l))
		}
		fmt.Fprintf(&b, `<text x="320" y="326" text-anchor="middle" font-family="Arial,Helvetica,sans-serif" font-size="17" fill="#fff" fill-opacity=".8">%s</text>`,
			html.EscapeString(ellipsize(a.SeasonTitle, 60)))
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// wrapWords greedily wraps text into at most maxLines lines of width runes,
// ellipsizing the last line if needed.
func wrapWords(text string, width, maxLines int) []string {
	var lines []string
	cur := ""
	for _, w := range strings.Fields(text) {
		if cur == "" {
			cur = w
		} else if len([]rune(cur))+1+len([]rune(w)) <= width {
			cur += " " + w
		} else {
			lines = append(lines, cur)
			cur = w
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines[maxLines-1] = ellipsize(lines[maxLines-1]+" …", width)
	}
	for i := range lines {
		lines[i] = ellipsize(lines[i], width)
	}
	return lines
}

func writeSVG(w http.ResponseWriter, svg string) {
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	// Defence in depth: the SVG is only ever used as an image, never as a document.
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write([]byte(svg))
}

// TeamBadge serves /img/team/{slug}.svg.
func (s *Server) TeamBadge(w http.ResponseWriter, r *http.Request) error {
	team, err := s.Store.GetTeamBySlug(r.Context(), strings.TrimSuffix(r.PathValue("file"), ".svg"))
	if err != nil {
		return err
	}
	writeSVG(w, teamBadgeSVG(team))
	return nil
}

// EpisodeArt serves /img/episode/{id}.svg. Unpublished episodes get no art
// for non-admins, matching who may see them.
var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func (s *Server) EpisodeArt(w http.ResponseWriter, r *http.Request) error {
	id := strings.TrimSuffix(r.PathValue("file"), ".svg")
	if !uuidPattern.MatchString(id) {
		return store.ErrNotFound
	}
	a, err := s.Store.GetEpisodeArtwork(r.Context(), id)
	if err != nil {
		return err
	}
	if !a.Episode.IsPublished {
		if u := s.Sessions.OptionalUser(r); u == nil || u.Role != models.RoleAdmin {
			return store.ErrNotFound
		}
	}
	writeSVG(w, episodeArtSVG(a))
	return nil
}

// teamLogo is the template helper for a team's image: an admin-provided
// logo URL wins, otherwise the generated badge.
func teamLogo(t models.Team) string {
	if t.LogoURL != nil && *t.LogoURL != "" {
		return *t.LogoURL
	}
	return "/img/team/" + t.Slug + ".svg"
}
