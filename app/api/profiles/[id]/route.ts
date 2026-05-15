import { NextRequest, NextResponse } from "next/server";
import { requireUser } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function PATCH(req:NextRequest,{params}:{params:Promise<{id:string}>}){const user=await requireUser();const {id}=await params;const body=await req.json();return NextResponse.json(await prisma.profile.update({where:{id,userId:user.id},data:body}))}
export async function DELETE(_:NextRequest,{params}:{params:Promise<{id:string}>}){const user=await requireUser();const {id}=await params;await prisma.profile.delete({where:{id,userId:user.id}});return NextResponse.json({ok:true})}
