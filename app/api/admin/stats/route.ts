import { NextResponse } from "next/server";
import { requireAdmin } from "@/lib/auth";
import { prisma } from "@/lib/prisma";
export async function GET(){await requireAdmin();const [users,active,media,watchtime]=await Promise.all([prisma.user.count(),prisma.user.count({where:{subscriptionStatus:'ACTIVE'}}),prisma.media.count(),prisma.watchProgress.aggregate({_sum:{progressSeconds:true}})]);return NextResponse.json({users,activeSubscriptions:active,media,watchtimeSeconds:watchtime._sum.progressSeconds??0})}
