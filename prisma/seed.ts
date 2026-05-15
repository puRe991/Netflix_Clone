import { PrismaClient } from "@prisma/client";
import bcrypt from "bcryptjs";
const prisma = new PrismaClient();
const sampleVideo = "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4";
async function main(){
  const [drama,sciFi,doc] = await Promise.all([
    prisma.genre.upsert({where:{slug:"drama"},update:{},create:{name:"Drama",slug:"drama"}}),
    prisma.genre.upsert({where:{slug:"sci-fi"},update:{},create:{name:"Sci-Fi",slug:"sci-fi"}}),
    prisma.genre.upsert({where:{slug:"dokumentation"},update:{},create:{name:"Dokumentation",slug:"dokumentation"}})
  ]);
  await prisma.user.upsert({where:{email:"admin@streamflix.local"},update:{},create:{email:"admin@streamflix.local",name:"Admin",role:"ADMIN",subscriptionStatus:"ACTIVE",passwordHash:await bcrypt.hash("StreamFlix123!",12),profiles:{create:{name:"Admin"}}}});
  await prisma.media.upsert({where:{slug:"ocean-city"},update:{},create:{title:"Ocean City",slug:"ocean-city",description:"Ein lizenzierter Demo-Film über Mut, Gemeinschaft und den Neubeginn an der Küste.",type:"MOVIE",thumbnailUrl:"https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=600&q=80",bannerUrl:"https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1800&q=80",videoUrl:sampleVideo,trailerUrl:sampleVideo,releaseYear:2026,duration:596,ageRating:"FSK 6",language:"Deutsch",tags:["featured","coast","family"],isPublished:true,genres:{create:[{genreId:drama.id},{genreId:doc.id}]}}});
  const seriesMedia=await prisma.media.upsert({where:{slug:"red-orbit"},update:{},create:{title:"Red Orbit",slug:"red-orbit",description:"Eine legale Sci-Fi-Serie über eine Forschungscrew auf dem Weg zu einem roten Planeten.",type:"SERIES",thumbnailUrl:"https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?auto=format&fit=crop&w=600&q=80",bannerUrl:"https://images.unsplash.com/photo-1462331940025-496dfbfc7564?auto=format&fit=crop&w=1800&q=80",trailerUrl:sampleVideo,releaseYear:2026,ageRating:"FSK 12",language:"Deutsch",tags:["featured","space"],isPublished:true,genres:{create:[{genreId:sciFi.id}]}}});
  const series=await prisma.series.upsert({where:{mediaId:seriesMedia.id},update:{},create:{mediaId:seriesMedia.id,title:"Red Orbit",description:"Staffel 1 der Forschungsmission."}});
  const season=await prisma.season.upsert({where:{seriesId_seasonNumber:{seriesId:series.id,seasonNumber:1}},update:{},create:{seriesId:series.id,seasonNumber:1,title:"Mission Start"}});
  for(const n of [1,2,3]) await prisma.episode.upsert({where:{seasonId_episodeNumber:{seasonId:season.id,episodeNumber:n}},update:{},create:{seasonId:season.id,episodeNumber:n,title:`Episode ${n}`,description:`Kapitel ${n} der Red-Orbit-Mission.`,duration:596,videoUrl:sampleVideo,thumbnailUrl:"https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?auto=format&fit=crop&w=800&q=80",isPublished:true}});
}
main().finally(()=>prisma.$disconnect());
