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

if (isLinux32Bit) {
  hasError = true;
  add(
    "ERROR",
    "32-bit Linux is not a reliable target for this Next.js/Prisma stack because current Node.js and native package binaries are not consistently published for linux-ia32.",
  );
  add(
    "INFO",
    "Use a 64-bit Linux host/VM/container for production, or use 32-bit Windows only for local evaluation.",
  );
}

if (isWindows32Bit) {
  add(
    "INFO",
    "32-bit Windows can be used for local evaluation with the 32-bit Node.js installer. Keep Prisma engine type set to binary.",
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
