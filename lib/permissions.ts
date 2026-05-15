import type { SessionUser } from "@/lib/auth";

export function canStreamFullContent(user: SessionUser | null) {
  return user?.subscriptionStatus === "ACTIVE" || user?.subscriptionStatus === "TRIALING" || user?.role === "ADMIN";
}

export function isAdmin(user: SessionUser | null) { return user?.role === "ADMIN"; }
