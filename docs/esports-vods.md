# eSports: Counter-Strike-VOD-Archiv (2005–2026)

Stand der Recherche: 7. Oktober 2026. Der Katalog `cmd/seed/esports_catalog.json` enthält **58 Turniere mit 2.408 Episoden (2.254 Matches) und 3.661 offiziellen Videos** in CS 1.6, CS:Source, CS:GO und CS2. Jedes Turnier ist so vollständig wie die Veranstalter es veröffentlicht haben, von den Qualifikationen (Minors, RMRs, Closed Qualifiers, Online-Cups) über Gruppen- und Swiss-Phasen bis zum Grand Final.

`go run ./cmd/seed` importiert den Katalog. Auf bestehenden Datenbanken ergänzt der Import fehlende Phasen und Matches, sortiert die Staffeln neu und lässt bestehende Episoden, Admin-Änderungen und Wiedergabefortschritte unangetastet. Ein zweiter Lauf ändert nichts.

## Abbildung im Katalog

| eSports | StreamFlix |
|---|---|
| Turnier | Medientitel vom Typ `SERIES` (Genre „eSports“, Tags `esports`, `counter-strike`, `csgo`/`cs2`) |
| Phase | Staffel in fester Reihenfolge: Qualifikation, Online-Cups, Opening Stage, Elimination Stage, Challengers Stage, Gruppenphase, Legends Stage, Playoffs, Relegation, Weitere Matches |
| Match | Episode plus `matches`-Zeile (Teams, Format, Ergebnis, Datum, weitere VOD-Teile) |
| VOD-Teile (Map 1, Map 2, …) | `episodes.video_url` und danach `matches.extra_video_urls` |

## Welche Videos wir einbinden – und welche nicht

Eingebunden werden **ausschließlich Uploads auf dem offiziellen YouTube-Kanal des Veranstalters** (ESL, ESL Archives, ESL Deutschland/„ESL Meisterschaft“, PGL, BLAST, StarLadder, FACEIT, MLG, World Cyber Games, ESEA, 99Damage). Re-Uploads von Fan-Kanälen, Restreams von Streamern („Watch Parties“) und Team-Kanäle sind ausgeschlossen, auch wenn sie oft besser auffindbar sind. Rechteinhaber der Übertragung ist der Veranstalter, nicht der Uploader.

Jede der 3.661 IDs wurde über YouTubes oEmbed-Endpunkt geprüft. Ein HTTP 200 bedeutet: Das Video existiert und Einbetten ist erlaubt. Alle 3.661 lieferten 200. Zusätzlich wurde für jedes Video der hochladende Kanal (`author_url`) geprüft, denn offizielle Playlists können auch fremde Videos enthalten. Alle stammen von einem dieser 14 Kanäle: `@esl`, `@ESLCS`, `@ESLArchives`, `@ESLDeutschland`, `@pgl`, `@PGL_VODs`, `@BLASTPremier`, `@starladder_cs`, `@starladder_cs_highlights`, `@FACEIT`, `@MajorLeagueGaming`, `@WorldCyberGamesOfficial`, `@ESEA`, `@99DMGCSGO`. `@PGL_VODs` wird vom PGL-Hauptkanal selbst verlinkt. Der Test `cmd/seed/esports_test.go` erzwingt diese Liste.

## Wie der Katalog entstanden ist

