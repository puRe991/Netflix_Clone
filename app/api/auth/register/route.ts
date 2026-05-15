import { NextRequest, NextResponse } from "next/server";
import { createSession, hashPassword } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { registerSchema } from "@/lib/validators";
export async function POST(req: NextRequest){const data=Object.fromEntries((await req.formData()).entries());const parsed=registerSchema.parse(data);const user=await prisma.user.create({data:{email:parsed.email.toLowerCase(),name:parsed.name,passwordHash:await hashPassword(parsed.password),profiles:{create:{name:parsed.name??"Profil"}}}});await createSession(user);return NextResponse.redirect(new URL('/profiles',req.url),303)}
