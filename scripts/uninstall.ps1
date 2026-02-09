$ErrorActionPreference = "Stop"

$BinaryName = "vibe.exe"
$InstallDir = if ($env:VIBESQL_INSTALL_DIR) { $env:VIBESQL_INSTALL_DIR } else {
    Join-Path $env:LOCALAPPDATA "VibeSQL"
}
$DataDir = ".\vibe-data"

Write-Host "VibeSQL Uninstaller" -ForegroundColor Yellow
Write-Host "================================"

$BinaryPath = Join-Path $InstallDir $BinaryName

if (-not (Test-Path $BinaryPath)) {
    Write-Host "VibeSQL is not installed at $BinaryPath"
    exit 0
}

Write-Host "This will remove:"
Write-Host "  - $BinaryPath"

$RemoveData = $false
if (Test-Path $DataDir) {
    Write-Host "  - $DataDir (DATABASE DATA)" -ForegroundColor Red
    Write-Host ""
    $Response = Read-Host "Remove database data too? [y/N]"
    if ($Response -eq "y" -or $Response -eq "Y") {
        $RemoveData = $true
    }
}

Write-Host ""
$Confirm = Read-Host "Proceed with uninstall? [y/N]"
if ($Confirm -ne "y" -and $Confirm -ne "Y") {
    Write-Host "Cancelled."
    exit 0
}

Remove-Item $BinaryPath -Force
Write-Host "Binary removed" -ForegroundColor Green

if ($RemoveData -and (Test-Path $DataDir)) {
    Remove-Item -Recurse -Force $DataDir
    Write-Host "Database data removed" -ForegroundColor Green
}

$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -like "*$InstallDir*") {
    $NewPath = ($UserPath -split ";" | Where-Object { $_ -ne $InstallDir }) -join ";"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Host "Removed $InstallDir from user PATH" -ForegroundColor Green
}

$RemainingFiles = Get-ChildItem $InstallDir -ErrorAction SilentlyContinue
if (-not $RemainingFiles) {
    Remove-Item $InstallDir -Force -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "VibeSQL uninstalled successfully" -ForegroundColor Green
