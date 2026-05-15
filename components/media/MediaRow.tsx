import { MovieCard } from "./MovieCard";

type RowMedia = { id: string; title: string; slug: string; thumbnailUrl: string; type: "MOVIE" | "SERIES"; releaseYear: number; ageRating: string };
export function MediaRow({ title, items }: { title: string; items: RowMedia[] }) {
  if (!items.length) return null;
  return <section className="space-y-4"><h2 className="text-xl font-extrabold md:text-2xl">{title}</h2><div className="scrollbar-hide flex gap-4 overflow-x-auto pb-4">{items.map((m) => <MovieCard key={m.id} media={m} />)}</div></section>;
}
