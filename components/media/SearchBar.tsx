"use client";
import { useMemo, useState } from "react";
import { MovieCard } from "./MovieCard";

type SearchMedia = { id: string; title: string; slug: string; description: string; thumbnailUrl: string; type: "MOVIE" | "SERIES"; releaseYear: number; ageRating: string; tags: string[] };
export function SearchBar({ media }: { media: SearchMedia[] }) {
  const [query, setQuery] = useState("");
  const results = useMemo(() => media.filter((m) => [m.title, m.description, ...m.tags].join(" ").toLowerCase().includes(query.toLowerCase())).slice(0, 10), [media, query]);
  return <div className="space-y-4"><input value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Titel, Genre oder Tag suchen..." className="w-full rounded-lg border border-white/10 bg-white/10 px-4 py-3 outline-none focus:border-red-500" />{query && <div className="flex gap-4 overflow-x-auto pb-3">{results.map((m) => <MovieCard key={m.id} media={m} />)}{!results.length && <p className="text-zinc-400">Keine Treffer gefunden.</p>}</div>}</div>;
}
