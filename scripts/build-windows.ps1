$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

$Version = if ($env:VERSION) { $env:VERSION } else { "1.0.0" }
$GitCommit = if ($env:GIT_COMMIT) { $env:GIT_COMMIT } else {
    try { git -C $ProjectRoot rev-parse --short HEAD 2>$null } catch { "dev" }
}
$BuildDate = if ($env:BUILD_DATE) { $env:BUILD_DATE } else {
    (Get-Date -Format "yyyy-MM-dd_HH:mm:ss").ToString()
}

$LdFlags = "-s -w"
$LdFlags += " -X github.com/vibesql/vibe/internal/version.Version=$Version"
$LdFlags += " -X github.com/vibesql/vibe/internal/version.GitCommit=$GitCommit"
$LdFlags += " -X github.com/vibesql/vibe/internal/version.BuildDate=$BuildDate"

$Output = Join-Path $ProjectRoot "vibe-windows-amd64.exe"

Write-Host "Building VibeSQL for Windows amd64..." -ForegroundColor Yellow
Write-Host "  Version:    $Version"
Write-Host "  Commit:     $GitCommit"
Write-Host "  Build Date: $BuildDate"

$RequiredEmbeds = @(
    "internal\postgres\embed\postgres_micro_windows_amd64.exe",
    "internal\postgres\embed\initdb_windows_amd64.exe",
    "internal\postgres\embed\pg_ctl_windows_amd64.exe",
    "internal\postgres\embed\libpq-5.dll",
    "internal\postgres\embed\share.tar.gz"
)

foreach ($f in $RequiredEmbeds) {
    $FullPath = Join-Path $ProjectRoot $f
    if (-not (Test-Path $FullPath)) {
        Write-Host "ERROR: Missing embedded file: $f" -ForegroundColor Red
        Write-Host "Run the PostgreSQL build first: .\build\build_postgres_windows.ps1"
        exit 1
    }
}

Push-Location $ProjectRoot
try {
    $env:CGO_ENABLED = "0"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    & go build -ldflags="$LdFlags" -o $Output ./cmd/vibe
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Build failed!" -ForegroundColor Red
        exit 1
    }
} finally {
    Pop-Location
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
}

$Size = (Get-Item $Output).Length
$SizeMB = [math]::Round($Size / 1MB, 2)

Write-Host ""
Write-Host "Build complete: $Output" -ForegroundColor Green
Write-Host "  Size: ${SizeMB}MB ($Size bytes)"

if ($Size -gt 26214400) {
    Write-Host "ERROR: Binary exceeds 25MB hard limit!" -ForegroundColor Red
    exit 1
} elseif ($Size -gt 20971520) {
    Write-Host "WARNING: Binary exceeds 20MB preferred target" -ForegroundColor Yellow
} else {
    Write-Host "  Size OK (under 20MB target)" -ForegroundColor Green
}

Write-Host ""
Write-Host "To test: $Output serve"
