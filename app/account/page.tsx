import { redirect } from "next/navigation";
import { Navbar } from "@/components/layout/Navbar";
import { getSession } from "@/lib/auth";
export default async function Account(){const session=await getSession();if(!session)redirect('/login');return <><Navbar/><main className="mx-auto min-h-screen max-w-4xl px-4 pt-28"><h1 className="text-4xl font-black">Konto</h1><div className="mt-8 rounded-xl border border-white/10 bg-zinc-950 p-6"><p><b>E-Mail:</b> {session.email}</p><p className="mt-2"><b>Abo:</b> {session.subscriptionStatus}</p><a href="/billing" className="mt-6 inline-block rounded bg-red-600 px-5 py-3 font-bold">Billing verwalten</a></div></main></>}
