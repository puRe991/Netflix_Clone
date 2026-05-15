import { NextResponse } from "next/server";
import { prisma } from "@/lib/prisma";
export async function GET(_:Request,{params}:{params:Promise<{slug:string}>}){const {slug}=await params;const media=await prisma.media.findUnique({where:{slug},include:{genres:{include:{genre:true}},series:{include:{seasons:{include:{episodes:true}}}}}});return media?NextResponse.json(media):NextResponse.json({error:'Not found'},{status:404})}
