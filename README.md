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


## 32-Bit-Systeme

StreamFlix nutzt Next.js, Prisma und native Build-/Datenbank-Komponenten. Für den produktiven Betrieb ist ein 64-Bit-System klar empfohlen. 32-Bit-Systeme haben wenig Adressraum und erhalten für moderne Node.js-Stacks nur eingeschränkt vorgebaute native Pakete.

Geeignete Software für einen 32-Bit-Testbetrieb:
- **Betriebssystem:** 32-Bit-Windows 10/11 nur für lokale Tests. 32-Bit-Linux ist für diesen Stack nicht zuverlässig, weil aktuelle Node.js- und native Paket-Binaries für `linux-ia32` nicht durchgängig verfügbar sind.
- **Node.js:** 32-Bit-Installer von Node.js 20 LTS oder neuer für Windows (`x86`/`ia32`).
- **Datenbank:** PostgreSQL möglichst extern betreiben, z. B. auf einem 64-Bit-Server, NAS, Docker-Host oder als Managed Database. Das entlastet den knappen Arbeitsspeicher des 32-Bit-Clients.
- **Prisma:** Die vorhandenen Skripte setzen auf 32-Bit-Systemen automatisch Prisma Binary Engines, damit keine Node-API-Engine mit inkompatibler nativer Bindung erzwungen wird.

Prüfung der lokalen Umgebung:
```bash
npm run doctor:32bit
```

Empfohlene Startsequenz auf einem unterstützten 32-Bit-Windows-Testsystem:
```bash
npm install
set NODE_OPTIONS=--max-old-space-size=1024
npm run doctor:32bit
npm run prisma:generate
npm run dev
```

Für Produktion, Stripe-Webhooks, Transcoding oder größere Medienkataloge sollte ein 64-Bit-Linux-Server verwendet werden. 32-Bit eignet sich hier nur als eingeschränkte Entwicklungs- oder Demo-Umgebung.

## Legale Nutzung
StreamFlix ist ausschließlich für eigene Videos, lizenzierte Filme/Serien, frei verwendbare Inhalte und Creator-/Partner-Content mit Nutzungsrechten vorgesehen. Piraterie, DRM-Bypass und illegale Quellen sind ausgeschlossen.

## Nächste Ausbaustufen
- HLS/adaptive Bitrate, signierte CDN-URLs und Transcoding-Pipeline
- Untertitel und mehrere Tonspuren
- E-Mail-Verifizierung und Passwort-Reset-Mailer
- Vollständige CRUD-Formulare für Serien/Staffeln/Episoden
- Erweiterte Analytics und Empfehlungssysteme
- Native Apps über dieselben APIs
