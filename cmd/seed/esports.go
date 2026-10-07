package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/pure991/streamflix/internal/models"
	"github.com/pure991/streamflix/internal/store"
)

// Counter-Strike grand finals from 2006 to today, each as a tournament
// (SERIES) with a "Playoffs" season and the grand final as its match
// episode.
//
// Every video is an upload on the organizer's own official YouTube channel
// (ESL, ESL Deutschland, ESL Archives, PGL, BLAST, StarLadder, FACEIT, MLG,
// WCG, ESEA, 99Damage), checked via YouTube oEmbed to exist and allow
// embedding. See docs/esports-vods.md for the full list, the channels and
// the events without an official VOD. Fan re-uploads are deliberately not
// used.
//
// Team order follows the official video titles (bracket order), never
// "winner first", so the listing itself doesn't give results away.

type seedTeam struct {
	Slug, Name, Short string
}

var esportsTeams = []seedTeam{
	{"again", "AGAiN", "AGAiN"},
	{"alternate", "ALTERNATE", "ALT"},
	{"alternate-attax", "ALTERNATE aTTaX", "ATN"},
	{"astralis", "Astralis", "AST"},
	{"avangar", "AVANGAR", "AVG"},
	{"berzerk", "Berzerk", "BRZ"},
	{"big", "BIG", "BIG"},
	{"dotpixels", "dotpiXels", "dpX"},
	{"dynamic", "Dynamic", "DYN"},
	{"emulate", "emuLate", "eMg"},
	{"ence", "ENCE", "ENCE"},
	{"envyus", "Team EnVyUs", "nV"},
	{"esc-gaming", "ESC Gaming", "ESC"},
	{"euronics-gaming", "EURONICS Gaming", "EUR"},
	{"expert-esport", "expert eSport", "EXP"},
	{"falcons", "Team Falcons", "FLC"},
	{"faze", "FaZe Clan", "FaZe"},
	{"fnatic", "Fnatic", "FNC"},
	{"frag-executors", "Frag eXecutors", "FX"},
	{"furia", "FURIA", "FUR"},
	{"g2", "G2 Esports", "G2"},
	{"gambit", "Gambit Esports", "GMB"},
	{"gamerlegion", "GamerLegion", "GL"},
	{"heroic", "Heroic", "HER"},
	{"immortals", "Immortals", "IMT"},
	{"leisure", "LeiSuRe", "LSR"},
	{"liquid", "Team Liquid", "TL"},
	{"luminosity", "Luminosity Gaming", "LG"},
	{"meet-your-makers", "Meet Your Makers", "MYM"},
	{"mortal-teamwork", "mTw", "mTw"},
	{"mousesports", "mousesports", "mouz"},
	{"n-faculty", "n!faculty", "n!f"},
	{"natus-vincere", "Natus Vincere", "NAVI"},
	{"nip", "Ninjas in Pyjamas", "NiP"},
	{"noa", "Team NoA", "NoA"},
	{"outsiders", "Outsiders", "OUT"},
	{"panthers-gaming", "PANTHERS Gaming", "PAN"},
	{"pentagram", "Pentagram", "PGS"},
	{"planetkey-dynamics", "Planetkey Dynamics", "PkD"},
	{"sk-gaming", "SK Gaming", "SK"},
	{"spirit", "Team Spirit", "SPR"},
	{"sprout", "Sprout", "SPR"},
	{"team-wild-fire", "Team Wild Fire", "TWF"},
	{"the-mongolz", "The MongolZ", "MGLZ"},
	{"verygames", "VeryGames", "VG"},
	{"virtus-pro", "Virtus.pro", "VP"},
	{"vitality", "Team Vitality", "VIT"},
	{"zomblerz", "Zomblerz", "ZMB"},
}

