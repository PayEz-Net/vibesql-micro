//go:build linux
// +build linux

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

//go:embed embed/postgres_micro_linux_amd64
//go:embed embed/initdb_linux_amd64
//go:embed embed/share.tar.gz
//go:embed embed/libpq.so.5
//go:embed embed/plpgsql.so
//go:embed embed/dict_snowball.so
var embeddedFiles embed.FS

// linuxLibs is the list of shared libraries to extract
var linuxLibs = []string{
	"libpq.so.5",
	"plpgsql.so",
	"dict_snowball.so",
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
	postgresPath = filepath.Join(binDir, "postgres")
	initdbPath = filepath.Join(binDir, "initdb")
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
	
	// Extract postgres
	if err := extractFile("embed/postgres_micro_linux_amd64", filepath.Join(binDir, "postgres")); err != nil {
		return "", "", "", false, err
	}
	
	// Extract initdb
	if err := extractFile("embed/initdb_linux_amd64", filepath.Join(binDir, "initdb")); err != nil {
		return "", "", "", false, err
	}
	
	// Extract shared libraries
	for _, lib := range linuxLibs {
		if err := extractFile("embed/"+lib, filepath.Join(binDir, lib)); err != nil {
			// Log but continue - some libs might already be on system
			fmt.Fprintf(os.Stderr, "Warning: could not extract %s: %v\n", lib, err)
		}
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
