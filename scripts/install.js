#!/usr/bin/env node
"use strict";

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");

const VERSION = require("../package.json").version;
const REPO = "PayEz-Net/vibesql-micro";
const BIN_DIR = path.join(__dirname, "..", "bin");

const PLATFORM_MAP = {
  "win32-x64": "vibesql-micro-windows-x64.exe",
  "linux-x64": "vibesql-micro-linux-x64",
  "darwin-x64": "vibesql-micro-macos-amd64",
  "darwin-arm64": "vibesql-micro-macos-arm64",
};

const platform = process.platform;
const arch = process.arch;
const key = `${platform}-${arch}`;
const binaryName = PLATFORM_MAP[key];

if (!binaryName) {
  console.error(`vibesql-micro: unsupported platform ${key}`);
  console.error(`Supported: ${Object.keys(PLATFORM_MAP).join(", ")}`);
  process.exit(1);
}

const ext = platform === "win32" ? ".exe" : "";
const dest = path.join(BIN_DIR, `vibesql-micro${ext}`);
const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;

// Check if correct version already installed
if (fs.existsSync(dest)) {
  try {
    const installed = execSync(`"${dest}" version`, { encoding: "utf8", timeout: 5000 });
    const match = installed.match(/Version:\s+(\S+)/);
    if (match && match[1] === VERSION) {
      console.log(`vibesql-micro: v${VERSION} already installed, skipping download`);
      process.exit(0);
    }
    console.log(`vibesql-micro: installed ${match ? match[1] : "unknown"} != ${VERSION}, updating...`);
    fs.unlinkSync(dest);
  } catch {
    console.log("vibesql-micro: existing binary corrupt or unreadable, replacing...");
    try { fs.unlinkSync(dest); } catch {}
  }
}

fs.mkdirSync(BIN_DIR, { recursive: true });

console.log(`vibesql-micro: downloading v${VERSION} for ${key}...`);

function download(url, dest, redirects) {
  if (redirects > 5) {
    console.error("vibesql-micro: too many redirects");
    process.exit(1);
  }

  const proto = url.startsWith("https") ? https : require("http");
  proto.get(url, { headers: { "User-Agent": "vibesql-micro-npm" } }, (res) => {
    if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
      download(res.headers.location, dest, redirects + 1);
      return;
    }

    if (res.statusCode !== 200) {
      console.error(`vibesql-micro: download failed (HTTP ${res.statusCode})`);
      console.error(`URL: ${url}`);
      console.error("");
      console.error("This usually means the release binaries haven't been uploaded yet.");
      console.error("Download manually from: https://github.com/PayEz-Net/vibesql-micro/releases");
      // Non-fatal — let the user run the CLI and get a clear error
      process.exit(0);
    }

    const file = fs.createWriteStream(dest);
    let downloaded = 0;
    const total = parseInt(res.headers["content-length"], 10) || 0;

    res.on("data", (chunk) => {
      downloaded += chunk.length;
      if (total > 0) {
        const pct = Math.round((downloaded / total) * 100);
        process.stdout.write(`\rvibesql-micro: ${pct}% (${Math.round(downloaded / 1024 / 1024)}MB)`);
      }
    });

    res.pipe(file);

    file.on("finish", () => {
      file.close();
      console.log("");

      if (platform !== "win32") {
        fs.chmodSync(dest, 0o755);
      }

      console.log(`vibesql-micro: v${VERSION} installed successfully`);
    });
  }).on("error", (err) => {
    console.error(`vibesql-micro: download error: ${err.message}`);
    console.error("Download manually from: https://github.com/PayEz-Net/vibesql-micro/releases");
    process.exit(0);
  });
}

download(url, dest, 0);