var esportsGames = []struct{ Slug, Name, Label string }{
	{"cs16", "Counter-Strike 1.6", "CS 1.6"},
	{"css", "Counter-Strike: Source", "CS:Source"},
	{"csgo", "Counter-Strike: Global Offensive", "CS:GO"},
	{"cs2", "Counter-Strike 2", "CS2"},
}

// Browse-row categories, stored as media tags.
const (
	catMajor   = "major"
	catClassic = "international"
	catDACH    = "dach-liga"
)

type seedFinal struct {
	Slug        string
	Title       string
	Year        int
	Location    string
	Game        string // cs16, css, csgo, cs2
	Category    string // catMajor, catClassic, catDACH
	Organizer   string
	ChannelURL  string
	PlayedOn    string // YYYY-MM-DD, "" when not reliably known
	TeamA       string
	TeamB       string
	BestOf      int
	Score       []int // {a, b}; nil when the result isn't verified
	DurationMin int
	VideoIDs    []string
	Language    string
	Note        string
}

func score(a, b int) []int { return []int{a, b} }

const (
	chESL           = "https://www.youtube.com/@esl"
	chESLCS         = "https://www.youtube.com/@ESLCS"
	chESLArchives   = "https://www.youtube.com/@ESLArchives"
	chESLDeutsch    = "https://www.youtube.com/@ESLDeutschland"
	chPGL           = "https://www.youtube.com/@pgl"
	chBLAST         = "https://www.youtube.com/@BLASTPremier"
	chWCG           = "https://www.youtube.com/@WorldCyberGamesOfficial"
	chESEA          = "https://www.youtube.com/@ESEA"
	ch99Damage      = "https://www.youtube.com/@99DMGCSGO"
	langEnglish     = "Englisch"
	langGerman      = "Deutsch"
	noteFullShow    = "Komplette Übertragung des Finales"
	noteResultOpen  = "Ergebnis nicht verifiziert und daher nicht hinterlegt"
	noteFullClassic = "ESL Classics: komplette Übertragung des Finales"
)

