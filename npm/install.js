#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const https = require('https');

const package = require('./package.json');
const version = package.version;

// Detect platform and architecture
const platform = process.platform;
const arch = process.arch;

// Map Node.js platform/arch to Go build names
const platformMap = {
  darwin: 'Darwin',
  linux: 'Linux',
  win32: 'Windows'
};

const archMap = {
  x64: 'x86_64',
  arm64: 'arm64'
};

const goPlatform = platformMap[platform];
const goArch = archMap[arch];

if (!goPlatform || !goArch) {
  console.error(`Unsupported platform: ${platform} ${arch}`);
  process.exit(1);
}

const ext = platform === 'win32' ? '.zip' : '.tar.gz';
const filename = `lanup_${version}_${goPlatform}_${goArch}${ext}`;
const downloadUrl = `https://github.com/raucheacho/lanup/releases/download/v${version}/${filename}`;

console.log(`Downloading lanup v${version} for ${platform} ${arch}...`);
console.log(`URL: ${downloadUrl}`);

const binDir = path.join(__dirname, 'bin');
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Download and extract
const tempFile = path.join(binDir, filename);

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        download(response.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      
      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode}`));
        return;
      }

      response.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
}

async function install() {
  try {
    await download(downloadUrl, tempFile);
    console.log('Download complete. Extracting...');

    // Extract
    if (platform === 'win32') {
      execSync(`powershell -command "Expand-Archive -Path '${tempFile}' -DestinationPath '${binDir}' -Force"`, {
        stdio: 'inherit'
      });
    } else {
      execSync(`tar -xzf "${tempFile}" -C "${binDir}"`, {
        stdio: 'inherit'
      });
    }

    // Make executable on Unix
    if (platform !== 'win32') {
      const binaryPath = path.join(binDir, 'lanup');
      fs.chmodSync(binaryPath, 0o755);
    }

    // Clean up
    fs.unlinkSync(tempFile);

    console.log('✓ lanup installed successfully!');
    console.log('\nRun: npx lanup --version');
  } catch (error) {
    console.error('Installation failed:', error.message);
    console.error('\nYou can manually download from:');
    console.error(`https://github.com/raucheacho/lanup/releases/tag/v${version}`);
    process.exit(1);
  }
}

install();
