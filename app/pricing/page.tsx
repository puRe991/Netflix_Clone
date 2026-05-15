import { Navbar } from "@/components/layout/Navbar";
import { PricingCards } from "@/components/billing/PricingCards";
export default function PricingPage(){return <><Navbar/><main className="mx-auto min-h-screen max-w-7xl px-4 pt-28"><h1 className="text-4xl font-black">Abo-Modelle</h1><p className="mt-3 text-zinc-400">Free für Trailer, Basic und Premium für Vollzugriff über Stripe Billing.</p><div className="mt-8"><PricingCards/></div></main></>}
