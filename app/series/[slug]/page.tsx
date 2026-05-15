import { notFound } from "next/navigation";
import { Navbar } from "@/components/layout/Navbar";
import { prisma } from "@/lib/prisma";

export default async function SeriesDetail({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const media = await prisma.media.findUnique({
    where: { slug },
    include: { series: { include: { seasons: { include: { episodes: true }, orderBy: { seasonNumber: "asc" } } } } }
  });
  if (!media || !media.series) notFound();
  return (
    <>
      <Navbar />
      <main className="mx-auto min-h-screen max-w-7xl px-4 pt-28">
        <h1 className="text-5xl font-black">{media.title}</h1>
        <p className="mt-4 max-w-3xl text-zinc-300">{media.description}</p>
        <div className="mt-8 space-y-8">
          {media.series.seasons.map((s) => (
            <section key={s.id}>
              <h2 className="text-2xl font-bold">Staffel {s.seasonNumber}: {s.title}</h2>
              <div className="mt-4 grid gap-4 md:grid-cols-2">
                {s.episodes.map((e) => (
                  <a href={`/watch/episode/${e.id}`} key={e.id} className="rounded-lg border border-white/10 bg-zinc-950 p-4 hover:border-red-600">
                    <h3 className="font-bold">{e.episodeNumber}. {e.title}</h3>
                    <p className="mt-2 text-sm text-zinc-400">{e.description}</p>
                  </a>
                ))}
              </div>
            </section>
          ))}
        </div>
      </main>
    </>
  );
}
