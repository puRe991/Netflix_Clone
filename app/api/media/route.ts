import { NextResponse } from "next/server";
import { prisma } from "@/lib/prisma";
export async function GET(){return NextResponse.json(await prisma.media.findMany({where:{isPublished:true},include:{genres:{include:{genre:true}}}}))}
