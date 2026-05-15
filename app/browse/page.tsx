import { Navbar } from "@/components/layout/Navbar";
import { MediaRow } from "@/components/media/MediaRow";
import { SearchBar } from "@/components/media/SearchBar";
import { prisma } from "@/lib/prisma";

export const dynamic = "force-dynamic";

export default async function Browse(){const media=await prisma.media.findMany({where:{isPublished:true},orderBy:{createdAt:"desc"}});const hero=media[0];return <><Navbar/><main className="min-h-screen bg-black pb-16">{hero&&<section className="relative min-h-[72vh] bg-cover bg-center" style={{backgroundImage:`url(${hero.bannerUrl})`}}><div className="absolute inset-0 hero-gradient"/><div className="relative mx-auto max-w-7xl px-4 pt-40"><p className="font-bold text-red-500">Neu im Katalog</p><h1 className="mt-2 max-w-2xl text-5xl font-black">{hero.title}</h1><p className="mt-4 max-w-xl text-zinc-200">{hero.description}</p><a href={hero.type==="SERIES"?`/series/${hero.slug}`:`/movie/${hero.slug}`} className="mt-6 inline-block rounded bg-red-600 px-6 py-3 font-bold">Details ansehen</a></div></section>}<div className="mx-auto mt-8 max-w-7xl space-y-10 px-4"><SearchBar media={media}/><MediaRow title="Weiterschauen" items={media.slice(0,4)}/><MediaRow title="Beliebt" items={media}/><MediaRow title="Filme" items={media.filter(m=>m.type==="MOVIE")}/><MediaRow title="Serien" items={media.filter(m=>m.type==="SERIES")}/><MediaRow title="Für dich empfohlen" items={media.filter(m=>m.tags.includes("featured"))}/></div></main></>}
