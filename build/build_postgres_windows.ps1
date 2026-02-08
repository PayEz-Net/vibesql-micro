$ErrorActionPreference = "Stop"

$PgVersion = "16.1"
$PgZipUrl = "https://get.enterprisedb.com/postgresql/postgresql-${PgVersion}-1-windows-x64-binaries.zip"
$BuildDir = Join-Path (Get-Location) "postgres-build-windows"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$EmbedDir = Join-Path $ProjectRoot "internal\postgres\embed"

function Log-Info($msg)  { Write-Host "[INFO] $msg" -ForegroundColor Green }
function Log-Warn($msg)  { Write-Host "[WARN] $msg" -ForegroundColor Yellow }
function Log-Error($msg) { Write-Host "[ERROR] $msg" -ForegroundColor Red }

Write-Host "========================================="
Write-Host "VibeSQL PostgreSQL Build (Windows amd64)"
Write-Host "========================================="
Write-Host "Version: PostgreSQL $PgVersion"
Write-Host "Target: Windows amd64"
Write-Host "Source: Official EDB binaries"
Write-Host ""

if (-not (Test-Path $BuildDir)) {
    New-Item -ItemType Directory -Path $BuildDir | Out-Null
}

$ZipFile = Join-Path $BuildDir "postgresql-${PgVersion}-windows-x64.zip"

if (-not (Test-Path $ZipFile)) {
    Log-Info "Downloading PostgreSQL $PgVersion Windows binaries..."
    Log-Info "URL: $PgZipUrl"
    try {
        Invoke-WebRequest -Uri $PgZipUrl -OutFile $ZipFile -UseBasicParsing
        Log-Info "Download complete"
    } catch {
        Log-Error "Download failed: $_"
        Log-Info "Please download PostgreSQL $PgVersion Windows x64 binaries manually:"
        Log-Info "  https://www.enterprisedb.com/download-postgresql-binaries"
        Log-Info "Place the zip at: $ZipFile"
        Log-Info "Then run this script again."
        exit 1
    }
} else {
    Log-Warn "Zip already exists, skipping download"
}

$ExtractDir = Join-Path $BuildDir "extracted"
if (Test-Path $ExtractDir) {
    Remove-Item -Recurse -Force $ExtractDir
}

Log-Info "Extracting PostgreSQL binaries..."
Expand-Archive -Path $ZipFile -DestinationPath $ExtractDir

$PgBinDir = Join-Path $ExtractDir "pgsql\bin"
$PgLibDir = Join-Path $ExtractDir "pgsql\lib"
$PgShareDir = Join-Path $ExtractDir "pgsql\share"

if (-not (Test-Path $PgBinDir)) {
    $PgBinDir = Get-ChildItem -Path $ExtractDir -Recurse -Filter "postgres.exe" |
        Select-Object -First 1 |
        ForEach-Object { $_.DirectoryName }
    if (-not $PgBinDir) {
        Log-Error "Could not find postgres.exe in extracted archive"
        exit 1
    }
    $PgLibDir = Join-Path (Split-Path $PgBinDir) "lib"
    $PgShareDir = Join-Path (Split-Path $PgBinDir) "share"
}

Log-Info "Found PostgreSQL at: $PgBinDir"

if (-not (Test-Path $EmbedDir)) {
    New-Item -ItemType Directory -Path $EmbedDir | Out-Null
}

Log-Info "Copying postgres.exe..."
Copy-Item (Join-Path $PgBinDir "postgres.exe") (Join-Path $EmbedDir "postgres_micro_windows_amd64.exe")

Log-Info "Copying initdb.exe..."
Copy-Item (Join-Path $PgBinDir "initdb.exe") (Join-Path $EmbedDir "initdb_windows_amd64.exe")

Log-Info "Copying pg_ctl.exe..."
Copy-Item (Join-Path $PgBinDir "pg_ctl.exe") (Join-Path $EmbedDir "pg_ctl_windows_amd64.exe")

