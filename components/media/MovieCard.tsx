import Link from "next/link";

type MediaCardProps = { title: string; slug: string; thumbnailUrl: string; type: "MOVIE" | "SERIES"; releaseYear: number; ageRating: string };

export function MovieCard({ media }: { media: MediaCardProps }) {
  const href = media.type === "SERIES" ? `/series/${media.slug}` : `/movie/${media.slug}`;
  return <Link href={href} className="card-hover group min-w-[170px] overflow-hidden rounded-lg border border-white/10 bg-zinc-950 sm:min-w-[220px]">
    <div className="aspect-[2/3] bg-cover bg-center" style={{ backgroundImage: `url(${media.thumbnailUrl})` }} />
    <div className="p-3"><h3 className="line-clamp-1 font-bold">{media.title}</h3><p className="mt-1 text-xs text-zinc-400">{media.releaseYear} • {media.ageRating}</p></div>
  </Link>;
}
