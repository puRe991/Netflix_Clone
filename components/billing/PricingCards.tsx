import { ButtonLink } from "@/components/ui/Button";

const plans = [
  { name: "Free", price: "0 €", features: ["Trailer und Demo-Inhalte", "1 Profil", "Kein Vollzugriff"], href: "/register" },
  { name: "Basic", price: "4,99 €", features: ["Alle Inhalte", "1 Profil", "Monatlich kündbar"], href: "/billing?plan=BASIC" },
  { name: "Premium", price: "9,99 €", features: ["Alle Inhalte", "Mehrere Profile", "Bessere Qualität vorbereitet"], href: "/billing?plan=PREMIUM" }
];
export function PricingCards() { return <div className="grid gap-6 md:grid-cols-3">{plans.map((p) => <div key={p.name} className="rounded-2xl border border-white/10 bg-zinc-950 p-6 shadow-2xl"><h3 className="text-2xl font-black">{p.name}</h3><p className="mt-3 text-4xl font-black">{p.price}<span className="text-sm text-zinc-400"> / Monat</span></p><ul className="my-6 space-y-2 text-zinc-300">{p.features.map((f) => <li key={f}>✓ {f}</li>)}</ul><ButtonLink href={p.href} className="w-full">Auswählen</ButtonLink></div>)}</div>; }
