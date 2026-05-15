import { z } from "zod";

export const registerSchema = z.object({ email: z.string().email(), password: z.string().min(8), name: z.string().min(2).max(80).optional() });
export const loginSchema = z.object({ email: z.string().email(), password: z.string().min(1) });
export const profileSchema = z.object({ name: z.string().min(2).max(40), avatarUrl: z.string().url().optional().or(z.literal("")), isKidsProfile: z.boolean().default(false) });
export const mediaSchema = z.object({
  title: z.string().min(2), slug: z.string().min(2).regex(/^[a-z0-9-]+$/), description: z.string().min(10), type: z.enum(["MOVIE", "SERIES"]),
  thumbnailUrl: z.string().url(), bannerUrl: z.string().url(), trailerUrl: z.string().url().optional().or(z.literal("")), videoUrl: z.string().url().optional().or(z.literal("")),
  releaseYear: z.coerce.number().int().min(1900), duration: z.coerce.number().int().optional(), ageRating: z.string().min(1), language: z.string().default("Deutsch"),
  tags: z.array(z.string()).default([]), isPublished: z.boolean().default(false), genreIds: z.array(z.string()).default([])
});
export const progressSchema = z.object({ mediaId: z.string(), profileId: z.string(), episodeId: z.string().optional(), progressSeconds: z.coerce.number().int().min(0), completed: z.boolean().default(false) });
