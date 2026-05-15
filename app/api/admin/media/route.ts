import { NextRequest, NextResponse } from "next/server";
import { requireAdmin } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { mediaSchema } from "@/lib/validators";
export async function POST(req:NextRequest){await requireAdmin();const form=Object.fromEntries((await req.formData()).entries());const parsed=mediaSchema.parse({...form,tags:[],genreIds:[],isPublished:false});await prisma.media.create({data:{...parsed,trailerUrl:parsed.trailerUrl||null,videoUrl:parsed.videoUrl||null,duration:parsed.duration||null,genres:{create:parsed.genreIds.map(genreId=>({genreId}))}}});return NextResponse.redirect(new URL('/admin/media',req.url),303)}
