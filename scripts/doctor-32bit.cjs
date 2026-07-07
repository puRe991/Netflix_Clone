const { arch, platform, versions } = require("node:process");

const MIN_NODE_MAJOR = 20;
const supportedArchitectures = new Set(["x64", "arm64", "ia32"]);
const nodeMajor = Number.parseInt(versions.node.split(".")[0] || "0", 10);
const is32Bit = arch === "ia32";
const isLinux32Bit = platform === "linux" && is32Bit;
const isWindows32Bit = platform === "win32" && is32Bit;

const messages = [];
let hasError = false;

function add(status, message) {
  messages.push(`${status} ${message}`);
}

if (!supportedArchitectures.has(arch)) {
  hasError = true;
  add("ERROR", `Unsupported CPU architecture: ${arch}. Use x64, arm64, or ia32.`);
} else {
  add("OK", `CPU architecture detected: ${arch}.`);
}

if (Number.isNaN(nodeMajor) || nodeMajor < MIN_NODE_MAJOR) {
  hasError = true;
  add("ERROR", `Node.js ${MIN_NODE_MAJOR} or newer is required. Detected: ${versions.node}.`);
} else {
  add("OK", `Node.js version detected: ${versions.node}.`);
}

if (nodeMajor >= 23 && is32Bit) {
  hasError = true;
  add(
    "ERROR",
    `Node.js ${versions.node} has no 32-bit Windows build. Node.js dropped ia32 Windows support starting with Node.js 23; install Node.js 20.19+ or 22.12+ instead.`,
  );
}

if (isLinux32Bit) {
  hasError = true;
  add(
    "ERROR",
    "32-bit Linux (linux-ia32) is not a published Node.js target and is not supported by this stack.",
  );
  add(
    "INFO",
    "Use a 64-bit Linux host/VM/container for production, or use 32-bit Windows only for local evaluation.",
  );
}

if (isWindows32Bit) {
  add(
    "OK",
    "32-bit Windows can run the app itself. Prisma ORM 7 uses a WASM query compiler with a JS driver adapter (@prisma/adapter-pg) at runtime instead of a native, architecture-specific engine binary, so `npm run dev`/`npm start` do not need any ia32-specific workaround.",
  );
  add(
    "WARN",
    "Prisma's schema engine (used by `prisma migrate`, `db push`, `db pull`, `studio`) only ships a 64-bit Windows binary and will likely fail to start on 32-bit Windows. Run those commands from a 64-bit machine (or WSL) against the same database, then only run the app itself on the 32-bit machine.",
  );
}

if (is32Bit) {
  add(
    "INFO",
    "Recommended low-resource settings: NODE_OPTIONS=--max-old-space-size=1024 and a Postgres database on another machine or managed service.",
  );
}

console.log(messages.join("\n"));
process.exit(hasError ? 1 : 0);
