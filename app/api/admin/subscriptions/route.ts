import { NextResponse } from "next/server";
import { requireAdmin } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function GET(){await requireAdmin();return NextResponse.json(await prisma.subscription.findMany({include:{user:{select:{email:true,name:true}}},orderBy:{createdAt:'desc'}}))}
