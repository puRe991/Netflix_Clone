import { cookies } from "next/headers";
import { SignJWT, jwtVerify } from "jose";
import bcrypt from "bcryptjs";
import { prisma } from "@/lib/prisma";

const cookieName = "streamflix_session";
const developmentSecret = "dev-only-secret-change-me";

function getJwtSecret() {
  const configuredSecret = process.env.JWT_SECRET;

  const isProductionRuntime = process.env.NODE_ENV === "production" && process.env.NEXT_PHASE !== "phase-production-build";

  if (isProductionRuntime) {
    if (!configuredSecret) {
      throw new Error("JWT_SECRET must be set in production.");
    }

    if (configuredSecret.length < 32) {
      throw new Error("JWT_SECRET must be at least 32 characters in production.");
    }
  }

  return configuredSecret ?? developmentSecret;
}

const secret = new TextEncoder().encode(getJwtSecret());
export type SessionUser = { id: string; email: string; role: "USER" | "ADMIN"; subscriptionStatus: "FREE" | "TRIALING" | "ACTIVE" | "PAST_DUE" | "CANCELED" | "INCOMPLETE"; name?: string | null };

export async function hashPassword(password: string) { return bcrypt.hash(password, 12); }
export async function verifyPassword(password: string, hash: string) { return bcrypt.compare(password, hash); }

export async function createSession(user: SessionUser) {
  const token = await new SignJWT(user).setProtectedHeader({ alg: "HS256" }).setIssuedAt().setExpirationTime("7d").sign(secret);
  (await cookies()).set(cookieName, token, { httpOnly: true, sameSite: "lax", secure: process.env.NODE_ENV === "production", path: "/", maxAge: 60 * 60 * 24 * 7 });
}

export async function clearSession() { (await cookies()).delete(cookieName); }

export async function getSession(): Promise<SessionUser | null> {
  const token = (await cookies()).get(cookieName)?.value;
  if (!token) return null;
  try {
    const { payload } = await jwtVerify(token, secret);
    return payload as SessionUser;
  } catch { return null; }
}

export async function requireUser() {
  const session = await getSession();
  if (!session) throw new Error("UNAUTHENTICATED");
  return session;
}

export async function requireAdmin() {
  const session = await requireUser();
  if (session.role !== "ADMIN") throw new Error("FORBIDDEN");
  return session;
}

export async function hasActiveSubscription(userId: string) {
  const user = await prisma.user.findUnique({ where: { id: userId }, select: { subscriptionStatus: true } });
  return user?.subscriptionStatus === "ACTIVE" || user?.subscriptionStatus === "TRIALING";
}
