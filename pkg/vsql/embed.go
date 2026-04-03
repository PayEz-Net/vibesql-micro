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
//go:embed embed/libpq-5.dll
//go:embed embed/libcrypto-3-x64.dll
//go:embed embed/libssl-3-x64.dll
//go:embed embed/libiconv-2.dll
//go:embed embed/libintl-9.dll
//go:embed embed/libwinpthread-1.dll
//go:embed embed/zlib1.dll
//go:embed embed/liblz4.dll
//go:embed embed/libzstd.dll
//go:embed embed/libxml2.dll
//go:embed embed/icudt67.dll
//go:embed embed/icuin67.dll
//go:embed embed/icuuc67.dll
//go:embed embed/icuio67.dll
//go:embed embed/icutu67.dll
//go:embed embed/dict_snowball.dll
//go:embed embed/plpgsql.dll
var embeddedFiles embed.FS

// dlls is the list of DLLs to extract
var dlls = []string{
	"vcruntime140.dll",
	"msvcp140.dll",
	"libpq-5.dll",
	"libcrypto-3-x64.dll",
	"libssl-3-x64.dll",
	"libiconv-2.dll",
	"libintl-9.dll",
	"libwinpthread-1.dll",
	"zlib1.dll",
	"liblz4.dll",
	"libzstd.dll",
	"libxml2.dll",
	"icudt67.dll",
	"icuin67.dll",
	"icuuc67.dll",
	"icuio67.dll",
	"icutu67.dll",
	"dict_snowball.dll",
	"plpgsql.dll",
}

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
	
	// Extract all DLLs
	for _, dll := range dlls {
		if err := extractFile("embed/"+dll, filepath.Join(binDir, dll)); err != nil {
			// Log but continue - some DLLs might already be on system
			fmt.Fprintf(os.Stderr, "Warning: could not extract %s: %v\n", dll, err)
		}
	}
	
	// Also copy libpq to LIBPQ.dll (uppercase) - Windows PostgreSQL looks for this name
	libpqData, err := embeddedFiles.ReadFile("embed/libpq-5.dll")
	if err == nil {
		_ = os.WriteFile(filepath.Join(binDir, "LIBPQ.dll"), libpqData, 0644)
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
