import { NextRequest, NextResponse } from "next/server";
import { requireUser } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { progressSchema } from "@/lib/validators";
export async function POST(req:NextRequest){const user=await requireUser();const parsed=progressSchema.parse(await req.json());const profile=await prisma.profile.findFirstOrThrow({where:{id:parsed.profileId,userId:user.id}});const progress=await prisma.watchProgress.upsert({where:{profileId_mediaId_episodeId:{profileId:profile.id,mediaId:parsed.mediaId,episodeId:parsed.episodeId??null}},create:{...parsed,userId:user.id,profileId:profile.id,episodeId:parsed.episodeId??null},update:{progressSeconds:parsed.progressSeconds,completed:parsed.completed}});return NextResponse.json(progress)}
