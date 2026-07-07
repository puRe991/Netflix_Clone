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

Frühere Versionen dieses Projekts nutzten Prisma 5, dessen Query-Engine **grundsätzlich keine Binaries für 32-Bit-Architekturen** (`linux-ia32`, `windows-ia32`) veröffentlicht – unabhängig davon, ob `PRISMA_CLIENT_ENGINE_TYPE` auf `binary` oder `library` gesetzt wurde. Dadurch schlug jeder Datenbankzugriff auf 32-Bit-Systemen fehl.

Seit dem Umstieg auf **Prisma ORM 7** läuft die Query-Engine als WASM + reines TypeScript (keine native, architekturspezifische Engine mehr) und verbindet sich über einen JS-Datenbanktreiber (`@prisma/adapter-pg` + `pg`). Damit funktioniert StreamFlix nun auch auf 32-Bit-Systemen, sofern die folgenden Voraussetzungen erfüllt sind:

- **Betriebssystem:** 32-Bit-Windows 10/11. 32-Bit-Linux (`linux-ia32`) bleibt ungeeignet, weil aktuelle Node.js-Versionen dafür keine Builds mehr veröffentlichen.
- **Node.js:** Node.js **20.19+ oder 22.12+** (32-Bit-Windows-Installer). Node.js hat ab Version 23 die 32-Bit-Windows-Builds eingestellt – nicht auf Node 23/24+ aktualisieren.
- **Datenbank:** PostgreSQL möglichst extern betreiben, z. B. auf einem 64-Bit-Server, NAS, Docker-Host oder als Managed Database. Das entlastet den knappen Arbeitsspeicher des 32-Bit-Clients.

Prüfung der lokalen Umgebung:
```bash
npm run doctor:32bit
```

Empfohlene Startsequenz auf einem unterstützten 32-Bit-Windows-System:
```bash
npm install
set NODE_OPTIONS=--max-old-space-size=1024
npm run doctor:32bit
npm run prisma:generate
npm run dev
```

Für Produktion, Stripe-Webhooks, Transcoding oder größere Medienkataloge sollte weiterhin ein 64-Bit-Linux-Server verwendet werden. 32-Bit eignet sich nur als Entwicklungs- oder Demo-Umgebung.

## Legale Nutzung
StreamFlix ist ausschließlich für eigene Videos, lizenzierte Filme/Serien, frei verwendbare Inhalte und Creator-/Partner-Content mit Nutzungsrechten vorgesehen. Piraterie, DRM-Bypass und illegale Quellen sind ausgeschlossen.

## Nächste Ausbaustufen
- HLS/adaptive Bitrate, signierte CDN-URLs und Transcoding-Pipeline
- Untertitel und mehrere Tonspuren
- E-Mail-Verifizierung und Passwort-Reset-Mailer
- Vollständige CRUD-Formulare für Serien/Staffeln/Episoden
- Erweiterte Analytics und Empfehlungssysteme
- Native Apps über dieselben APIs
