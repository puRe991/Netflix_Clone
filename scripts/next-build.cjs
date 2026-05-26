const { spawnSync } = require("node:child_process");
const { existsSync } = require("node:fs");
const { join } = require("node:path");

const localNext = join(
  __dirname,
  "..",
  "node_modules",
  ".bin",
  process.platform === "win32" ? "next.cmd" : "next",
);
const command = existsSync(localNext) ? localNext : "next";

const result = spawnSync(command, ["build"], {
  env: { ...process.env, STREAMFLIX_SKIP_DB: "1" },
  shell: process.platform === "win32",
  stdio: "inherit",
});

if (result.error) {
  console.error(`Next.js build failed: ${result.error.message}`);
  process.exit(1);
}

process.exit(result.status || 0);
