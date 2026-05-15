import { NextRequest, NextResponse } from "next/server";
import { requireAdmin } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function POST(req:NextRequest){await requireAdmin();const body=await req.json();return NextResponse.json(await prisma.genre.create({data:{name:body.name,slug:body.slug}}))}
