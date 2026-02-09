$ErrorActionPreference = "Stop"

$BinaryName = "vibe.exe"
$InstallDir = if ($env:VIBESQL_INSTALL_DIR) { $env:VIBESQL_INSTALL_DIR } else {
    Join-Path $env:LOCALAPPDATA "VibeSQL"
}
$BinaryFile = "vibe-windows-amd64.exe"

Write-Host "VibeSQL Installer" -ForegroundColor Green
Write-Host "================================"

$Arch = $env:PROCESSOR_ARCHITECTURE
if ($Arch -ne "AMD64") {
    Write-Host "Unsupported architecture: $Arch" -ForegroundColor Red
    Write-Host "VibeSQL supports: AMD64 (x86_64)"
    exit 1
}

Write-Host "Platform: windows/amd64"
Write-Host "Install to: $InstallDir\$BinaryName"
Write-Host ""

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$ExistingBin = Join-Path $InstallDir $BinaryName
if (Test-Path $ExistingBin) {
    try {
        $ExistingVersion = & $ExistingBin version 2>&1
        Write-Host "Existing installation found: $ExistingVersion" -ForegroundColor Yellow
        Write-Host "Upgrading..."
    } catch {
        Write-Host "Existing installation found (version unknown)" -ForegroundColor Yellow
    }
}

if (Test-Path ".\$BinaryFile") {
    Write-Host "Installing from local file..."
} else {
    Write-Host "Binary not found: .\$BinaryFile" -ForegroundColor Red
    Write-Host ""
    Write-Host "To install VibeSQL:"
    Write-Host "  1. Build from source:  .\scripts\build-windows.ps1"
    Write-Host "  2. Place the binary in the current directory"
    Write-Host "  3. Run this script again"
    exit 1
}

Write-Host "Verifying binary..."
try {
    $null = & ".\$BinaryFile" version 2>&1
} catch {
    Write-Host "Binary verification failed" -ForegroundColor Red
    exit 1
}

Copy-Item ".\$BinaryFile" (Join-Path $InstallDir $BinaryName) -Force

$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "Adding $InstallDir to user PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path = "$env:Path;$InstallDir"
    Write-Host "PATH updated (restart terminal for full effect)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "VibeSQL installed successfully!" -ForegroundColor Green
Write-Host ""
& (Join-Path $InstallDir $BinaryName) version
Write-Host ""
Write-Host "Quick start:"
Write-Host "  vibe serve        Start the server"
Write-Host "  vibe version      Show version info"
Write-Host "  vibe help         Show help"
Write-Host ""
Write-Host "API endpoint: http://127.0.0.1:5173/v1/query"
