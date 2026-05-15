import { NextRequest, NextResponse } from "next/server";
import { requireUser } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function GET(){const user=await requireUser();return NextResponse.json(await prisma.watchlist.findMany({where:{userId:user.id},include:{media:true}}))}
export async function POST(req:NextRequest){const user=await requireUser();let mediaId='';try{const form=await req.formData();mediaId=String(form.get('mediaId')??'')}catch{mediaId=(await req.json()).mediaId}let profile=await prisma.profile.findFirst({where:{userId:user.id}});if(!profile)profile=await prisma.profile.create({data:{userId:user.id,name:user.name??'Profil'}});await prisma.watchlist.upsert({where:{profileId_mediaId:{profileId:profile.id,mediaId}},create:{userId:user.id,profileId:profile.id,mediaId},update:{}});return NextResponse.redirect(new URL('/my-list',req.url),303)}
