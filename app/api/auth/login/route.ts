import { NextRequest, NextResponse } from "next/server";
import { createSession, verifyPassword } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { loginSchema } from "@/lib/validators";
export async function POST(req: NextRequest){const data=Object.fromEntries((await req.formData()).entries());const parsed=loginSchema.parse(data);const user=await prisma.user.findUnique({where:{email:parsed.email.toLowerCase()}});if(!user||!(await verifyPassword(parsed.password,user.passwordHash)))return NextResponse.json({error:'Invalid credentials'},{status:401});await createSession(user);return NextResponse.redirect(new URL('/browse',req.url),303)}
