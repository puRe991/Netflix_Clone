import { redirect } from "next/navigation";
import { Navbar } from "@/components/layout/Navbar";
import { MovieCard } from "@/components/media/MovieCard";
import { getSession } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export default async function MyList(){const session=await getSession();if(!session)redirect('/login');const items=await prisma.watchlist.findMany({where:{userId:session.id},include:{media:true},orderBy:{createdAt:'desc'}});return <><Navbar/><main className="mx-auto min-h-screen max-w-7xl px-4 pt-28"><h1 className="text-4xl font-black">Meine Liste</h1><div className="mt-8 flex flex-wrap gap-4">{items.map(i=><MovieCard key={i.id} media={i.media}/>)}</div>{!items.length&&<p className="mt-8 text-zinc-400">Noch keine Titel gespeichert.</p>}</main></>}
