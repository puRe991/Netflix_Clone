import { NextRequest, NextResponse } from "next/server";
import { requireUser } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
import { stripe } from "@/lib/stripe";
export async function POST(req:NextRequest){const user=await requireUser();const dbUser=await prisma.user.findUniqueOrThrow({where:{id:user.id}});if(!dbUser.stripeCustomerId)return NextResponse.redirect(new URL('/pricing',req.url),303);const portal=await stripe.billingPortal.sessions.create({customer:dbUser.stripeCustomerId,return_url:`${process.env.NEXT_PUBLIC_APP_URL}/account`});return NextResponse.redirect(portal.url,303)}
