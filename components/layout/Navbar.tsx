import Link from "next/link";
import { getSession } from "@/lib/auth";

export async function Navbar() {
  const session = await getSession();
  return (
    <header className="fixed inset-x-0 top-0 z-50 border-b border-white/10 bg-black/70 backdrop-blur">
      <nav className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3">
        <Link href={session ? "/browse" : "/"} className="text-2xl font-black tracking-tight text-red-600">StreamFlix</Link>
        <div className="flex items-center gap-4 text-sm text-zinc-200">
          {session ? <>
            <Link href="/browse" className="hover:text-white">Browse</Link><Link href="/my-list" className="hover:text-white">Meine Liste</Link><Link href="/account" className="hover:text-white">Konto</Link>
            {session.role === "ADMIN" && <Link href="/admin" className="rounded bg-red-600 px-3 py-1 font-semibold">Admin</Link>}
            <form action="/api/auth/logout" method="post"><button className="hover:text-white">Logout</button></form>
          </> : <><Link href="/pricing">Preise</Link><Link href="/login" className="rounded bg-red-600 px-4 py-2 font-bold">Einloggen</Link></>}
        </div>
      </nav>
    </header>
  );
}
