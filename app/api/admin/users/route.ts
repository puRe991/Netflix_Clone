import { NextResponse } from "next/server";
import { requireAdmin } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function GET(){await requireAdmin();return NextResponse.json(await prisma.user.findMany({select:{id:true,email:true,name:true,role:true,subscriptionStatus:true,createdAt:true}}))}
