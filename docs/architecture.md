# StreamFlix Architektur

StreamFlix ist ein legales Streaming-MVP auf Basis von Go und PostgreSQL, gebaut als einzelnes, selbst-enthaltenes Binary ohne Laufzeitabhängigkeiten (kein Node.js, keine native Datenbank-Engine). Die App nutzt serverseitige Zugriffskontrollen, rollenbasierte Adminbereiche und Stripe Billing. Videoquellen sind eigene, lizenzierte oder frei verwendbare URLs.

## Schichten
- UI: serverseitig gerenderte `html/template`-Seiten, dunkles, responsives Netflix-ähnliches Design, minimales Vanilla-JS für Video-Fortschritt und Live-Suche.
- HTTP: `net/http` mit method+pattern Routing (Go 1.22+), zentrale Fehlerbehandlung (`internal/httpserver/errors.go`) mit sauberen 401/403/404/422-Antworten statt roher 500er.
- Daten: handgeschriebenes SQL pro Aggregat über `pgx` (reiner Go-Treiber, kein cgo) — kein ORM.
- Auth: E-Mail/Passwort, bcrypt, HTTP-only JWT-Cookie, CSRF-Schutz (Double-Submit-Cookie) auf allen mutierenden Formularen. Abo-Status wird nie aus dem Cookie vertraut, sondern bei jeder sicherheitsrelevanten Prüfung frisch aus der Datenbank gelesen.
- Billing: Stripe Checkout, Customer Portal und signaturgeprüftes Webhook (kein unverifizierter Fallback, falls das Secret fehlt — der Server lehnt in dem Fall jedes Webhook-Event ab).

## Warum kein ORM/Framework
Ziel ist ein Binary, das garantiert auf 32-Bit-Windows und 32-Bit-Linux läuft (`GOOS=windows|linux GOARCH=386`, `CGO_ENABLED=0`). Jede zusätzliche Abhängigkeit mit nativem Code (cgo, vorkompilierte Engine-Binaries) würde dieses Ziel wieder gefährden. Alle verwendeten Third-Party-Pakete (`pgx`, `golang-jwt`, `bcrypt`, `google/uuid`, `stripe-go`) sind reines Go und cross-kompilieren nachweislich für `386`.
