#!/usr/bin/env node
// Launcher for the anthropica CLI. npm installs exactly one platform package
// (an optionalDependency filtered by os/cpu); this script finds its native
// binary and runs it with the same arguments, stdio and exit code.
"use strict";

const { spawnSync } = require("child_process");
const fs = require("fs");
const path = require("path");

const PLATFORMS = {
  "win32-x64": "anthropica-win32-x64",
  "win32-arm64": "anthropica-win32-arm64",
  "darwin-x64": "anthropica-darwin-x64",
  "darwin-arm64": "anthropica-darwin-arm64",
  "linux-x64": "anthropica-linux-x64",
  "linux-arm64": "anthropica-linux-arm64",
};

function fail(message) {
  process.stderr.write(`anthropica: ${message}\n`);
  process.exit(1);
}

function binaryPath() {
  // Escape hatch for testing or unusual platforms.
  if (process.env.ANTHROPICA_BINARY) {
    return process.env.ANTHROPICA_BINARY;
  }

  const key = `${process.platform}-${process.arch}`;
  const pkg = PLATFORMS[key];
  if (!pkg) {
    fail(
      `no prebuilt binary for ${key}. Supported: ${Object.keys(PLATFORMS).join(", ")}.\n` +
        "Build from source instead: https://github.com/jamal-gml/anthropica/tree/main/cli"
    );
  }

  let dir;
  try {
    dir = path.dirname(require.resolve(`${pkg}/package.json`));
  } catch {
    fail(
      `the platform package ${pkg} is not installed.\n` +
        "It is an optional dependency; reinstall without --omit=optional / --no-optional:\n" +
        "  npm install -g anthropica@latest"
    );
  }

  const exe = path.join(dir, "bin", process.platform === "win32" ? "anthropica.exe" : "anthropica");
  if (!fs.existsSync(exe)) {
    fail(`binary missing from ${pkg}: ${exe}. Reinstall with: npm install -g anthropica@latest`);
  }

  // Packages published from Windows lose the executable bit; restore it.
  if (process.platform !== "win32") {
    try {
      fs.accessSync(exe, fs.constants.X_OK);
    } catch {
      try {
        fs.chmodSync(exe, 0o755);
      } catch (err) {
        fail(`${exe} is not executable and could not be fixed (${err.code}). Run: chmod +x "${exe}"`);
      }
    }
  }
  return exe;
}

const result = spawnSync(binaryPath(), process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  fail(`could not start the binary: ${result.error.message}`);
}
if (result.signal) {
  // Re-raise the signal so the parent shell sees the same termination.
  process.kill(process.pid, result.signal);
}
process.exit(result.status === null ? 1 : result.status);
