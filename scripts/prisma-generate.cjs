const { spawnSync } = require("node:child_process");
const { existsSync } = require("node:fs");
const { join } = require("node:path");

function runPrismaGenerate({ optional = false } = {}) {
  const localPrisma = join(
    __dirname,
    "..",
    "node_modules",
    ".bin",
    process.platform === "win32" ? "prisma.cmd" : "prisma",
  );

  const command = existsSync(localPrisma) ? localPrisma : "prisma";

  const result = spawnSync(command, ["generate"], {
    env: process.env,
    shell: process.platform === "win32",
    stdio: "inherit",
  });

  if (result.error) {
    const message = `Prisma Client generation failed: ${result.error.message}`;
    if (optional) {
      console.warn(`${message}. Skipping optional postinstall generation.`);
      return 0;
    }

    console.error(message);
    return 1;
  }

  if (result.status !== 0) {
    if (optional) {
      console.warn("Prisma Client generation failed during postinstall and was skipped.");
      console.warn("Run `npm run prisma:generate` after fixing the Prisma environment.");
      return 0;
    }

    return result.status || 1;
  }

  return 0;
}

if (require.main === module) {
  process.exit(runPrismaGenerate());
}

module.exports = { runPrismaGenerate };
