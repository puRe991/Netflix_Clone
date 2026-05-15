export type MediaKind = "MOVIE" | "SERIES";
export type MediaCard = { id: string; title: string; slug: string; thumbnailUrl: string; type: MediaKind; releaseYear: number; ageRating: string };
