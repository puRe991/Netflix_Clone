import type { PrismaClient as PrismaClientType } from "@prisma/client";

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
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const { PrismaClient } = require("@prisma/client") as { PrismaClient: new (options?: unknown) => PrismaClientType };
  return new PrismaClient({ log: ["error", "warn"] });
}

export const prisma = globalForPrisma.prisma ?? makeClient();

if (process.env.NODE_ENV !== "production") globalForPrisma.prisma = prisma;
