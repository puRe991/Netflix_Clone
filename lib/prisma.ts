import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaClient, type PrismaClient as PrismaClientType } from "./generated/prisma/client";

type PrismaLike = PrismaClientType;

function createBuildMock(): PrismaLike {
  const model = new Proxy({}, {
    get: (_target, prop) => {
      if (prop === "count") return async () => 0;
      if (prop === "aggregate") return async () => ({ _sum: { progressSeconds: 0 } });
      if (prop === "findMany") return async () => [];
      if (prop === "findFirst" || prop === "findUnique" || prop === "findFirstOrThrow" || prop === "findUniqueOrThrow") return async () => null;
      if (["create", "update", "delete", "deleteMany", "upsert"].includes(String(prop))) return async () => ({});
      return async () => ({});
    }
  });
  return new Proxy({}, { get: (_target, prop) => prop === "$disconnect" ? async () => undefined : model }) as PrismaLike;
}

const skipDatabase = process.env.STREAMFLIX_SKIP_DB === "1";
const globalForPrisma = globalThis as unknown as { prisma?: PrismaLike };

function makeClient(): PrismaLike {
  if (skipDatabase) return createBuildMock();
  const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL });
  return new PrismaClient({ adapter, log: ["error", "warn"] });
}

export const prisma = globalForPrisma.prisma ?? makeClient();

if (process.env.NODE_ENV !== "production") globalForPrisma.prisma = prisma;