1. **Sammeln:** Für jedes Turnier wurden die Kanalsuche und die offiziellen Turnier-Playlists des Veranstalterkanals ausgelesen, nicht die allgemeine YouTube-Suche. Wegen schwankender Suchergebnisse lief das zweimal, die Ergebnisse wurden zusammengeführt.
2. **Filtern:** Highlights, Talkshows, Interviews, Showmatches, andere Spiele (Hearthstone, StarCraft, LoL, Black Ops …), Videos anderer Jahrgänge, fremdsprachige Doppel-Uploads und 99Damage-Co-Streams fremder Ligen (z. B. ESL Pro League) fliegen raus.
3. **Zuordnen:** Aus den Titeln werden Phase, Tag/Runde/Gruppe, Paarung und Map-Nummer gelesen. Fehlt die Phase im Titel, gilt der Name der offiziellen Playlist. Videos derselben Paarung und Phase werden zu einem Match mit mehreren Teilen zusammengefasst. Eine wiederholte Map-Nummer beginnt ein neues Match.
4. **Teams:** Schreibweisen werden zusammengeführt (z. B. „NaVi“, „Na'Vi“ und „Natus Vincere“ oder „Alernate aTTaX“ und „ALTERNATE aTTaX“). Eine Paarung wird **nur dann zum Match mit Teamseiten, wenn beide Teams eindeutig erkannt sind.** WCG-Nationalteams („POL vs. CAN“), Platzhalter („TBA“) und Streams mit mehreren Begegnungen („AHG & BIG vs. Sprout“) werden normale Episoden.
5. **Übertragungen:** Ganztägige Streams (über 5 Stunden) werden nur dort als „Übertragung“ übernommen, wo es in der Phase keine Einzelmatch-Videos gibt, also vor allem bei BLAST. Das vermeidet Dubletten.
6. **WCG 2005–2010** wurde wegen sehr uneinheitlicher Titel von Hand erfasst. Doppel-Uploads desselben Spiels wurden zusammengeführt, und die koreanischen Uploads („결승 1라운드“ = Finale, 1. Runde) halfen bei der Reihenfolge.
7. **Format und Ergebnis:** Das Serienformat (Bo1/Bo3/Bo5) wird nur gesetzt, wenn Titel oder Map-Nummern es belegen, sonst steht „–“. Ergebnisse gibt es nur bei den verifizierten Grand Finals (Tabellen 1, 3, 4). Alle anderen Matches haben bewusst kein Ergebnis, statt eines geratenen.

## Archiv-Übersicht (aus dem Katalog erzeugt)

| Turnier | Spiel | Phasen (Episoden) | Matches | Videos | Kanäle |
|---|---|---|---|---|---|
| EMS One Katowice 2014 | CS:GO | Gruppenphase (8), Playoffs (7) | 15 | 23 | @esl |
| ESL One Cologne 2014 | CS:GO | Gruppenphase (16), Playoffs (7) | 23 | 32 | @ESLCS, @esl |
| ESL One Cologne 2015 | CS:GO | Qualifikation (38), Gruppenphase (20), Playoffs (7) | 65 | 78 | @ESLArchives |
| ESL One Katowice 2015 | CS:GO | Qualifikation (22), Gruppenphase (6), Playoffs (7) | 35 | 45 | @ESLCS, @esl |
| ESL One Cologne 2016 | CS:GO | Qualifikation (33), Gruppenphase (18), Playoffs (7) | 58 | 73 | @ESLArchives |
| MLG Major Championship: Columbus 2016 | CS:GO | Gruppenphase (39), Playoffs (7) | 31 | 52 | @MajorLeagueGaming |
| PGL Major Kraków 2017 | CS:GO | Gruppenphase (32), Playoffs (7) | 39 | 50 | @pgl |
| FACEIT Major: London 2018 | CS:GO | Challengers Stage (34), Legends Stage (38), Playoffs (15) | 87 | 88 | @FACEIT |
| IEM Katowice Major 2019 | CS:GO | Qualifikation (60), Challengers Stage (33), Gruppenphase (4), Legends Stage (33), Playoffs (11) | 141 | 238 | @ESLArchives |
| StarLadder Major: Berlin 2019 | CS:GO | Qualifikation (196), Challengers Stage (33), Legends Stage (33), Playoffs (7) | 261 | 478 | @starladder_cs_highlights |
| PGL Major Stockholm 2021 | CS:GO | Challengers Stage (33), Legends Stage (29), Playoffs (7) | 69 | 114 | @pgl |
| IEM Rio Major 2022 | CS:GO | Qualifikation (157), Challengers Stage (33), Legends Stage (32), Playoffs (7) | 229 | 351 | @ESLArchives |
| PGL Major Antwerp 2022 | CS:GO | Challengers Stage (33), Legends Stage (7), Playoffs (5), Weitere Matches (37) | 71 | 128 | @pgl |
| BLAST.tv Paris Major 2023 | CS:GO | Qualifikation (27), Challengers Stage (7), Legends Stage (7), Playoffs (4) | 1 | 45 | @BLASTPremier |
| PGL CS2 Major Copenhagen 2024 | CS2 | Opening Stage (27), Elimination Stage (34), Playoffs (7) | 68 | 110 | @pgl |
| Perfect World Shanghai Major 2024 | CS2 | Opening Stage (33), Elimination Stage (39), Playoffs (1) | 73 | 73 | @PGL_VODs, @pgl |
| BLAST.tv Austin Major 2025 | CS2 | Opening Stage (8), Challengers Stage (7), Legends Stage (7), Playoffs (4), Weitere Matches (1) | 1 | 27 | @BLASTPremier |
| StarLadder Budapest Major 2025 | CS2 | Opening Stage (34), Challengers Stage (33), Legends Stage (33), Playoffs (7) | 107 | 107 | @starladder_cs |
| IEM Cologne Major 2026 | CS2 | Opening Stage (33), Challengers Stage (33), Legends Stage (26), Playoffs (14) | 106 | 106 | @ESLArchives |
| World Cyber Games 2005 | CS:Source | Gruppenphase (2), Playoffs (2), Weitere Matches (2) | 4 | 9 | @WorldCyberGamesOfficial |
| World Cyber Games 2006 | CS 1.6 | Gruppenphase (2), Playoffs (3), Weitere Matches (8) | 12 | 20 | @WorldCyberGamesOfficial |
| World Cyber Games 2007 | CS 1.6 | Gruppenphase (6), Playoffs (1) | 1 | 9 | @WorldCyberGamesOfficial |
| Intel Extreme Masters III – Global Challenge Dubai 2008 | CS 1.6 | Playoffs (1), Weitere Matches (11) | 12 | 12 | @ESLArchives |
| World Cyber Games 2008 | CS 1.6 | Gruppenphase (6), Playoffs (2) | 8 | 11 | @WorldCyberGamesOfficial |
| Intel Extreme Masters III – European Championship 2009 | CS 1.6 | Weitere Matches (10) | 10 | 10 | @ESLArchives |
| Intel Extreme Masters IV – Global Challenge Chengdu 2009 | CS 1.6 | Playoffs (1), Weitere Matches (6) | 7 | 7 | @ESLArchives, @esl |
| Intel Extreme Masters IV – Global Challenge Dubai 2009 | CS 1.6 | Playoffs (4) | 4 | 4 | @esl |
| World Cyber Games 2009 | CS 1.6 | Qualifikation (3), Gruppenphase (17), Playoffs (3) | 23 | 30 | @WorldCyberGamesOfficial |
| Intel Extreme Masters IV – World Championship 2010 | CS 1.6 | Playoffs (1), Weitere Matches (10) | 9 | 11 | @ESLArchives, @esl |
| Intel Extreme Masters V – Global Challenge Cologne 2010 | CS 1.6 | Playoffs (1) | 1 | 1 | @esl |
| World Cyber Games 2010 | CS 1.6 | Gruppenphase (3), Playoffs (2) | 5 | 8 | @WorldCyberGamesOfficial |
| Intel Extreme Masters V – European Championship Finals 2011 | CS 1.6 | Playoffs (1) | 1 | 3 | @ESLArchives |
| Intel Extreme Masters V – World Championship 2011 | CS 1.6 | Playoffs (1) | 1 | 2 | @ESLArchives |
| Intel Extreme Masters VI – Global Challenge Guangzhou 2011 | CS 1.6 | Playoffs (1) | 1 | 3 | @ESLArchives |
| Intel Extreme Masters VI – Global Challenge New York 2011 | CS 1.6 | Playoffs (1), Weitere Matches (1) | 1 | 3 | @ESLArchives, @esl |
| ESEA LAN Season 10 | CS:Source | Playoffs (1), Weitere Matches (2) | 3 | 6 | @ESEA |
| ESEA LAN Season 11 | CS:Source | Playoffs (1), Weitere Matches (12) | 13 | 29 | @ESEA |
| Intel Extreme Masters VI – Global Challenge Kiev 2012 | CS 1.6 | Playoffs (1), Weitere Matches (1) | 1 | 2 | @ESLArchives |
| Intel Extreme Masters VI – World Championship 2012 | CS 1.6 | Playoffs (1) | 1 | 2 | @esl |
| RaidCall EMS One Fall 2013 | CS:GO | Online-Cups (24), Gruppenphase (16) | 40 | 41 | @esl |
| RaidCall EMS One Finals 2013 | CS:GO | Online-Cups (11), Gruppenphase (12), Playoffs (7) | 30 | 30 | @esl |
| RaidCall EMS One Summer 2013 | CS:GO | Online-Cups (11) | 11 | 11 | @esl |
| ESL Pro Series Summer Finals 2012 (CS 1.6) | CS 1.6 | Playoffs (1) | 1 | 2 | @ESLDeutschland |
| ESL Pro Series Summer Finals 2012 (CS:Source) | CS:Source | Playoffs (1) | 1 | 3 | @ESLDeutschland |
| ESL Pro Series Finals Winter 2014 | CS:GO | Playoffs (10) | 10 | 11 | @ESLDeutschland, @esl |
| ESL Pro Series Germany Summer 2014 | CS:GO | Gruppenphase (8), Playoffs (2), Weitere Matches (4) | 13 | 21 | @ESLDeutschland, @esl |
| ESL Pro Series Spring 2014 | CS:GO | Playoffs (5) | 5 | 11 | @ESLDeutschland, @esl |
| ESL Meisterschaft Frühling 2016 | CS:GO | Playoffs (3) | 3 | 3 | @ESLDeutschland |
| ESL Meisterschaft Winter 2016 | CS:GO | Playoffs (3) | 3 | 7 | @ESLArchives |
| 99Damage Liga Saison 10 | CS:GO | Gruppenphase (25), Playoffs (6), Relegation (18), Weitere Matches (29) | 78 | 162 | @99DMGCSGO |
| 99Damage Liga Saison 9 | CS:GO | Gruppenphase (33), Playoffs (6), Relegation (16), Weitere Matches (27) | 79 | 173 | @99DMGCSGO |
| 99Damage Liga Saison 11 | CS:GO | Gruppenphase (53), Playoffs (3), Relegation (10) | 66 | 134 | @99DMGCSGO |
| 99Damage Liga Saison 12 | CS:GO | Gruppenphase (58), Relegation (6), Weitere Matches (3) | 63 | 131 | @99DMGCSGO |
| 99Damage Liga Saison 13 | CS:GO | Gruppenphase (72), Playoffs (3), Relegation (4) | 79 | 154 | @99DMGCSGO |
| 99Damage Liga Saison 14 | CS:GO | Gruppenphase (81), Playoffs (3), Relegation (3) | 87 | 165 | @99DMGCSGO |
| 99Damage Liga Saison 15 | CS:GO | Gruppenphase (64), Playoffs (2), Relegation (4) | 66 | 137 | @99DMGCSGO |
| ESL Meisterschaft 2020 – Saison 1 | CS:GO | Gruppenphase (21), Playoffs (3), Weitere Matches (1) | 3 | 25 | @ESLDeutschland |
| 99Damage Liga Saison 16 | CS:GO | Gruppenphase (20), Playoffs (1) | 18 | 40 | @99DMGCSGO |

**Hinweise zu einzelnen Turnieren** (auch im Rechte-Vermerk des jeweiligen Titels hinterlegt):

- **World Cyber Games 2005 (CS:Source):** Das Finale Team3D gegen k23 gibt es in zwei offiziellen Teilen mit identischem Titel. Die Reihenfolge der Teile ist nicht gesichert. Laut Wikipedia gewann Team3D, ein Ergebnis ist nicht hinterlegt.
- **ESEA LAN Season 11 (CS:Source):** Das Grand Final Dynamic gegen Fully Torqued ist in sechs offiziellen Teilen eingebunden. Dynamic kam aus dem Lower Bracket. Das genaue Serienformat ist nicht belegt, deshalb steht es auf „–“.
- **99Damage Liga Saison 12–16:** Alle offiziellen Spieltage, Relegationen und Offline-Finals. Bei den Offline-Finals sagt der Videotitel nicht, welche Begegnung das Endspiel war. Sie stehen deshalb als „Offline-Finals“ ohne „Grand Final“-Markierung im Katalog.
- **IEM European Championships:** Die Finals 1/3–3/3 Fnatic gegen mTw gehören zur **IEM V European Championship 2011 (Kiew)**. Grundlage dafür sind die Platzierungen auf esportsearnings.com (1. Fnatic, 2. mTw) und der Upload-Zeitpunkt 2011. Die Uploads mit dem Titel „IEM 2009 Hanover Season II European Championship“ stehen als **IEM III European Championship 2009**. Die ESL betitelt sie als Season II, laut Wikipedia gehört das Turnier zu Season III.
- **Neu dazugekommen, weil die offiziellen Archive sie enthalten:** IEM III Global Challenge Dubai 2008 sowie RaidCall EMS One Summer und Fall 2013 (Online-Cups und Gruppen). Die Cups der Spring Season gehören zu „RaidCall EMS One Finals 2013“.

## Rechtliche Leitplanken (bitte vor Livegang anwaltlich prüfen lassen)

1. **Keine Paywall.** Die YouTube API Developer Policies verbieten es, Nutzern für das Ansehen im eingebetteten Player Geld zu berechnen oder den Zugang an andere Handlungen als den Klick auf „Play“ zu knüpfen: „API Clients must not charge users to watch content in an embedded YouTube player. API Clients must not otherwise gate access to a video by requiring a user to take an action other than clicking the play button“. Deshalb laufen Episoden, deren Videoteile *alle* YouTube-URLs sind, **ohne Abo und ohne Login**. Die Rechtequelle dafür ist `OFFICIAL_EMBED`. Selbst gehostete Videos bleiben wie bisher hinter dem Abo.
2. **Player nicht verändern oder überdecken.** Der Spoiler-Modus nutzt nur den dokumentierten Player-Parameter `controls=0` und die offizielle IFrame Player API für eigene Knöpfe (Pause, ±10 s/30 s, +2 min). YouTube-Logo, Titel, Werbung und Endscreen bleiben unangetastet.
3. **Datenschutz.** Der Player läuft über `youtube-nocookie.com` und wird erst nach Klick auf „Video abspielen“ geladen (Zwei-Klick-Lösung). Ein Hinweis steht auf `/legal/datenschutz`.
4. **Keine Werbung auf/um den Player verkaufen**, solange keine schriftliche Genehmigung von YouTube vorliegt.
5. Verschwindet ein Video oder wird es gesperrt, zeigt der Player einen Hinweis mit Link zu YouTube. Die IDs sollten regelmäßig erneut geprüft werden (siehe unten).

## Tabelle 1: eingebundene Grand Finals (19 Majors)

Die Reihenfolge der Teams folgt dem offiziellen Videotitel (Bracket-Reihenfolge) und verrät deshalb nicht den Sieger.

| Major | Grand Final | Ergebnis | Kanal | YouTube-IDs |
|---|---|---|---|---|
| EMS One Katowice 2014 | Virtus.pro vs. NiP | 2:0 | [ESL](https://www.youtube.com/@esl) | `do4FUYUHmUM` (komplettes Finale) |
| ESL One Cologne 2014 | NiP vs. Fnatic | 2:1 | [ESL Counter-Strike](https://www.youtube.com/@ESLCS) | `-5laY7bjono` (ESL Classics, komplett) |
| ESL One Katowice 2015 | Fnatic vs. NiP | 2:1 | [ESL Counter-Strike](https://www.youtube.com/@ESLCS) | `cqmdRLNRZHA` (ESL Classics, komplett) |
| ESL One Cologne 2015 | Fnatic vs. Team EnVyUs | 2:0 | [ESL Archives](https://www.youtube.com/@ESLArchives) | `02I5vVxlJhU`, `_Ef6cXH1oC8` |
| MLG Major Columbus 2016 | Natus Vincere vs. Luminosity | 0:2 | [Major League Gaming](https://www.youtube.com/@MajorLeagueGaming) | `7nWxZin959o` (komplett) |
| ESL One Cologne 2016 | SK Gaming vs. Team Liquid | 2:0 | [ESL Archives](https://www.youtube.com/@ESLArchives) | `AK5RY8W-AcM`, `TC6N1zHvXpU` |
| PGL Major Kraków 2017 | Gambit vs. Immortals | 2:1 | [PGL](https://www.youtube.com/@pgl) | `ttbe_A8Ee50`, `sIQ1Eh11Quk`, `w2MdCEZSBtw` |
| FACEIT Major London 2018 | Natus Vincere vs. Astralis | 0:2 | [FACEIT](https://www.youtube.com/@FACEIT) | `z7IDPdmYw0A`, `9CcrjvC4SkE` |
| IEM Katowice Major 2019 | ENCE vs. Astralis | 0:2 | [ESL Archives](https://www.youtube.com/@ESLArchives) | `kgitmggEgrA`, `A_wk85oA0Rc` |
| StarLadder Major Berlin 2019 | AVANGAR vs. Astralis | 0:2 | [StarLadder](https://www.youtube.com/@starladder_cs_highlights) | `c2uQV59-tbI`, `hkVQXKCVSFQ` |
| PGL Major Stockholm 2021 | G2 vs. Natus Vincere | 0:2 | [PGL](https://www.youtube.com/@pgl) | `UKYnAujcFac`, `pogK6a25LO0` |
| PGL Major Antwerp 2022 | FaZe vs. Natus Vincere | 2:0 | [PGL](https://www.youtube.com/@pgl) | `O9u1FQKZHVw`, `CMBAVBkC8Zw` |
| IEM Rio Major 2022 | Outsiders vs. Heroic | 2:0 | [ESL Archives](https://www.youtube.com/@ESLArchives) | `9pscaxA1QBs`, `oVH782n-32U` |
| BLAST.tv Paris Major 2023 | Vitality vs. GamerLegion | 2:0 | [BLAST Premier](https://www.youtube.com/@BLASTPremier) | `y6Y0l1Y-PTo` (Finaltag inkl. Vorberichten) |
| PGL CS2 Major Copenhagen 2024 | FaZe vs. Natus Vincere | 1:2 | [PGL](https://www.youtube.com/@pgl) | `09N-bQbLrx8`, `L9eBTSwDwaM`, `vdXveK6G70g` |
| Perfect World Shanghai Major 2024 | Team Spirit vs. FaZe | 2:1 | [PGL](https://www.youtube.com/@pgl) | `Vv7eXxl--LA` (inkl. Showmatch) |
| BLAST.tv Austin Major 2025 | Vitality vs. The MongolZ | 2:1 | [BLAST Premier](https://www.youtube.com/@BLASTPremier) | `meijWUaxNYg` (Finaltag komplett) |
| StarLadder Budapest Major 2025 | Vitality vs. FaZe (Bo5) | 3:1 | [StarLadder CS2](https://www.youtube.com/@starladder_cs) | `9qqiaW7sjD0` |
| IEM Cologne Major 2026 | FURIA vs. Team Falcons (Bo5) | 0:3 | [ESL Archives](https://www.youtube.com/@ESLArchives) | `p6U28DEBn8A` |

Die Ergebnisse wurden mit der Major-Übersicht auf [Wikipedia](https://en.wikipedia.org/wiki/Counter-Strike_Major_Championships) abgeglichen, die beiden jüngsten zusätzlich mit Turnierberichten ([Budapest](https://esportsworldcup.com/en/news/vitality-caps-dominant-cs2-season-with-starladder-major-victory), [Köln 2026](https://www.gfinityesports.com/article/team-falcons-become-iem-cologne-major-2026-champions)).

## Tabelle 3: Internationale Klassiker (CS 1.6, CS:Source, frühes CS:GO)

| Turnier | Grand Final | Ergebnis | Kanal | YouTube-IDs |
|---|---|---|---|---|
| World Cyber Games 2006 (Monza) | NiP vs. Pentagram | 1:2 | [WCG](https://www.youtube.com/@WorldCyberGamesOfficial) | `YSEtpv7n8tw`, `n89xjEA6VdE`, `h6gILRcMFVA` |
| World Cyber Games 2007 (Seattle) | emuLate vs. Team NoA | 2:1 | WCG | `DEQRNpv6fHM`, `gJsTe85XzNQ`, `FEQCLTTy5QY` |
| World Cyber Games 2008 (Köln) | SK Gaming vs. mTw | 1:2 | WCG | `J5DbuSDp5t8`, `P5-ubodo6A0`, `XLHlDEN7HHc` |
| World Cyber Games 2009 (Chengdu) | Fnatic vs. AGAiN | offen¹ | WCG | `AaUSjLtG4aM`, `kPt1zrKbviY` |
| World Cyber Games 2010 (Los Angeles) | mTw („Denmark2“) vs. Natus Vincere | 1:2 | WCG | `kVa3-QcDDlc`, `cVz_yRkxsUg`, `Q_JSeHJq9Cw` (englischer Kommentar) |
| IEM IV Global Challenge Chengdu 2009 | Fnatic vs. SK Gaming | offen¹ | [ESL](https://www.youtube.com/@esl) | `L3dYILqFunE` (ESL Classics, komplett) |
| IEM IV Global Challenge Dubai 2009 | Fnatic vs. Meet Your Makers | offen¹ | ESL | `hZeysLd0QSE` (ESL Classics, komplett) |
| IEM IV World Championship 2010 (Hannover) | Fnatic vs. Natus Vincere | offen¹ | ESL | `LzWsMTSdY3A` (ESL Classics, komplett) |
| IEM V Global Challenge Cologne 2010 | Fnatic vs. mousesports | offen¹ | ESL | `BJbo61X3D58` (ESL Classics, komplett) |
| IEM V World Championship 2011 (Hannover) | Natus Vincere vs. Frag eXecutors | 2:0 | [ESL Archives](https://www.youtube.com/@ESLArchives) | `oicR25CGIcc`, `LGNsKZsK9BI` |
| IEM VI Global Challenge Guangzhou 2011 | Fnatic vs. mousesports | 2:1 | ESL Archives | `YuALZm8OkCc`, `z7RxpTAYDZ4`, `Dhpo3TF9JQc` |
| IEM VI Global Challenge New York 2011 | SK Gaming vs. Fnatic | offen¹ | ESL | `71Kf35IrPKc`, `iw5dVBKnuyw` |
| IEM VI Global Challenge Kiev 2012 | Natus Vincere vs. SK Gaming | 2:0 | ESL Archives | `mMVk1YubLO0` (komplett) |
| IEM VI World Championship 2012 (Hannover) | ESC Gaming vs. Natus Vincere | 2:0 | ESL | `S82eTX3Wiq0`, `fiF42bT206o` |
| ESEA LAN Season 10 (Dallas, CS:Source) | Dynamic vs. Zomblerz | 2:1 | [ESEA](https://www.youtube.com/@ESEA) | `Y6A7ZgrekBE`, `vz0TjGTO-gM`, `uTgITPAW__8` |
| RaidCall EMS One Finals 2013 (CS:GO) | VeryGames vs. Virtus.pro | offen¹ | ESL | `N-GiZJbwD24` (komplett) |

Hinweis: Die ESL hat die Finals von 2011 und 2012 in den Videotiteln teils falsch datiert („IEM CeBIT 2010“ für NaVi vs. FX, das Finale der IEM V World Championship 2011). Im Katalog stehen die korrigierten Turniernamen.

## Tabelle 4: Deutsche Ligen (EPS, ESL Meisterschaft, 99Damage Liga) – deutscher Kommentar

| Turnier | Grand Final | Ergebnis | Kanal | YouTube-IDs |
|---|---|---|---|---|
| ESL Pro Series Summer Finals 2012 – CS 1.6 | ALTERNATE vs. dotpiXels | offen¹ | [ESL Deutschland](https://www.youtube.com/@ESLDeutschland) | `kUkruN8tu4U`, `VdHE-rRmHuE` |
| ESL Pro Series Summer Finals 2012 – CS:Source | mTw vs. n!faculty | 2:1 | ESL Deutschland | `oO2GYAFjLX0`, `FPVlvpg6FNg`, `2E-YaN4Zhto` |
| ESL Pro Series Spring 2014 (CS:GO) | mousesports vs. Team Wild Fire | offen¹ | ESL | `FA-RZoQqCSg`, `3ERnzg8mlXI` |
| ESL Pro Series Germany Summer 2014 | Planetkey Dynamics vs. mousesports | 1:2 | ESL | `GtZxiroCkso`, `ZSxLYH1E27s`, `zxdX1YCmceg` |
| ESL Pro Series Finals Winter 2014 | mousesports vs. Planetkey Dynamics | offen¹ | ESL | `Kz536sTE0OA` (komplett) |
| ESL Meisterschaft Frühling 2016 | LeiSuRe vs. ALTERNATE aTTaX | offen¹ | ESL Deutschland | `9INJXsB0TCg` (komplett) |
| ESL Meisterschaft Winter 2016 | ALTERNATE aTTaX vs. EURONICS Gaming | offen¹ | ESL Archives | `FMNDWEygNGg`, `JU2P3WUk64g` |
| ESL Meisterschaft 2020 – Saison 1 | Sprout vs. BIG | offen¹ | ESL Deutschland | `FWLjfawkozI` (komplett) |
| 99Damage Liga Saison 9 | ALTERNATE aTTaX vs. EURONICS Gaming | offen¹ | [99Damage](https://www.youtube.com/@99DMGCSGO) | `bRz6L9yxKUk`, `ytKSrRv1H6M`, `J-Y51NAC0OE` |
| 99Damage Liga Saison 10 | PANTHERS Gaming vs. expert eSport | offen¹ | 99Damage | `83XULIsBo5w`, `XqXBDmbpDfM`, `qDf0VphAfgY` |
| 99Damage Liga Saison 11 | Berzerk vs. expert eSport | offen¹ | 99Damage | `d-pDpNs2t5U`, `MGoPmccrlh0` |

¹ **Ergebnis offen:** Liquipedia und HLTV blockieren automatisierte Abfragen, und andere Quellen nannten für diese Finals kein eindeutiges Ergebnis. Statt zu raten, bleibt das Ergebnis leer. Im Katalog steht dann „nicht erfasst“. Nachtragen lässt es sich im Admin unter *Serien → Staffel → Episode bearbeiten → eSports-Match*. Wo ein Ergebnis eingetragen ist, stammt es aus einer Quelle (Wikipedia, Turnierberichte) oder folgt eindeutig aus dem belegten Sieger und den nummerierten Map-Uploads (z. B. „Set 1–3“ → 2:1). Ein Test (`cmd/seed/esports_test.go`) prüft, dass Ergebnis, Format und Anzahl der Map-Uploads zusammenpassen.

## Tabelle 5: Majors und Turniere ohne verwendbares offizielles VOD (nicht eingebunden)

| Major | Grand Final | Befund |
|---|---|---|
| DreamHack Winter 2013 (Jönköping) | Fnatic vs. NiP (2:1) | Nur Fan-Re-Uploads (z. B. „Swoier“, „CS Pro Play“) |
| DreamHack Winter 2014 (Jönköping) | Team LDLC vs. NiP (2:1) | Nur Fan-Re-Uploads. Der Kanal „Pgl Studio“ ist nicht eindeutig PGL zuzuordnen und wurde deshalb nicht genutzt |
| DreamHack Open Cluj-Napoca 2015 | Team EnVyUs vs. Natus Vincere (2:0) | Nur Fan-Re-Uploads |
| ELEAGUE Major Atlanta 2017 | Astralis vs. Virtus.pro (2:1) | Kein VOD auf einem ELEAGUE-Kanal (ELEAGUE ist eingestellt) |
| ELEAGUE Major Boston 2018 | Cloud9 vs. FaZe (2:1) | Nur der Team-Kanal „Cloud9 VODs“, der nicht Rechteinhaber der Übertragung ist |

| CPL, ESWC, WCG 2000–2004 und weitere Turniere vor ca. 2006 | – | YouTube gab es erst ab 2005. Von dieser Zeit existieren nur Fan-Uploads, Fragmovies und Demo-Aufnahmen, keine Uploads der Veranstalter. Der ESWC-Kanal hat keine CS-Finals aus dieser Zeit |
| World Cyber Games 2011 | ESC vs. SK Gaming | Nur Fan-Uploads |
| DreamHack-Finals CS 1.6 (2006–2012) | – | Nur Fan- und Team-Uploads |
| Einzelvideos ohne sichere Zuordnung | – | z. B. „SK Gaming vs. WinFakt – IEM GC New York“ (Jahr unklar) und „Knochen Executioners vs. wNv“ (WCG 2009, nur Teil 2). Dazu WCG-Uploads, die eindeutig Doppel-Uploads desselben Spiels sind |
| Spiele ohne offizielles VOD innerhalb eingebundener Turniere | – | Manche Maps oder Matches haben die Veranstalter nie hochgeladen. Solche Matches sind mit den vorhandenen Teilen eingebunden |

Für diese Turniere wäre eine Lizenzanfrage beim Veranstalter (DreamHack/ESL FACEIT Group, Warner Bros. Discovery für ELEAGUE, Samsung für die WCG-Archive) der richtige Weg.

## Spoiler-Schutz

- Neue Profile haben den Spoiler-Schutz standardmäßig **an** (`profiles.hide_spoilers`). Anonyme Besucher steuern ihn per Cookie.
- Bei aktivem Schutz gibt es kein Ergebnis, keine Siegermarkierung, keine Episodendauer (sie würde die Anzahl der Maps verraten), keine Zeitleiste und keine Gesamtlänge im Player. Das Ergebnis lässt sich pro Match gezielt aufdecken („Ergebnis anzeigen“).
- Mehrteilige VODs laufen automatisch nacheinander. Der nächste Teil wird erst nach dem Ende des vorherigen geladen, die Anzahl der Teile wird vorher nie angezeigt.
- Grenzen: YouTube zeigt beim Pausieren bzw. am Ende eigene Vorschläge (mit `rel=0` nur vom selben Kanal) und im Player den Videotitel. Beides darf laut YouTube-Richtlinien nicht überdeckt werden. Die Videotitel der eingebundenen Uploads enthalten keine Ergebnisse.

## IDs erneut prüfen

```bash
# Alle IDs aus dem Katalog:
python3 -c "import json;c=json.load(open('cmd/seed/esports_catalog.json'));print('\n'.join(sorted({v for t in c['tournaments'] for s in t['stages'] for e in s['episodes'] for v in e['videos']})))" > /tmp/ids
# … oder die Major-IDs direkt:
for id in do4FUYUHmUM -5laY7bjono cqmdRLNRZHA 02I5vVxlJhU _Ef6cXH1oC8 7nWxZin959o AK5RY8W-AcM TC6N1zHvXpU \
  ttbe_A8Ee50 sIQ1Eh11Quk w2MdCEZSBtw z7IDPdmYw0A 9CcrjvC4SkE kgitmggEgrA A_wk85oA0Rc c2uQV59-tbI hkVQXKCVSFQ \
  UKYnAujcFac pogK6a25LO0 O9u1FQKZHVw CMBAVBkC8Zw 9pscaxA1QBs oVH782n-32U y6Y0l1Y-PTo 09N-bQbLrx8 L9eBTSwDwaM \
  vdXveK6G70g Vv7eXxl--LA meijWUaxNYg 9qqiaW7sjD0 p6U28DEBn8A; do
  printf '%s ' "$id"; curl -s -o /dev/null -w '%{http_code}\n' "https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=$id&format=json"
done   # 200 = ok, 401 = Einbetten deaktiviert, 404 = gelöscht/privat
```
