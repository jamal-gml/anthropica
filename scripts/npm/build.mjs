#!/usr/bin/env node
// Assemble the npm packages for the anthropica CLI into npm/dist/.
//
//   node scripts/npm/build.mjs --version 1.2.0 --binaries cli/build [--platforms win32-x64,linux-x64]
//
// --binaries points at a folder containing anthropica_<goos>_<goarch>/anthropica[.exe],
// which is the layout the release workflow builds and the release archives unpack to.
// Output: npm/dist/anthropica (launcher) and npm/dist/anthropica-<os>-<cpu> (one binary each).
// Publish the platform packages before the launcher.
import { chmodSync, copyFileSync, existsSync, mkdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { parseArgs } from "node:util";

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "..", "..");
const OUT = join(ROOT, "npm", "dist");

export const PLATFORMS = [
  { os: "win32", cpu: "x64", goos: "windows", goarch: "amd64" },
  { os: "win32", cpu: "arm64", goos: "windows", goarch: "arm64" },
  { os: "darwin", cpu: "x64", goos: "darwin", goarch: "amd64" },
  { os: "darwin", cpu: "arm64", goos: "darwin", goarch: "arm64" },
  { os: "linux", cpu: "x64", goos: "linux", goarch: "amd64" },
  { os: "linux", cpu: "arm64", goos: "linux", goarch: "arm64" },
];

const COMMON = {
  license: "MIT",
  author: "Jamal-GML (https://github.com/jamal-gml)",
  homepage: "https://github.com/jamal-gml/anthropica#readme",
  repository: { type: "git", url: "git+https://github.com/jamal-gml/anthropica.git", directory: "npm/anthropica" },
  bugs: { url: "https://github.com/jamal-gml/anthropica/issues" },
};

function writeJSON(path, data) {
  writeFileSync(path, JSON.stringify(data, null, 2) + "\n");
}

function main() {
  const { values } = parseArgs({
    options: {
      version: { type: "string" },
      binaries: { type: "string" },
      platforms: { type: "string" },
    },
  });
  const version = values.version?.replace(/^v/, "");
  if (!version || !/^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$/.test(version)) {
    throw new Error("--version must be semver, e.g. 1.2.0");
  }
  if (!values.binaries) throw new Error("--binaries is required");
  const binaries = resolve(values.binaries);

  const wanted = values.platforms ? values.platforms.split(",").map((s) => s.trim()) : null;
  const platforms = PLATFORMS.filter((p) => !wanted || wanted.includes(`${p.os}-${p.cpu}`));
  if (wanted && platforms.length !== wanted.length) {
    throw new Error(`unknown platform in --platforms; valid: ${PLATFORMS.map((p) => `${p.os}-${p.cpu}`).join(", ")}`);
  }

  rmSync(OUT, { recursive: true, force: true });
  const license = join(ROOT, "LICENSE");

  for (const p of platforms) {
    const name = `anthropica-${p.os}-${p.cpu}`;
    const exe = p.os === "win32" ? "anthropica.exe" : "anthropica";
    const src = join(binaries, `anthropica_${p.goos}_${p.goarch}`, exe);
    if (!existsSync(src)) throw new Error(`missing binary: ${src}`);

    const dir = join(OUT, name);
    mkdirSync(join(dir, "bin"), { recursive: true });
    copyFileSync(src, join(dir, "bin", exe));
    if (p.os !== "win32") chmodSync(join(dir, "bin", exe), 0o755);
    copyFileSync(license, join(dir, "LICENSE"));
    writeFileSync(
      join(dir, "README.md"),
      `# ${name}\n\nThe ${p.os}-${p.cpu} binary for the [anthropica](https://www.npmjs.com/package/anthropica) CLI. ` +
        "Install `anthropica` instead; npm picks this package automatically.\n"
    );
    writeJSON(join(dir, "package.json"), {
      name,
      version,
      description: `The ${p.os}-${p.cpu} binary for the anthropica CLI.`,
      ...COMMON,
      os: [p.os],
      cpu: [p.cpu],
      files: ["bin/", "README.md", "LICENSE"],
      preferUnplugged: true,
    });
    console.log(`built ${name}@${version}`);
  }

  const main = join(OUT, "anthropica");
  mkdirSync(join(main, "bin"), { recursive: true });
  copyFileSync(join(ROOT, "npm", "anthropica", "bin", "anthropica.js"), join(main, "bin", "anthropica.js"));
  chmodSync(join(main, "bin", "anthropica.js"), 0o755);
  copyFileSync(join(ROOT, "npm", "anthropica", "README.md"), join(main, "README.md"));
  copyFileSync(license, join(main, "LICENSE"));
  writeJSON(join(main, "package.json"), {
    name: "anthropica",
    version,
    description:
      "Provider-neutral CLI for preparing projects for AI coding assistants and working with Anthropica role-based skill packs.",
    ...COMMON,
    keywords: ["cli", "agents-md", "ai-agents", "claude-code", "agent-skills", "devops"],
    bin: { anthropica: "bin/anthropica.js" },
    files: ["bin/", "README.md", "LICENSE"],
    engines: { node: ">=18" },
    // Always all six, even when only some were built, so the published launcher is complete.
    optionalDependencies: Object.fromEntries(PLATFORMS.map((p) => [`anthropica-${p.os}-${p.cpu}`, version])),
  });
  console.log(`built anthropica@${version}`);
}

try {
  main();
} catch (err) {
  console.error(`error: ${err.message}`);
  process.exit(1);
}
