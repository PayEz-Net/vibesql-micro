#!/usr/bin/env node
"use strict";

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");

const ext = process.platform === "win32" ? ".exe" : "";
const binary = path.join(__dirname, `vibesql-micro${ext}`);

if (!fs.existsSync(binary)) {
  console.error("vibesql-micro: binary not found. Run `npm install` or download from:");
  console.error("https://github.com/PayEz-Net/vibesql-micro/releases");
  process.exit(1);
}

const child = spawn(binary, process.argv.slice(2), { stdio: "inherit" });

function forward(signal) {
  if (child.pid) {
    try { process.kill(child.pid, signal); } catch {}
  }
}

process.on("SIGINT", () => forward("SIGINT"));
process.on("SIGTERM", () => forward("SIGTERM"));

child.on("exit", (code, signal) => {
  process.exit(code ?? (signal ? 1 : 0));
});
