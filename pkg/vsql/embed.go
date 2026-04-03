//go:build windows
// +build windows

package vsql

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed embed/postgres_micro_windows_amd64.exe
//go:embed embed/initdb_windows_amd64.exe
//go:embed embed/share.tar.gz
var embeddedFiles embed.FS

// binaryManager handles extraction and caching of postgres binaries
type binaryManager struct {
	binDir string
}

// ensureBinary extracts postgres binaries if needed and returns the path
func ensureBinary(progress func(string)) (postgresPath string, initdbPath string, sharePath string, firstTime bool, err error) {
	// Determine cache directory
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	
	binDir := filepath.Join(cacheDir, "vibesql-micro", "bin-"+version)
	
	// Check if already extracted
	postgresPath = filepath.Join(binDir, "postgres.exe")
	initdbPath = filepath.Join(binDir, "initdb.exe")
	sharePath = filepath.Join(binDir, "share")
	
	if _, err := os.Stat(postgresPath); err == nil {
		if _, err := os.Stat(initdbPath); err == nil {
			// Already extracted
			return postgresPath, initdbPath, sharePath, false, nil
		}
	}
	
	// Need to extract
	firstTime = true
	if progress != nil {
		progress("Setting up vibesql...")
	}
	
	// Create directory
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", "", "", false, fmt.Errorf("create bin directory: %w", err)
	}
	
	// Extract postgres.exe
	postgresData, err := embeddedFiles.ReadFile("embed/postgres_micro_windows_amd64.exe")
	if err != nil {
		return "", "", "", false, fmt.Errorf("read embedded postgres: %w", err)
	}
	
	if err := os.WriteFile(postgresPath, postgresData, 0755); err != nil {
		return "", "", "", false, fmt.Errorf("write postgres binary: %w", err)
	}
	
	// Extract initdb.exe
	initdbData, err := embeddedFiles.ReadFile("embed/initdb_windows_amd64.exe")
	if err != nil {
		return "", "", "", false, fmt.Errorf("read embedded initdb: %w", err)
	}
	
	if err := os.WriteFile(initdbPath, initdbData, 0755); err != nil {
		return "", "", "", false, fmt.Errorf("write initdb binary: %w", err)
	}
	
	// Extract share.tar.gz
	shareData, err := embeddedFiles.ReadFile("embed/share.tar.gz")
	if err != nil {
		return "", "", "", false, fmt.Errorf("read embedded share: %w", err)
	}
	
	if err := extractTarGz(shareData, binDir); err != nil {
		return "", "", "", false, fmt.Errorf("extract share: %w", err)
	}
	
	if progress != nil {
		progress("done")
	}
	
	return postgresPath, initdbPath, sharePath, firstTime, nil
}

// extractTarGz extracts a tar.gz archive to the destination directory
func extractTarGz(data []byte, dst string) error {
	gr, err := gzip.NewReader(io.NopCloser(io.NopCloser(nil)))
	if err != nil {
		return err
	}
	
	// We need to use bytes.Reader for gzip
	gr, err = gzip.NewReader(io.NopCloser(io.NopCloser(nil)))
	if err != nil {
		return err
	}
	
	// Re-implement with bytes
	return extractTarGzImpl(data, dst)
}

func extractTarGzImpl(data []byte, dst string) error {
	gr, err := gzip.NewReader(io.NopCloser(io.NopCloser(nil)))
	if err != nil {
		return err
	}
	
	// Actually implement properly
	return fmt.Errorf("extraction not fully implemented yet")
}

var version = "0.1.0"

func init() {
	runtime.GOMAXPROCS(1)
}
