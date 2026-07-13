# Konzept: Adaptive-Streaming-Pipeline

Aktuell trägt ein Admin im Formular (`/admin/media/new`) eine fertige `videoUrl` ein, die im Browser 1:1 als `<video src="…">` landet (`web/templates/pages/watch.html`, `internal/httpserver/api_watch.go`). Das funktioniert für Demos, hat aber drei strukturelle Grenzen:

- **Eine feste Datei/Qualität** — kein Wechsel je nach Bandbreite, kein Mobil-/Desktop-Unterschied.
- **Kein Schutz** — jede erreichbare URL kann direkt heruntergeladen werden.
- **Keine Verteilung** — die Datei liegt genau dort, wo die URL hinzeigt; bei vielen gleichzeitigen Zuschauern skaliert das nicht.

Dieses Dokument beschreibt, welche Komponenten für einen produktionsnahen Ausbau nötig wären, und in welcher Reihenfolge man sie sinnvoll einführt. Es enthält bewusst noch keinen Code — das ist eine Grundlage für eine spätere Entscheidung.

## Leitplanke: bestehende Architektur respektieren

Laut `docs/architecture.md` ist das Ziel ein **einzelnes, cgo-freies Go-Binary**, das auch auf 32-Bit-Windows/Linux läuft. Das schließt aus, Transcoding (ffmpeg o. ä.) *in* den `streamflix`-Prozess einzubetten — das wäre eine native, nicht cross-kompilierbare Abhängigkeit. Transcoding muss deshalb als **externer Schritt** laufen (separates Skript, Worker-Prozess oder Cloud-Dienst), der nur über HTTP/Storage mit dem bestehenden Backend spricht. Der `streamflix`-Server bleibt reiner Ausliefer- und Metadaten-Layer.

## Zielkomponenten

**1. Objektspeicher (Storage/Origin)**
Ein S3-kompatibler Bucket (AWS S3, Cloudflare R2, Backblaze B2, MinIO selbst gehostet) als Ablage für Rohmaster und die erzeugten HLS-Segmente. Ersetzt die heutige Praxis, irgendeine externe URL einzutragen.

**2. Transcoding (extern, offline)**
Ein `ffmpeg`-basierter Schritt, der eine hochgeladene Master-Datei in mehrere Qualitätsstufen (z. B. 480p/720p/1080p, ggf. 4K) plus ein HLS-Manifest (`.m3u8` + `.ts`/`.m4s`-Segmente) umwandelt. Läuft **nicht** im `streamflix`-Binary, sondern:
- als eigenständiges CLI/Makefile-Target, das ein Admin manuell anstößt, oder
- als separater Worker (eigenes kleines Go- oder Shell-Programm), der einen Upload-Ordner/Queue abarbeitet und das Ergebnis in den Objektspeicher schreibt.

**3. Packaging & Manifest**
Ergebnis des Transcodings ist pro Titel ein Ordner mit Master-Playlist (`master.m3u8`) plus einer Playlist je Qualitätsstufe. Das ersetzt das heutige `Media.VideoURL`-Feld: statt einer einzelnen Datei-URL zeigt es künftig auf die Master-Playlist.

**4. Player (Frontend)**
`web/static/js/watch.js` bräuchte `hls.js` (oder nativen HLS-Support in Safari), um die Master-Playlist zu laden und je nach Bandbreite automatisch zwischen den Qualitätsstufen zu wechseln. Fallback auf das heutige Verhalten (direktes `<video src>`) bleibt sinnvoll für Trailer oder Alt-Content ohne HLS-Aufbereitung.

**5. Verteilung (CDN)**
Vor den Objektspeicher gehört ein CDN (Cloudflare, CloudFront, Bunny CDN), das die Segmente in der Nähe der Zuschauer cached. Ohne CDN skaliert der Origin-Speicher bei vielen gleichzeitigen Streams nicht.

**6. DRM (optional, nur bei echtem Lizenzinhalt nötig)**
Für frei nutzbare/eigene Inhalte (aktueller Anspruch laut `docs/architecture.md`) nicht zwingend. Erst relevant, falls jemals fremdlizenzierter Content mit Studio-Vorgaben dazukommt — dann Widevine/FairPlay über einen DRM-as-a-Service-Anbieter, kein Eigenbau.

## Auswirkung auf das Datenmodell

- `models.Media.VideoURL *string` → Bedeutung ändert sich von "direkte Videodatei" zu "URL der HLS-Master-Playlist". Feldname/Typ bliebe gleich, nur die Admin-Beschreibung im Formular müsste angepasst werden.
- `models.Episode.VideoURL string` → analog.
- Kein neues DB-Feld zwingend nötig für Phase 1; erst falls mehrere Sprachfassungen/Untertitelspuren pro Titel abgebildet werden sollen, bräuchte es eine eigene Tabelle statt eines einzelnen URL-Strings.

## Empfohlene Reihenfolge (inkrementell, jede Phase einzeln nutzbar)

1. **Player-Vorbereitung**: `hls.js` einbinden, `watch.js` erkennt `.m3u8`-Endungen und nutzt adaptives Streaming, alles andere läuft wie bisher weiter. Kein Risiko für bestehende Inhalte.
2. **Manuelles Transcoding-Skript**: Ein `ffmpeg`-Aufruf (z. B. als `make transcode file=movie.mp4`), der lokal eine HLS-Variante erzeugt. Admin lädt Ergebnis manuell irgendwohin hoch und trägt die `master.m3u8`-URL wie bisher als `videoUrl` ein — noch kein Storage-Anbindung im Code.
3. **Objektspeicher-Anbindung**: Erst wenn Schritt 2 sich bewährt hat, das Transcoding-Skript direkt in einen S3-kompatiblen Bucket schreiben lassen und ggf. einen Upload-Button im Admin-Formular ergänzen, der den Transcoding-Job anstößt statt nur eine URL entgegenzunehmen.
4. **CDN davor**: Sobald echte Zuschauerzahlen das rechtfertigen.

Jede Phase ist unabhängig nutzbar und wirft nichts aus vorherigen Phasen um — insbesondere Phase 1 lässt sich isoliert umsetzen, ohne dass Storage oder Transcoding existieren müssen.
