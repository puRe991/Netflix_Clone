# StreamFlix Architektur

StreamFlix ist ein legales Streaming-MVP auf Basis von Next.js, Prisma und PostgreSQL. Die App nutzt serverseitige Zugriffskontrollen, rollenbasierte Adminbereiche und Stripe Billing. Videoquellen sind eigene, lizenzierte oder frei verwendbare URLs; spätere signierte URLs, HLS, CDN, Transcoding und DRM über offizielle Anbieter sind vorbereitet.

## Schichten
- UI: App Router, Tailwind, responsive Netflix-ähnliche Komponenten.
- API: Next.js Route Handler als REST-Endpunkte.
- Daten: Prisma ORM mit PostgreSQL.
- Auth: E-Mail/Passwort, bcrypt, HTTP-only JWT-Cookie.
- Billing: Stripe Checkout, Customer Portal und Webhook-Synchronisation.
