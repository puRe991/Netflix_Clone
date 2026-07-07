# Deployment

1. PostgreSQL bei Neon, Supabase, Railway oder einem eigenen Server anlegen.
2. Binary für die Zielplattform bauen: `make build-all` oder gezielt `GOOS=<windows|linux> GOARCH=<386|amd64> CGO_ENABLED=0 go build -o dist/streamflix ./cmd/streamflix`.
3. Umgebungsvariablen aus `.env.example` setzen (`DATABASE_URL`, `JWT_SECRET`, `APP_URL`, Stripe Keys).
4. Binary starten — Migrationen laufen beim Start automatisch gegen `DATABASE_URL`.
5. Einmalig `go run ./cmd/seed` (oder das gebaute `seed`-Binary) gegen die Ziel-Datenbank ausführen.
6. Stripe Price IDs für Basic und Premium eintragen.
7. Stripe Webhook `/api/webhooks/stripe` mit Signatur-Secret (`STRIPE_WEBHOOK_SECRET`) konfigurieren — ohne dieses Secret verifiziert und akzeptiert der Server keine Webhook-Events.
8. Medien nur aus eigenen, lizenzierten oder frei nutzbaren Quellen veröffentlichen.

## 32-Bit-Zielsysteme

`GOOS=windows GOARCH=386 CGO_ENABLED=0 go build ...` bzw. `GOOS=linux GOARCH=386 ...` erzeugen ein einzelnes statisch gelinktes Binary ohne jede weitere Laufzeitabhängigkeit. Auf dem 32-Bit-System selbst wird nichts weiter installiert — nur das Binary plus eine erreichbare PostgreSQL-Datenbank (kann und sollte extern laufen).
