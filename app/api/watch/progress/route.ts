import { NextRequest, NextResponse } from "next/server";
import { requireUser } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { progressSchema } from "@/lib/validators";

export async function POST(req: NextRequest) {
  const user = await requireUser();
  const parsed = progressSchema.parse(await req.json());
  const profile = await prisma.profile.findFirstOrThrow({
    where: { id: parsed.profileId, userId: user.id },
  });
  const progressData = {
    progressSeconds: parsed.progressSeconds,
    completed: parsed.completed,
  };

  if (parsed.episodeId) {
    const progress = await prisma.watchProgress.upsert({
      where: {
        profileId_mediaId_episodeId: {
          profileId: profile.id,
          mediaId: parsed.mediaId,
          episodeId: parsed.episodeId,
        },
      },
      create: {
        ...progressData,
        userId: user.id,
        profileId: profile.id,
        mediaId: parsed.mediaId,
        episodeId: parsed.episodeId,
      },
      update: progressData,
    });

    return NextResponse.json(progress);
  }

  const existingProgress = await prisma.watchProgress.findFirst({
    where: { profileId: profile.id, mediaId: parsed.mediaId, episodeId: null },
  });
  const progress = existingProgress
    ? await prisma.watchProgress.update({
        where: { id: existingProgress.id },
        data: progressData,
      })
    : await prisma.watchProgress.create({
        data: {
          ...progressData,
          userId: user.id,
          profileId: profile.id,
          mediaId: parsed.mediaId,
          episodeId: null,
        },
      });

  return NextResponse.json(progress);
}
