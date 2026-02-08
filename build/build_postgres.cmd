@echo off
REM Build minimal PostgreSQL binary using Docker
REM Target: Linux x64 (amd64)
REM Size Goal: <=20MB

setlocal enabledelayedexpansion

echo [INFO] Building minimal PostgreSQL binary for Linux x64...
echo.

REM Check if Docker is available
docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not installed or not running
    echo [ERROR] Please install Docker Desktop for Windows
    exit /b 1
)

echo [INFO] Docker found
echo.

REM Build the builder image
echo [INFO] Building PostgreSQL builder Docker image...
docker build -f Dockerfile.postgres-builder -t vibesql-postgres-builder .
if errorlevel 1 (
    echo [ERROR] Docker build failed
    exit /b 1
)

echo.
echo [INFO] Docker image built successfully
echo.

REM Run the builder and extract binary
echo [INFO] Compiling PostgreSQL...
docker run --rm -v "%cd%:/output" vibesql-postgres-builder
if errorlevel 1 (
    echo [ERROR] PostgreSQL compilation failed
    exit /b 1
)

echo.
echo [INFO] Checking binary size...

REM Get file size (Windows compatible)
for %%F in (postgres_micro_linux_amd64) do set SIZE=%%~zF
set /a SIZE_MB=!SIZE! / 1024 / 1024

echo [INFO] Binary size: !SIZE_MB!MB (!SIZE! bytes)
echo.

REM Check if within 20MB limit
set /a MAX_SIZE=20 * 1024 * 1024
if !SIZE! GTR !MAX_SIZE! (
    echo [ERROR] Binary size !SIZE_MB!MB exceeds 20MB limit
    exit /b 1
)

echo [SUCCESS] PostgreSQL binary built successfully!
echo [SUCCESS] Output: postgres_micro_linux_amd64 (!SIZE_MB!MB)
echo.
echo Next steps:
echo   1. Copy binary to internal/postgres/embed/
echo   2. Update embedding code in internal/postgres/embed.go
echo   3. Build VibeSQL with embedded PostgreSQL

exit /b 0