var csFinals = []seedFinal{
	// --- Valve Majors (CS:GO / CS2) ---
	{"ems-one-katowice-2014", "EMS One Katowice 2014", 2014, "Katowice", "csgo", catMajor, "ESL", chESL,
		"2014-03-16", "virtus-pro", "nip", 3, score(2, 0), 150, []string{"do4FUYUHmUM"}, langEnglish, ""},
	{"esl-one-cologne-2014", "ESL One Cologne 2014", 2014, "Köln", "csgo", catMajor, "ESL", chESLCS,
		"2014-08-17", "nip", "fnatic", 3, score(2, 1), 215, []string{"-5laY7bjono"}, langEnglish, noteFullClassic},
	{"esl-one-katowice-2015", "ESL One Katowice 2015", 2015, "Katowice", "csgo", catMajor, "ESL", chESLCS,
		"2015-03-15", "fnatic", "nip", 3, score(2, 1), 200, []string{"cqmdRLNRZHA"}, langEnglish, noteFullClassic},
	{"esl-one-cologne-2015", "ESL One Cologne 2015", 2015, "Köln", "csgo", catMajor, "ESL", chESLArchives,
		"2015-08-23", "fnatic", "envyus", 3, score(2, 0), 142, []string{"02I5vVxlJhU", "_Ef6cXH1oC8"}, langEnglish, ""},
	{"mlg-major-columbus-2016", "MLG Major Championship: Columbus 2016", 2016, "Columbus", "csgo", catMajor, "Major League Gaming", "https://www.youtube.com/@MajorLeagueGaming",
		"2016-04-03", "natus-vincere", "luminosity", 3, score(0, 2), 203, []string{"7nWxZin959o"}, langEnglish, ""},
	{"esl-one-cologne-2016", "ESL One Cologne 2016", 2016, "Köln", "csgo", catMajor, "ESL", chESLArchives,
		"2016-07-10", "sk-gaming", "liquid", 3, score(2, 0), 145, []string{"AK5RY8W-AcM", "TC6N1zHvXpU"}, langEnglish, ""},
	{"pgl-major-krakow-2017", "PGL Major Kraków 2017", 2017, "Kraków", "csgo", catMajor, "PGL", chPGL,
		"2017-07-23", "gambit", "immortals", 3, score(2, 1), 226, []string{"ttbe_A8Ee50", "sIQ1Eh11Quk", "w2MdCEZSBtw"}, langEnglish, ""},
	{"faceit-major-london-2018", "FACEIT Major: London 2018", 2018, "London", "csgo", catMajor, "FACEIT", "https://www.youtube.com/@FACEIT",
		"2018-09-23", "natus-vincere", "astralis", 3, score(0, 2), 171, []string{"z7IDPdmYw0A", "9CcrjvC4SkE"}, langEnglish, ""},
	{"iem-katowice-major-2019", "IEM Katowice Major 2019", 2019, "Katowice", "csgo", catMajor, "ESL", chESLArchives,
		"2019-03-03", "ence", "astralis", 3, score(0, 2), 151, []string{"kgitmggEgrA", "A_wk85oA0Rc"}, langEnglish, ""},
	{"starladder-major-berlin-2019", "StarLadder Major: Berlin 2019", 2019, "Berlin", "csgo", catMajor, "StarLadder", "https://www.youtube.com/@starladder_cs_highlights",
		"2019-09-08", "avangar", "astralis", 3, score(0, 2), 188, []string{"c2uQV59-tbI", "hkVQXKCVSFQ"}, langEnglish, ""},
	{"pgl-major-stockholm-2021", "PGL Major Stockholm 2021", 2021, "Stockholm", "csgo", catMajor, "PGL", chPGL,
		"2021-11-07", "g2", "natus-vincere", 3, score(0, 2), 146, []string{"UKYnAujcFac", "pogK6a25LO0"}, langEnglish, ""},
	{"pgl-major-antwerp-2022", "PGL Major Antwerp 2022", 2022, "Antwerpen", "csgo", catMajor, "PGL", chPGL,
		"2022-05-22", "faze", "natus-vincere", 3, score(2, 0), 136, []string{"O9u1FQKZHVw", "CMBAVBkC8Zw"}, langEnglish, ""},
	{"iem-rio-major-2022", "IEM Rio Major 2022", 2022, "Rio de Janeiro", "csgo", catMajor, "ESL", chESLArchives,
		"2022-11-13", "outsiders", "heroic", 3, score(2, 0), 113, []string{"9pscaxA1QBs", "oVH782n-32U"}, langEnglish, ""},
	{"blast-paris-major-2023", "BLAST.tv Paris Major 2023", 2023, "Paris", "csgo", catMajor, "BLAST", chBLAST,
		"2023-05-21", "vitality", "gamerlegion", 3, score(2, 0), 410, []string{"y6Y0l1Y-PTo"}, langEnglish, "Komplette Übertragung des Finaltags inkl. Vorberichterstattung"},
	{"pgl-major-copenhagen-2024", "PGL CS2 Major Copenhagen 2024", 2024, "Kopenhagen", "cs2", catMajor, "PGL", chPGL,
		"2024-03-31", "faze", "natus-vincere", 3, score(1, 2), 112, []string{"09N-bQbLrx8", "L9eBTSwDwaM", "vdXveK6G70g"}, langEnglish, ""},
	{"perfect-world-shanghai-major-2024", "Perfect World Shanghai Major 2024", 2024, "Shanghai", "cs2", catMajor, "Perfect World / PGL", chPGL,
		"2024-12-15", "spirit", "faze", 3, score(2, 1), 413, []string{"Vv7eXxl--LA"}, langEnglish, "Komplette Übertragung inkl. Showmatch vor dem Finale"},
	{"blast-austin-major-2025", "BLAST.tv Austin Major 2025", 2025, "Austin", "cs2", catMajor, "BLAST", chBLAST,
		"2025-06-22", "vitality", "the-mongolz", 3, score(2, 1), 382, []string{"meijWUaxNYg"}, langEnglish, "Komplette Übertragung des Finaltags"},
	{"starladder-budapest-major-2025", "StarLadder Budapest Major 2025", 2025, "Budapest", "cs2", catMajor, "StarLadder", "https://www.youtube.com/@starladder_cs",
		"2025-12-14", "vitality", "faze", 5, score(3, 1), 169, []string{"9qqiaW7sjD0"}, langEnglish, "Erstes Major-Finale im Best-of-5-Format"},
	{"iem-cologne-major-2026", "IEM Cologne Major 2026", 2026, "Köln", "cs2", catMajor, "ESL", chESLArchives,
		"2026-06-21", "furia", "falcons", 5, score(0, 3), 164, []string{"p6U28DEBn8A"}, langEnglish, ""},

	// --- International classics: WCG (CS 1.6), Intel Extreme Masters (CS 1.6), ESEA LAN (CS:Source), EMS One ---
	{"wcg-2006-monza", "World Cyber Games 2006", 2006, "Monza", "cs16", catClassic, "World Cyber Games", chWCG,
		"", "nip", "pentagram", 3, score(1, 2), 102, []string{"YSEtpv7n8tw", "n89xjEA6VdE", "h6gILRcMFVA"}, langEnglish, ""},
	{"wcg-2007-seattle", "World Cyber Games 2007", 2007, "Seattle", "cs16", catClassic, "World Cyber Games", chWCG,
		"", "emulate", "noa", 3, score(2, 1), 123, []string{"DEQRNpv6fHM", "gJsTe85XzNQ", "FEQCLTTy5QY"}, langEnglish, ""},
	{"wcg-2008-cologne", "World Cyber Games 2008", 2008, "Köln", "cs16", catClassic, "World Cyber Games", chWCG,
		"", "sk-gaming", "mortal-teamwork", 3, score(1, 2), 127, []string{"J5DbuSDp5t8", "P5-ubodo6A0", "XLHlDEN7HHc"}, langEnglish, ""},
	{"wcg-2009-chengdu", "World Cyber Games 2009", 2009, "Chengdu", "cs16", catClassic, "World Cyber Games", chWCG,
		"", "fnatic", "again", 3, nil, 136, []string{"AaUSjLtG4aM", "kPt1zrKbviY"}, langEnglish, noteResultOpen},
	{"wcg-2010-los-angeles", "World Cyber Games 2010", 2010, "Los Angeles", "cs16", catClassic, "World Cyber Games", chWCG,
		"", "mortal-teamwork", "natus-vincere", 3, score(1, 2), 158, []string{"kVa3-QcDDlc", "cVz_yRkxsUg", "Q_JSeHJq9Cw"}, langEnglish, "Team mTw trat als „Denmark2“ an"},
	{"iem-iv-chengdu-2009", "Intel Extreme Masters IV – Global Challenge Chengdu 2009", 2009, "Chengdu", "cs16", catClassic, "ESL", chESL,
		"", "fnatic", "sk-gaming", 3, nil, 154, []string{"L3dYILqFunE"}, langEnglish, noteFullClassic + "; " + noteResultOpen},
	{"iem-iv-dubai-2009", "Intel Extreme Masters IV – Global Challenge Dubai 2009", 2009, "Dubai", "cs16", catClassic, "ESL", chESL,
		"", "fnatic", "meet-your-makers", 3, nil, 114, []string{"hZeysLd0QSE"}, langEnglish, noteFullClassic + "; " + noteResultOpen},
	{"iem-iv-world-championship-2010", "Intel Extreme Masters IV – World Championship 2010", 2010, "Hannover", "cs16", catClassic, "ESL", chESL,
		"", "fnatic", "natus-vincere", 3, nil, 184, []string{"LzWsMTSdY3A"}, langEnglish, noteFullClassic + "; " + noteResultOpen},
	{"iem-v-cologne-2010", "Intel Extreme Masters V – Global Challenge Cologne 2010", 2010, "Köln", "cs16", catClassic, "ESL", chESL,
		"", "fnatic", "mousesports", 3, nil, 188, []string{"BJbo61X3D58"}, langEnglish, noteFullClassic + "; " + noteResultOpen},
	{"iem-v-world-championship-2011", "Intel Extreme Masters V – World Championship 2011", 2011, "Hannover", "cs16", catClassic, "ESL", chESLArchives,
		"", "natus-vincere", "frag-executors", 3, score(2, 0), 135, []string{"oicR25CGIcc", "LGNsKZsK9BI"}, langEnglish, ""},
	{"iem-vi-guangzhou-2011", "Intel Extreme Masters VI – Global Challenge Guangzhou 2011", 2011, "Guangzhou", "cs16", catClassic, "ESL", chESLArchives,
		"", "fnatic", "mousesports", 3, score(2, 1), 146, []string{"YuALZm8OkCc", "z7RxpTAYDZ4", "Dhpo3TF9JQc"}, langEnglish, ""},
	{"iem-vi-new-york-2011", "Intel Extreme Masters VI – Global Challenge New York 2011", 2011, "New York", "cs16", catClassic, "ESL", chESL,
		"", "sk-gaming", "fnatic", 3, nil, 112, []string{"71Kf35IrPKc", "iw5dVBKnuyw"}, langEnglish, noteResultOpen},
	{"iem-vi-kiev-2012", "Intel Extreme Masters VI – Global Challenge Kiev 2012", 2012, "Kiew", "cs16", catClassic, "ESL", chESLArchives,
		"", "natus-vincere", "sk-gaming", 3, score(2, 0), 128, []string{"mMVk1YubLO0"}, langEnglish, ""},
	{"iem-vi-world-championship-2012", "Intel Extreme Masters VI – World Championship 2012", 2012, "Hannover", "cs16", catClassic, "ESL", chESL,
		"", "esc-gaming", "natus-vincere", 3, score(2, 0), 149, []string{"S82eTX3Wiq0", "fiF42bT206o"}, langEnglish, ""},
	{"esea-lan-season-10", "ESEA LAN Season 10", 2012, "Dallas", "css", catClassic, "ESEA", chESEA,
		"2012-03-03", "dynamic", "zomblerz", 3, score(2, 1), 157, []string{"Y6A7ZgrekBE", "vz0TjGTO-gM", "uTgITPAW__8"}, langEnglish, ""},
	{"raidcall-ems-one-finals-2013", "RaidCall EMS One Finals 2013", 2013, "Katowice", "csgo", catClassic, "ESL", chESL,
		"", "verygames", "virtus-pro", 3, nil, 121, []string{"N-GiZJbwD24"}, langEnglish, noteResultOpen},

	// --- German leagues: ESL Pro Series (EPS), ESL Meisterschaft, 99Damage Liga ---
	{"eps-summer-2012-cs16", "ESL Pro Series Summer Finals 2012 (CS 1.6)", 2012, "Deutschland", "cs16", catDACH, "ESL Deutschland", chESLDeutsch,
		"2012-09-01", "alternate", "dotpixels", 3, nil, 101, []string{"kUkruN8tu4U", "VdHE-rRmHuE"}, langGerman, noteResultOpen},
	{"eps-summer-2012-source", "ESL Pro Series Summer Finals 2012 (CS:Source)", 2012, "Deutschland", "css", catDACH, "ESL Deutschland", chESLDeutsch,
		"2012-09-01", "mortal-teamwork", "n-faculty", 3, score(2, 1), 144, []string{"oO2GYAFjLX0", "FPVlvpg6FNg", "2E-YaN4Zhto"}, langGerman, ""},
	{"eps-spring-2014", "ESL Pro Series Spring 2014", 2014, "Deutschland", "csgo", catDACH, "ESL", chESL,
		"", "mousesports", "team-wild-fire", 3, nil, 133, []string{"FA-RZoQqCSg", "3ERnzg8mlXI"}, langGerman, noteResultOpen},
	{"eps-germany-summer-2014", "ESL Pro Series Germany Summer 2014", 2014, "Deutschland", "csgo", catDACH, "ESL", chESL,
		"", "planetkey-dynamics", "mousesports", 3, score(1, 2), 148, []string{"GtZxiroCkso", "ZSxLYH1E27s", "zxdX1YCmceg"}, langGerman, ""},
	{"eps-finals-winter-2014", "ESL Pro Series Finals Winter 2014", 2014, "Deutschland", "csgo", catDACH, "ESL", chESL,
		"", "mousesports", "planetkey-dynamics", 3, nil, 89, []string{"Kz536sTE0OA"}, langGerman, noteResultOpen},
	{"esl-meisterschaft-spring-2016", "ESL Meisterschaft Frühling 2016", 2016, "Deutschland", "csgo", catDACH, "ESL Deutschland", chESLDeutsch,
		"", "leisure", "alternate-attax", 3, nil, 186, []string{"9INJXsB0TCg"}, langGerman, noteFullShow + "; " + noteResultOpen},
	{"esl-meisterschaft-winter-2016", "ESL Meisterschaft Winter 2016", 2016, "Deutschland", "csgo", catDACH, "ESL", chESLArchives,
		"", "alternate-attax", "euronics-gaming", 3, nil, 71, []string{"FMNDWEygNGg", "JU2P3WUk64g"}, langGerman, noteResultOpen},
	{"esl-meisterschaft-2020-s1", "ESL Meisterschaft 2020 – Saison 1", 2020, "Online", "csgo", catDACH, "ESL Deutschland", chESLDeutsch,
		"", "sprout", "big", 3, nil, 213, []string{"FWLjfawkozI"}, langGerman, noteFullShow + "; " + noteResultOpen},
	{"99damage-liga-season-9", "99Damage Liga Saison 9", 2018, "Deutschland", "csgo", catDACH, "99Damage", ch99Damage,
		"", "alternate-attax", "euronics-gaming", 3, nil, 149, []string{"bRz6L9yxKUk", "ytKSrRv1H6M", "J-Y51NAC0OE"}, langGerman, noteResultOpen},
	{"99damage-liga-season-10", "99Damage Liga Saison 10", 2018, "Deutschland", "csgo", catDACH, "99Damage", ch99Damage,
		"", "panthers-gaming", "expert-esport", 3, nil, 148, []string{"83XULIsBo5w", "XqXBDmbpDfM", "qDf0VphAfgY"}, langGerman, noteResultOpen},
	{"99damage-liga-season-11", "99Damage Liga Saison 11", 2019, "Deutschland", "csgo", catDACH, "99Damage", ch99Damage,
		"", "berzerk", "expert-esport", 3, nil, 120, []string{"d-pDpNs2t5U", "MGoPmccrlh0"}, langGerman, noteResultOpen},
}

