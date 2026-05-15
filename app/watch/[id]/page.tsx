import { notFound, redirect } from "next/navigation";
import { VideoPlayer } from "@/components/player/VideoPlayer";
import { getSession } from "@/lib/auth";
import { canStreamFullContent } from "@/lib/permissions";
import { prisma } from "@/lib/prisma";
export default async function WatchMovie({params}:{params:Promise<{id:string}>}){const session=await getSession();if(!session)redirect("/login");if(!canStreamFullContent(session))redirect("/pricing");const {id}=await params;const media=await prisma.media.findUnique({where:{id}});if(!media?.videoUrl)notFound();let profile=await prisma.profile.findFirst({where:{userId:session.id}});if(!profile)profile=await prisma.profile.create({data:{userId:session.id,name:session.name??"Profil"}});const progress=await prisma.watchProgress.findFirst({where:{profileId:profile.id,mediaId:id,episodeId:null}});return <VideoPlayer src={media.videoUrl} mediaId={media.id} profileId={profile.id} initialProgress={progress?.progressSeconds??0}/>}
