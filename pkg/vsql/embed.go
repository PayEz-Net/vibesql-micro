//go:build windows
// +build windows

package vsql

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:embed embed/postgres_micro_windows_amd64.exe
//go:embed embed/initdb_windows_amd64.exe
//go:embed embed/share.tar.gz
//go:embed embed/vcruntime140.dll
//go:embed embed/msvcp140.dll
var embeddedFiles embed.FS

// ensureBinary extracts postgres binaries if needed and returns the paths
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
			if _, err := os.Stat(sharePath); err == nil {
				// Already extracted
				return postgresPath, initdbPath, sharePath, false, nil
			}
		}
	}
	
	// Need to extract
	firstTime = true
	if progress != nil {
		progress("Setting up vibesql")
	}
	
	// Create directory
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return "", "", "", false, fmt.Errorf("create bin directory: %w", err)
	}
	
	// Extract postgres.exe
	if err := extractFile("embed/postgres_micro_windows_amd64.exe", filepath.Join(binDir, "postgres.exe")); err != nil {
		return "", "", "", false, err
	}
	
	// Extract initdb.exe
	if err := extractFile("embed/initdb_windows_amd64.exe", filepath.Join(binDir, "initdb.exe")); err != nil {
		return "", "", "", false, err
	}
	
	// Extract MSVC runtime DLLs
	if err := extractFile("embed/vcruntime140.dll", filepath.Join(binDir, "vcruntime140.dll")); err != nil {
		// Non-fatal - may already be on system
		fmt.Fprintf(os.Stderr, "Warning: could not extract vcruntime140.dll: %v\n", err)
	}
	if err := extractFile("embed/msvcp140.dll", filepath.Join(binDir, "msvcp140.dll")); err != nil {
		// Non-fatal - may already be on system
		fmt.Fprintf(os.Stderr, "Warning: could not extract msvcp140.dll: %v\n", err)
	}
	
	// Extract share.tar.gz
	shareData, err := embeddedFiles.ReadFile("embed/share.tar.gz")
	if err != nil {
		return "", "", "", false, fmt.Errorf("read embedded share: %w", err)
	}
	
	if err := extractTarGz(shareData, binDir); err != nil {
		return "", "", "", false, fmt.Errorf("extract share: %w", err)
	}
	
	return postgresPath, initdbPath, sharePath, firstTime, nil
}

// extractFile extracts a single embedded file
func extractFile(embedPath, destPath string) error {
	data, err := embeddedFiles.ReadFile(embedPath)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", embedPath, err)
	}
	
	if err := os.WriteFile(destPath, data, 0755); err != nil {
		return fmt.Errorf("write %s: %w", destPath, err)
	}
	
	return nil
}

// extractTarGz extracts a tar.gz archive
func extractTarGz(data []byte, dst string) error {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gr.Close()
	
	tr := tar.NewReader(gr)
	
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		
		target := filepath.Join(dst, header.Name)
		
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	
	return nil
}

var version = "0.1.0"
