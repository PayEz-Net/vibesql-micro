#!/usr/bin/env node
'use strict';

// Postinstall: download the native binary for the current platform from the
// matching GitHub Release and verify its SHA256.
//
// The upstream release tag this npm version wraps is the `BINARY_VERSION`
// constant below. Bump it alongside the npm version bump when you want the
// package to pull a newer native binary.

const { createHash } = require('crypto');
const { createWriteStream, existsSync, chmodSync, mkdirSync, readFileSync, statSync } = require('fs');
const { get } = require('https');
const { URL } = require('url');
const { join } = require('path');

const BINARY_VERSION = 'v0.3.1';

// Currently only Linux x64 is published. Windows/macOS builds will be added to
// both this map and the upstream release assets together.
const PUBLISHED_ASSETS = {
  'linux-x64': 'vibesql-micro-linux-x64',
};

const USER_AGENT = `vibesql-micro-npm/${require('./package.json').version}`;

function log(msg)  { console.log(`[vibesql-micro] ${msg}`); }
function warn(msg) { console.error(`[vibesql-micro] ${msg}`); }
function die(msg)  { warn(msg); process.exit(1); }

function platformTag(platform, arch) {
  if (platform === 'linux'  && arch === 'x64')   return 'linux-x64';
  if (platform === 'darwin' && arch === 'x64')   return 'darwin-x64';
  if (platform === 'darwin' && arch === 'arm64') return 'darwin-arm64';
  if (platform === 'win32'  && arch === 'x64')   return 'win32-x64';
  return null;
}

function request(url, headers, onResponse, redirects = 0) {
  if (redirects > 5) return onResponse(new Error('too many redirects'));
  const opts = {
    headers: { 'User-Agent': USER_AGENT, 'Accept': '*/*', ...headers },
  };
  get(url, opts, (res) => {
    if ([301, 302, 307, 308].includes(res.statusCode) && res.headers.location) {
      res.resume();
      const next = new URL(res.headers.location, url).toString();
      return request(next, headers, onResponse, redirects + 1);
    }
    onResponse(null, res);
  }).on('error', (err) => onResponse(err));
}

function downloadToFile(url, destPath) {
  return new Promise((resolve, reject) => {
    request(url, {}, (err, res) => {
      if (err) return reject(err);
      if (res.statusCode !== 200) {
        return reject(new Error(`HTTP ${res.statusCode} downloading ${url}`));
      }
      const out = createWriteStream(destPath);
      res.pipe(out);
      out.on('finish', () => out.close(() => resolve()));
      out.on('error', reject);
    });
  });
}

function fetchText(url) {
  return new Promise((resolve, reject) => {
    request(url, {}, (err, res) => {
      if (err) return reject(err);
      if (res.statusCode !== 200) {
        return reject(new Error(`HTTP ${res.statusCode} fetching ${url}`));
      }
      let body = '';
      res.setEncoding('utf8');
      res.on('data', (chunk) => { body += chunk; });
      res.on('end', () => resolve(body));
      res.on('error', reject);
    });
  });
}

function sha256File(path) {
  return createHash('sha256').update(readFileSync(path)).digest('hex');
}

async function main() {
  const tag = platformTag(process.platform, process.arch);
  if (!tag || !PUBLISHED_ASSETS[tag]) {
    warn(`Platform/arch ${process.platform}/${process.arch} is not yet published as a release asset.`);
    warn(`Linux x64 is the current target; Windows/macOS builds are in progress.`);
    warn(`The bin shim will report the same message if the binary is ever invoked.`);
    return; // do not fail install — let the bin shim handle invocation-time messaging
  }

  const assetName = PUBLISHED_ASSETS[tag];
  const binDir = join(__dirname, 'bin');
  mkdirSync(binDir, { recursive: true });
  const ext = process.platform === 'win32' ? '.exe' : '';
  const binPath = join(binDir, `vibesql-micro-${tag}${ext}`);

  // Already installed at the expected binary-version? Keep it.
  if (existsSync(binPath) && statSync(binPath).size > 0) {
    log(`Binary already present (${binPath}), skipping download.`);
    return;
  }

  const base = `https://github.com/PayEz-Net/vibesql-micro/releases/download/${BINARY_VERSION}`;
  log(`Downloading ${assetName} from ${BINARY_VERSION}...`);

  try {
    await downloadToFile(`${base}/${assetName}`, binPath);
  } catch (err) {
    die(`Download failed: ${err.message}. You can retry with: npm rebuild vibesql-micro --foreground-scripts`);
  }

  // Verify SHA256 against the sidecar checksum file from the same release.
  let expected;
  try {
    expected = (await fetchText(`${base}/${assetName}.sha256`)).trim().split(/\s+/)[0];
  } catch (err) {
    warn(`Could not fetch sidecar checksum (${err.message}). Proceeding without verification is not best practice — aborting.`);
    process.exit(1);
  }

  const actual = sha256File(binPath);
  if (expected.toLowerCase() !== actual.toLowerCase()) {
    die(`Checksum mismatch for ${assetName}. expected=${expected} actual=${actual}`);
  }

  if (process.platform !== 'win32') {
    chmodSync(binPath, 0o755);
  }
  log(`Installed ${assetName} (${BINARY_VERSION}, sha256 ${actual.slice(0, 16)}...).`);
}

main().catch((err) => die(err.message || String(err)));
