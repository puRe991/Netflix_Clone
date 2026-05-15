import { NextRequest, NextResponse } from "next/server";
import { requireAdmin } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function PATCH(req:NextRequest,{params}:{params:Promise<{id:string}>}){await requireAdmin();const {id}=await params;return NextResponse.json(await prisma.media.update({where:{id},data:await req.json()}))}
export async function POST(req:NextRequest,{params}:{params:Promise<{id:string}>}){await requireAdmin();const {id}=await params;const form=Object.fromEntries((await req.formData()).entries());await prisma.media.update({where:{id},data:{title:String(form.title),description:String(form.description),thumbnailUrl:String(form.thumbnailUrl),bannerUrl:String(form.bannerUrl)}});return NextResponse.redirect(new URL('/admin/media',req.url),303)}
export async function DELETE(_:NextRequest,{params}:{params:Promise<{id:string}>}){await requireAdmin();const {id}=await params;await prisma.media.delete({where:{id}});return NextResponse.json({ok:true})}
