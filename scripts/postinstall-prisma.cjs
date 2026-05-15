const { spawnSync } = require("node:child_process");
const { existsSync } = require("node:fs");
const { join } = require("node:path");

const localPrisma = join(
  __dirname,
  "..",
  "node_modules",
  ".bin",
  process.platform === "win32" ? "prisma.cmd" : "prisma",
);

const command = existsSync(localPrisma) ? localPrisma : "prisma";
const result = spawnSync(command, ["generate"], {
  shell: process.platform === "win32",
  stdio: "inherit",
});

if (result.error) {
  console.warn(`Prisma Client generation was skipped: ${result.error.message}`);
  process.exit(0);
}

if (result.status !== 0) {
  console.warn("Prisma Client generation failed during postinstall and was skipped.");
  console.warn("Run `npm run prisma:generate` after fixing the Prisma environment.");
}

process.exit(0);
