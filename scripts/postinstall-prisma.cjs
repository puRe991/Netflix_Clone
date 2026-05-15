const { runPrismaGenerate } = require("./prisma-generate.cjs");

process.exit(runPrismaGenerate({ optional: true }));
