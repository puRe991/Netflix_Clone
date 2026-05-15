import { NextResponse } from "next/server";
import { requireUser, hasActiveSubscription } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function GET(_:Request,{params}:{params:Promise<{id:string}>}){const user=await requireUser();if(!(await hasActiveSubscription(user.id))&&user.role!=='ADMIN')return NextResponse.json({error:'Subscription required'},{status:402});const {id}=await params;const media=await prisma.media.findUnique({where:{id},select:{id:true,title:true,videoUrl:true,trailerUrl:true}});return media?NextResponse.json(media):NextResponse.json({error:'Not found'},{status:404})}
