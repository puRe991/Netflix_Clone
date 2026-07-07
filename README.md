# StreamFlix

Professionelles Web-MVP für eine legale Netflix-ähnliche Streaming-Plattform mit Konten, Profilen, Medienkatalog, Watchlist, Videoplayer, vollständigem Adminbereich und Stripe-Abos — als einzelnes, selbst-enthaltenes Go-Binary.

## Warum Go statt Next.js/Prisma

Frühere Versionen dieses Projekts liefen auf Next.js + Prisma. Für 32-Bit-Systeme stellte sich das als grundsätzliche Sackgasse heraus:

- Node.js veröffentlicht keine `linux-ia32`-Builds mehr und hat `win32-ia32` ab Node 23 eingestellt.
- Prismas native Query- und Schema-Engines wurden nie für 32-Bit-Architekturen veröffentlicht — unabhängig vom Engine-Typ (`binary`/`library`).

Go dagegen kompiliert `GOOS=windows GOARCH=386` und `GOOS=linux GOARCH=386` offiziell und mit `CGO_ENABLED=0` zu einem einzelnen statisch gelinkten Binary **ohne jede Laufzeitabhängigkeit** — kein Node.js, keine native Engine, nichts zusätzlich zu installieren. Genau das macht dieses Projekt jetzt aus: Go-Server, serverseitig gerendertes HTML (`html/template`), PostgreSQL über den reinen Go-Treiber `pgx` (kein cgo), Templates/CSS/JS im Binary eingebettet (`embed`).

## Funktionen

- Registrierung, Login, Logout mit bcrypt und HTTP-only JWT-Cookie, CSRF-Schutz auf allen mutierenden Formularen
- Rollen: USER und ADMIN
- Profile, Watchlist, Watch Progress und Resume-Funktion (race-freier Upsert über partielle Unique-Indexe)
- Filme, Serien, Staffeln, Episoden, Genres und Tags
- Dunkles, responsives Streaming-UI mit Hero und Content-Reihen
- Vollständiger Adminbereich: Dashboard, Medien-CRUD, Genre-CRUD, Serien/Staffel/Episoden-CRUD, Nutzerverwaltung (Rollenwechsel), Abo-Übersicht, Analytics
- Stripe Checkout, Customer Portal und signaturgeprüftes Webhook (kein unverifizierter Fallback)
- Abo-Status wird bei jeder sicherheitsrelevanten Prüfung frisch aus der Datenbank gelesen, nie aus dem (bis zu 7 Tage alten) Session-Cookie

## Setup

Voraussetzungen: Go 1.24+, PostgreSQL.

```bash
cp .env.example .env
# DATABASE_URL, JWT_SECRET etc. in .env anpassen
go run ./cmd/seed      # legt Schema an (Migrationen laufen automatisch) und seeded Demo-Daten
go run ./cmd/streamflix
```

Admin-Testkonto nach dem Seed:
- E-Mail: `admin@streamflix.local`
- Passwort: `StreamFlix123!`

## Wichtige ENV-Variablen

Siehe `.env.example`: `DATABASE_URL`, `JWT_SECRET`, `APP_URL`, `APP_ENV`, `PORT`, Stripe Keys/Price IDs, `STRIPE_WEBHOOK_SECRET`.

`JWT_SECRET` ist außerhalb von `APP_ENV=development` zwingend erforderlich — der Server startet ohne gesetztes Secret in Produktion gar nicht erst. `STRIPE_WEBHOOK_SECRET` wird für jede Webhook-Verifikation gebraucht; ohne Secret schlägt die Signaturprüfung immer fehl (kein unsicherer Klartext-Fallback).

## Bauen (auch für 32-Bit)

```bash
make build          # natives Binary in dist/
make build-all       # windows/386, linux/386, windows/amd64, linux/amd64, darwin/arm64
```

Jedes erzeugte Binary ist einzeln lauffähig — keine weiteren Dateien, kein Node.js, keine DLLs. Auf dem 32-Bit-Zielsystem wird nur das jeweilige Binary plus eine erreichbare PostgreSQL-Datenbank benötigt.

## Deployment

1. PostgreSQL bei einem Anbieter der Wahl anlegen (Neon, Supabase, Railway, eigener Server, …).
2. Binary für die Zielplattform bauen (`make build-all` oder gezielt `GOOS=... GOARCH=... CGO_ENABLED=0 go build ./cmd/streamflix`).
3. Umgebungsvariablen aus `.env.example` setzen, Binary starten — Migrationen laufen beim Start automatisch.
4. `go run ./cmd/seed` einmalig gegen die Ziel-Datenbank ausführen (oder eigene Seed-Daten einspielen).
5. Stripe Price IDs für Basic und Premium eintragen, Webhook `/api/webhooks/stripe` mit Signatur-Secret konfigurieren.
6. Medien nur aus eigenen, lizenzierten oder frei nutzbaren Quellen veröffentlichen.

## Legale Nutzung

StreamFlix ist ausschließlich für eigene Videos, lizenzierte Filme/Serien, frei verwendbare Inhalte und Creator-/Partner-Content mit Nutzungsrechten vorgesehen. Piraterie, DRM-Bypass und illegale Quellen sind ausgeschlossen.

## Projektstruktur

```
cmd/streamflix/        Server-Entrypoint
cmd/seed/               Demo-Daten-Seed
internal/config/        Env-Konfiguration
internal/db/             Postgres-Pool + eingebettete SQL-Migrationen
internal/models/         Domänentypen
internal/store/          Handgeschriebenes SQL pro Aggregat (kein ORM)
internal/auth/           Sessions (JWT-Cookie), Passwort-Hashing, CSRF
internal/billing/        Stripe-Integration
internal/httpserver/     Routing, Handler, Validierung, Template-Rendering
web/templates/           html/template-Seiten und Partials
web/static/              CSS/JS, ins Binary eingebettet
```