$LibpqSrc = Join-Path $PgBinDir "libpq.dll"
if (-not (Test-Path $LibpqSrc)) {
    $LibpqSrc = Get-ChildItem -Path $ExtractDir -Recurse -Filter "libpq*.dll" |
        Select-Object -First 1 |
        ForEach-Object { $_.FullName }
}

if ($LibpqSrc -and (Test-Path $LibpqSrc)) {
    Log-Info "Copying libpq DLL..."
    Copy-Item $LibpqSrc (Join-Path $EmbedDir "libpq-5.dll")
} else {
    Log-Warn "libpq DLL not found -- initdb may fail at runtime"
}

$DependencyDlls = @(
    "libintl*.dll",
    "libssl*.dll",
    "libcrypto*.dll",
    "libiconv*.dll",
    "zlib*.dll",
    "icu*.dll",
    "libwinpthread*.dll",
    "libzstd*.dll",
    "liblz4*.dll",
    "libxml2*.dll"
)

foreach ($pattern in $DependencyDlls) {
    $dlls = Get-ChildItem -Path $PgBinDir -Filter $pattern -ErrorAction SilentlyContinue
    foreach ($dll in $dlls) {
        Log-Info "Copying dependency: $($dll.Name)"
        Copy-Item $dll.FullName (Join-Path $EmbedDir $dll.Name)
    }
}

# Copy essential PostgreSQL extension DLLs from lib directory (for $libdir)
$LibExtDlls = @(
    "plpgsql.dll",
    "dict_snowball.dll"
)

foreach ($dllName in $LibExtDlls) {
    $dllPath = Join-Path $PgLibDir $dllName
    if (Test-Path $dllPath) {
        Log-Info "Copying lib extension: $dllName"
        Copy-Item $dllPath (Join-Path $EmbedDir $dllName)
    } else {
        Log-Warn "Extension DLL not found: $dllName"
    }
}

$ShareTar = Join-Path $EmbedDir "share.tar.gz"
if (-not (Test-Path $ShareTar)) {
    Log-Info "Creating share.tar.gz..."
    if (Get-Command tar -ErrorAction SilentlyContinue) {
        Push-Location (Split-Path $PgShareDir)
        tar -czf $ShareTar share/
        Pop-Location
        Log-Info "share.tar.gz created"
    } else {
        Log-Warn "tar not available -- share.tar.gz must be created manually"
        Log-Info "You can use Git Bash or WSL:"
        Log-Info "  cd $(Split-Path $PgShareDir) && tar -czf $ShareTar share/"
    }
} else {
    Log-Info "share.tar.gz already exists, skipping"
}

Write-Host ""
Write-Host "========================================="
Write-Host "Windows PostgreSQL Build Complete"
Write-Host "========================================="
Write-Host "Binaries placed in: $EmbedDir"
Get-ChildItem $EmbedDir | ForEach-Object {
    $sizeMB = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  $($_.Name)  ${sizeMB}MB"
}

$TotalSize = (Get-ChildItem $EmbedDir | Measure-Object -Property Length -Sum).Sum
$TotalMB = [math]::Round($TotalSize / 1MB, 2)
Write-Host ""
Write-Host "Total embed size: ${TotalMB}MB"

if ($TotalMB -gt 70) {
    Log-Error "Total embed size exceeds 70MB limit!"
} elseif ($TotalMB -gt 65) {
    Log-Warn "Total embed size exceeds 65MB preferred target"
} else {
    Log-Info "Total embed size within target"
}

if ($env:CI -eq "true" -or $env:NONINTERACTIVE -eq "true") {
    Remove-Item -Recurse -Force $BuildDir
    Log-Info "Cleanup complete (non-interactive)"
} else {
    $Cleanup = Read-Host "Remove build artifacts? (y/N)"
    if ($Cleanup -eq "y" -or $Cleanup -eq "Y") {
        Remove-Item -Recurse -Force $BuildDir
        Log-Info "Cleanup complete"
    }
}
