import { NextRequest, NextResponse } from "next/server";
import { requireUser } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { profileSchema } from "@/lib/validators";
export async function GET(){const user=await requireUser();return NextResponse.json(await prisma.profile.findMany({where:{userId:user.id}}))}
export async function POST(req:NextRequest){const user=await requireUser();const form=Object.fromEntries((await req.formData()).entries());const parsed=profileSchema.parse({name:form.name,avatarUrl:form.avatarUrl??'',isKidsProfile:form.isKidsProfile==='on'});await prisma.profile.create({data:{...parsed,userId:user.id}});return NextResponse.redirect(new URL('/profiles',req.url),303)}
