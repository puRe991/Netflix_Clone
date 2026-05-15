# Deployment

1. PostgreSQL bei Neon, Supabase oder Railway anlegen.
2. Variablen aus `.env.example` in Vercel setzen.
3. `npm install`, `npx prisma migrate deploy`, `npx prisma db seed` ausführen.
4. Stripe Price IDs für Basic und Premium eintragen.
5. Stripe Webhook `/api/webhooks/stripe` mit Signatur-Secret konfigurieren.
6. Medien nur aus eigenen, lizenzierten oder frei nutzbaren Quellen veröffentlichen.
