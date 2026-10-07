# eSports: Counter-Strike-VODs (2006–2026)

Stand der Recherche: 7. Oktober 2026. Der Seed (`go run ./cmd/seed`) legt alle Turniere aus den Tabellen 1, 3 und 4 automatisch an: **46 Grand Finals** aus CS 1.6, CS:Source, CS:GO und CS2. Er kann auf bestehenden Datenbanken erneut laufen und legt nur fehlende Turniere an.

## Abbildung im Katalog

| eSports | StreamFlix |
|---|---|
| Turnier | Medientitel vom Typ `SERIES` (Genre „eSports“, Tags `esports`, `counter-strike`, `csgo`/`cs2`) |
| Phase (z. B. Playoffs) | Staffel |
| Match | Episode plus `matches`-Zeile (Teams, Format, Ergebnis, Datum, weitere VOD-Teile) |
| VOD-Teile (Map 1, Map 2, …) | `episodes.video_url` und danach `matches.extra_video_urls` |

## Welche Videos wir einbinden – und welche nicht

Eingebunden werden **ausschließlich Uploads auf dem offiziellen YouTube-Kanal des Veranstalters** (ESL, ESL Archives, ESL Deutschland/„ESL Meisterschaft“, PGL, BLAST, StarLadder, FACEIT, MLG, World Cyber Games, ESEA, 99Damage). Re-Uploads von Fan-Kanälen, Restreams von Streamern („Watch Parties“) und Team-Kanäle sind ausgeschlossen, auch wenn sie oft besser auffindbar sind. Rechteinhaber der Übertragung ist der Veranstalter, nicht der Uploader.

Jede ID wurde über YouTubes oEmbed-Endpunkt geprüft. Ein HTTP 200 bedeutet: Das Video existiert und Einbetten ist erlaubt. Die Antwort enthält den Kanal (`author_url`), der unten jeweils angegeben ist.

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
| World Cyber Games 2005 (CS:Source) | Team3D vs. k23 | Offizielle WCG-Uploads vorhanden, aber zwei Teile mit identischem Titel und ohne erkennbare Reihenfolge |
| World Cyber Games 2011 | ESC vs. SK Gaming | Nur Fan-Uploads |
| DreamHack-Finals CS 1.6 (2006–2012) | – | Nur Fan- und Team-Uploads |
| IEM European Championship (Fnatic vs. mTw) | – | Offizieller ESL-Archives-Upload, aber Saison und Jahr ließen sich nicht eindeutig zuordnen |
| ESEA LAN Season 11 (CS:Source) | Dynamic vs. Fully Torqued | Offizielle 6 Teile, aber das Format (Bracket-Reset) ist nicht sicher belegt |
| 99Damage Liga Saison 12–16 | – | Offline-Finals mit mehreren Begegnungen, die eigentliche Finalpaarung war nicht eindeutig |

Für diese Turniere wäre eine Lizenzanfrage beim Veranstalter (DreamHack/ESL FACEIT Group, Warner Bros. Discovery für ELEAGUE, Samsung für die WCG-Archive) der richtige Weg.

## Spoiler-Schutz

- Neue Profile haben den Spoiler-Schutz standardmäßig **an** (`profiles.hide_spoilers`). Anonyme Besucher steuern ihn per Cookie.
- Bei aktivem Schutz gibt es kein Ergebnis, keine Siegermarkierung, keine Episodendauer (sie würde die Anzahl der Maps verraten), keine Zeitleiste und keine Gesamtlänge im Player. Das Ergebnis lässt sich pro Match gezielt aufdecken („Ergebnis anzeigen“).
- Mehrteilige VODs laufen automatisch nacheinander. Der nächste Teil wird erst nach dem Ende des vorherigen geladen, die Anzahl der Teile wird vorher nie angezeigt.
- Grenzen: YouTube zeigt beim Pausieren bzw. am Ende eigene Vorschläge (mit `rel=0` nur vom selben Kanal) und im Player den Videotitel. Beides darf laut YouTube-Richtlinien nicht überdeckt werden. Die Videotitel der eingebundenen Uploads enthalten keine Ergebnisse.

## IDs erneut prüfen

```bash
# Alle IDs aus dem Seed:
grep -oE '"[A-Za-z0-9_-]{11}"' cmd/seed/esports.go | tr -d '"' | sort -u > /tmp/ids
# … oder die Major-IDs direkt:
for id in do4FUYUHmUM -5laY7bjono cqmdRLNRZHA 02I5vVxlJhU _Ef6cXH1oC8 7nWxZin959o AK5RY8W-AcM TC6N1zHvXpU \
  ttbe_A8Ee50 sIQ1Eh11Quk w2MdCEZSBtw z7IDPdmYw0A 9CcrjvC4SkE kgitmggEgrA A_wk85oA0Rc c2uQV59-tbI hkVQXKCVSFQ \
  UKYnAujcFac pogK6a25LO0 O9u1FQKZHVw CMBAVBkC8Zw 9pscaxA1QBs oVH782n-32U y6Y0l1Y-PTo 09N-bQbLrx8 L9eBTSwDwaM \
  vdXveK6G70g Vv7eXxl--LA meijWUaxNYg 9qqiaW7sjD0 p6U28DEBn8A; do
  printf '%s ' "$id"; curl -s -o /dev/null -w '%{http_code}\n' "https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=$id&format=json"
done   # 200 = ok, 401 = Einbetten deaktiviert, 404 = gelöscht/privat
```
