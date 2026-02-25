package postgres

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// relocateBinaries patches the compiled-in PKGDATADIR and LIBDIR paths in the
// extracted postgres and initdb binaries so they point at the actual runtime
// extraction directory.
//
// PostgreSQL bakes --prefix paths into binaries at compile time. When we extract
// the embedded binary to a temp dir, the compiled-in /opt/postgres_micro/share
// (or whatever prefix was used) doesn't exist. initdb's -L flag only helps
// initdb itself find catalog templates — but the `postgres --boot` subprocess
// that initdb spawns resolves timezone files from the compiled-in PKGDATADIR.
//
// The fix: after extraction, scan each binary for compiled-in paths and
// overwrite each full path in-place with the corresponding new path
// (null-padded to the original length). This is the same technique used by
// embedded-postgres-go and other production embedded PG libraries.
func relocateBinaries(m *Manager) error {
	if m.shareDir == "" || m.tmpDir == "" {
		return nil // system install or no extraction — nothing to relocate
	}

	// Detect the compiled-in share path by scanning the postgres binary.
	compiledShareDir, err := detectCompiledShareDir(m.postgresBinPath)
	if err != nil {
		log.Printf("[WARN] Could not detect compiled-in share dir: %v", err)
		return nil // non-fatal — the binary might work if paths happen to match
	}

	if compiledShareDir == m.shareDir {
		return nil // paths already match
	}

	// Derive the compiled-in prefix (parent of share)
	// e.g. /opt/postgres_micro/share → /opt/postgres_micro
	compiledPrefix := compiledShareDir[:len(compiledShareDir)-len("/share")]

	log.Printf("[INFO] Relocating compiled-in prefix %q -> %q", compiledPrefix, m.tmpDir)

	// PostgreSQL bakes these full paths as null-terminated strings:
	//   <prefix>/share    (PKGDATADIR — timezone, catalog templates)
	//   <prefix>/lib      (PKGLIBDIR — extensions like plpgsql.so)
	//   <prefix>/etc      (SYSCONFDIR)
	//   <prefix>/include  (INCLUDEDIR)
	//   <prefix>/bin      (BINDIR)
	//   <prefix>/include/server (PKGINCLUDEDIR)
	//   <prefix>/share/locale   (LOCALEDIR)
	//   <prefix>/share/man      (MANDIR)
	//   <prefix>/share/doc      (DOCDIR)
	//
	// We replace each full path individually so the null terminator stays correct.
	pathSuffixes := []string{
		"/share/locale",
		"/share/man",
		"/share/doc",
		"/share",
		"/include/server",
		"/include",
		"/lib",
		"/etc",
		"/bin",
	}

	// Build the replacement map: old full path → new full path
	replacements := make([][2]string, 0, len(pathSuffixes))
	for _, suffix := range pathSuffixes {
		oldPath := compiledPrefix + suffix
		newPath := m.tmpDir + suffix
		if len(newPath) > len(oldPath) {
			log.Printf("[WARN] New path %q longer than compiled %q — skipping (set TMPDIR shorter)", newPath, oldPath)
			continue
		}
		replacements = append(replacements, [2]string{oldPath, newPath})
	}

	if len(replacements) == 0 {
		return fmt.Errorf(
			"runtime path %q is longer than compiled prefix %q — "+
				"cannot relocate; set TMPDIR to a shorter path",
			m.tmpDir, compiledPrefix)
	}

	// Patch each binary
	binaries := []string{m.postgresBinPath, m.initdbBinPath}
	if m.pgCtlBinPath != "" {
		binaries = append(binaries, m.pgCtlBinPath)
	}

	for _, binPath := range binaries {
		if binPath == "" {
			continue
		}
		if err := patchBinaryPaths(binPath, replacements); err != nil {
			return fmt.Errorf("failed to patch %s: %w", binPath, err)
		}
	}

	return nil
}

// detectCompiledShareDir scans a postgres binary for the compiled-in
// PKGDATADIR. It searches for a null-terminated string ending with "/share\0"
// that looks like an absolute install path.
func detectCompiledShareDir(binaryPath string) (string, error) {
	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return "", fmt.Errorf("read binary: %w", err)
	}

	// Search for "/share\x00" — the null-terminated PKGDATADIR string.
	marker := []byte("/share\x00")
	offset := 0
	for {
		idx := bytes.Index(data[offset:], marker)
		if idx < 0 {
			break
		}
		absIdx := offset + idx

		// Walk backwards to find the start of the null-terminated string
		start := absIdx
		for start > 0 && data[start-1] != 0 {
			start--
		}

		candidate := string(data[start : absIdx+len("/share")])

		// Must be an absolute path with at least 2 slashes (e.g. /opt/x/share)
		if len(candidate) > 6 && candidate[0] == '/' {
			slashCount := 0
			for _, c := range candidate {
				if c == '/' {
					slashCount++
				}
			}
			if slashCount >= 2 {
				return candidate, nil
			}
		}

		offset = absIdx + len(marker)
	}

	return "", fmt.Errorf("compiled-in share path not found in binary")
}

// patchBinaryPaths replaces multiple full path strings in a binary file.
// Each new path is null-padded to match the original's byte length so the
// binary layout (and all offsets) remain unchanged.
func patchBinaryPaths(path string, replacements [][2]string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	totalPatched := 0
	for _, r := range replacements {
		oldPath, newPath := r[0], r[1]
		old := []byte(oldPath)
		count := bytes.Count(data, old)
		if count == 0 {
			continue
		}

		// Build replacement: new path null-padded to same byte length
		replacement := make([]byte, len(old))
		copy(replacement, []byte(newPath))

		data = bytes.ReplaceAll(data, old, replacement)
		totalPatched += count
	}

	if totalPatched == 0 {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}

	if err := os.WriteFile(path, data, info.Mode()); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("[INFO] Patched %d path occurrences in %s", totalPatched, path)
	return nil
}