func seedEsports(ctx context.Context, st *store.Store) error {
	genre, err := ensureGenre(ctx, st, "eSports", "esports")
	if err != nil {
		return err
	}

	games := map[string]models.Game{}
	labels := map[string]string{}
	for _, g := range esportsGames {
		game, err := st.EnsureGame(ctx, g.Name, g.Slug)
		if err != nil {
			return err
		}
		games[g.Slug] = game
		labels[g.Slug] = g.Label
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

	for _, f := range csFinals {
		if _, err := st.GetMediaBySlug(ctx, f.Slug, false); err == nil {
			continue
		} else if err != store.ErrNotFound {
			return err
		}
		if err := seedFinalMatch(ctx, st, f, genre.ID, games[f.Game], labels[f.Game], teams); err != nil {
			return err
		}
		log.Printf("created eSports tournament %s", f.Title)
	}
	return nil
}

func seedFinalMatch(ctx context.Context, st *store.Store, f seedFinal, genreID string, game models.Game, gameLabel string, teams map[string]models.Team) error {
	teamA, okA := teams[f.TeamA]
	teamB, okB := teams[f.TeamB]
	if !okA || !okB {
		log.Printf("skipping %s: unknown team slug", f.Slug)
		return nil
	}
	poster := "/static/img/esports/" + f.Slug + "-poster.svg"

	// Descriptions must stay spoiler-free: no winner, no score.
	description := "Das Grand Final von " + f.Title + " (" + gameLabel + ", " + f.Location + ") als offizielles VOD von " +
		f.Organizer + ". Ergebnisse bleiben verborgen, solange der Spoiler-Schutz aktiv ist."

	media, err := st.CreateMedia(ctx, store.MediaInput{
		Title:        f.Title,
		Slug:         f.Slug,
		Description:  description,
		Type:         models.MediaSeries,
		ThumbnailURL: poster,
		BannerURL:    "/static/img/esports/" + f.Slug + "-banner.svg",
		ReleaseYear:  f.Year,
		AgeRating:    "USK 18",
		Language:     f.Language,
		Tags:         []string{"esports", "counter-strike", f.Game, f.Category},
		IsPublished:  true,
		GenreIDs:     []string{genreID},
	})
	if err != nil {
		return err
	}

	var playedOn *time.Time
	validFrom := time.Date(f.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	if f.PlayedOn != "" {
		t, err := time.Parse("2006-01-02", f.PlayedOn)
		if err != nil {
			return err
		}
		playedOn, validFrom = &t, t
	}

	holder := f.Organizer + " (offizieller YouTube-Kanal)"
	channel := f.ChannelURL
	notes := "Offizielles, einbettbares YouTube-Video des Veranstalters. Wiedergabe ausschließlich über den YouTube-Embed-Player, ohne Abo-Pflicht (YouTube API Developer Policies)."
	if f.Note != "" {
		notes += " " + f.Note + "."
	}
	if _, err := st.CreateRightsInfo(ctx, media.ID, store.RightsInfoInput{
		Source:        models.RightsOfficialEmbed,
		RightsHolder:  &holder,
		LicenseDocURL: &channel,
		ValidFrom:     validFrom,
		Notes:         &notes,
	}); err != nil {
		return err
	}

	series, err := st.CreateSeries(ctx, media.ID, f.Title, "Playoffs von "+f.Title+".")
	if err != nil {
		return err
	}
	season, err := st.CreateSeason(ctx, series.ID, 1, "Playoffs")
	if err != nil {
		return err
	}

	urls := make([]string, len(f.VideoIDs))
	for i, id := range f.VideoIDs {
		urls[i] = "https://www.youtube.com/watch?v=" + id
	}

	episode, err := st.CreateEpisode(ctx, season.ID, "Grand Final",
		"Grand Final: "+teamA.Name+" gegen "+teamB.Name+" (Best-of-"+strconv.Itoa(f.BestOf)+").",
		1, f.DurationMin, urls[0], poster, true)
	if err != nil {
		return err
	}

	in := store.MatchInput{
		TeamAID:        teamA.ID,
		TeamBID:        teamB.ID,
		Stage:          "Grand Final",
		BestOf:         f.BestOf,
		PlayedOn:       playedOn,
		ExtraVideoURLs: urls[1:],
	}
	if game.ID != "" {
		gameID := game.ID
		in.GameID = &gameID
	}
	if len(f.Score) == 2 {
		a, b := f.Score[0], f.Score[1]
		in.ScoreA, in.ScoreB = &a, &b
	}
	return st.UpsertMatch(ctx, episode.ID, in)
}
