# StreamFlix

Professionelles Web-MVP für eine legale Netflix-ähnliche Streaming-Plattform mit Konten, Profilen, Medienkatalog, Watchlist, Videoplayer, Adminbereich und Stripe-Abos.

## Funktionen
- Registrierung, Login, Logout mit bcrypt und HTTP-only JWT-Cookie
- Rollen: USER und ADMIN
- Profile, Watchlist, Watch Progress und Resume-Funktion
- Filme, Serien, Staffeln, Episoden, Genres und Tags
- Dunkles, responsives Streaming-UI mit Hero und Content-Reihen
- Admin-Dashboard, Medienverwaltung und API für Nutzer/Abos/Statistiken
- Stripe Checkout, Customer Portal und Webhook-Verarbeitung
- Sicherheitsgrundlagen: serverseitige Validierung, Zugriffskontrolle, Webhook-Signatur, keine Secrets im Frontend

## Setup
```bash
npm install
cp .env.example .env
npx prisma migrate dev
npx prisma db seed
npm run dev
```

Admin-Testkonto nach Seed:
- E-Mail: `admin@streamflix.local`
- Passwort: `StreamFlix123!`

## Wichtige ENV-Variablen
Siehe `.env.example` für `DATABASE_URL`, `JWT_SECRET`, `NEXT_PUBLIC_APP_URL` und Stripe Keys/Price IDs.

## Legale Nutzung
StreamFlix ist ausschließlich für eigene Videos, lizenzierte Filme/Serien, frei verwendbare Inhalte und Creator-/Partner-Content mit Nutzungsrechten vorgesehen. Piraterie, DRM-Bypass und illegale Quellen sind ausgeschlossen.

## Nächste Ausbaustufen
- HLS/adaptive Bitrate, signierte CDN-URLs und Transcoding-Pipeline
- Untertitel und mehrere Tonspuren
- E-Mail-Verifizierung und Passwort-Reset-Mailer
- Vollständige CRUD-Formulare für Serien/Staffeln/Episoden
- Erweiterte Analytics und Empfehlungssysteme
- Native Apps über dieselben APIs
