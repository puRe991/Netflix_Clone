import { NextResponse } from "next/server";
import { requireUser } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function GET(_:Request,{params}:{params:Promise<{id:string}>}){const user=await requireUser();const {id}=await params;return NextResponse.json(await prisma.watchProgress.findMany({where:{userId:user.id,mediaId:id}}))}
