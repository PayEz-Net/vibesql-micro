package postgres

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDetectCompiledShareDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name      string
		sharePath string
	}{
		{"standard prefix", "/opt/postgres_micro/share"},
		{"usr local", "/usr/local/pgsql/share"},
		{"tmp install", "/tmp/postgres_micro_install/share"},
		{"deep nesting", "/home/user/builds/pg16/share"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build a fake binary: ... \x00 <path> \x00 ...
			fakeBinary := make([]byte, 1024)
			fakeBinary[99] = 0
			copy(fakeBinary[100:], []byte(tt.sharePath))
			fakeBinary[100+len(tt.sharePath)] = 0

			binPath := filepath.Join(tmpDir, "postgres-"+tt.name)
			if err := os.WriteFile(binPath, fakeBinary, 0755); err != nil {
				t.Fatalf("write: %v", err)
			}

			detected, err := detectCompiledShareDir(binPath)
			if err != nil {
				t.Fatalf("detectCompiledShareDir: %v", err)
			}
			if detected != tt.sharePath {
				t.Errorf("detected = %q, want %q", detected, tt.sharePath)
			}
		})
	}
}

func TestDetectCompiledShareDir_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fakeBinary := make([]byte, 256)
	copy(fakeBinary[10:], []byte("no share path here"))
	binPath := filepath.Join(tmpDir, "postgres")
	if err := os.WriteFile(binPath, fakeBinary, 0755); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, err = detectCompiledShareDir(binPath)
	if err == nil {
		t.Error("expected error when share path not found")
	}
}

func TestDetectCompiledShareDir_RejectsRelativePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// "share" without an absolute prefix should not match
	fakeBinary := make([]byte, 256)
	fakeBinary[49] = 0
	copy(fakeBinary[50:], []byte("relative/share"))
	fakeBinary[50+len("relative/share")] = 0
	binPath := filepath.Join(tmpDir, "postgres")
	if err := os.WriteFile(binPath, fakeBinary, 0755); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, err = detectCompiledShareDir(binPath)
	if err == nil {
		t.Error("expected error for relative path")
	}
}

func TestPatchBinaryPaths_SinglePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Simulate: \x00/opt/postgres_micro/share\x00
	oldPath := "/opt/postgres_micro/share"
	newPath := "/tmp/vb-1234567/share"
	content := make([]byte, 256)
	offset := 50
	content[offset-1] = 0
	copy(content[offset:], []byte(oldPath))
	content[offset+len(oldPath)] = 0

	binPath := filepath.Join(tmpDir, "test-binary")
	if err := os.WriteFile(binPath, content, 0755); err != nil {
		t.Fatalf("write: %v", err)
	}

	err = patchBinaryPaths(binPath, [][2]string{{oldPath, newPath}})
	if err != nil {
		t.Fatalf("patchBinaryPaths: %v", err)
	}

	patched, _ := os.ReadFile(binPath)

	// Old path should be gone
	if bytes.Contains(patched, []byte(oldPath)) {
		t.Error("patched binary still contains old path")
	}

	// Read the null-terminated string at the same offset
	end := offset
	for end < len(patched) && patched[end] != 0 {
		end++
	}
	result := string(patched[offset:end])

	if result != newPath {
		t.Errorf("patched path = %q, want %q", result, newPath)
	}
}

func TestPatchBinaryPaths_MultiplePaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	prefix := "/opt/postgres_micro"
	newPrefix := "/tmp/vb-12345"

	// Binary layout: ... \x00<share>\x00 ... \x00<lib>\x00 ... \x00<share/locale>\x00 ...
	content := make([]byte, 512)

	// Path 1: /opt/postgres_micro/share at offset 50
	p1 := prefix + "/share"
	content[49] = 0
	copy(content[50:], []byte(p1))
	content[50+len(p1)] = 0

	// Path 2: /opt/postgres_micro/lib at offset 150
	p2 := prefix + "/lib"
	content[149] = 0
	copy(content[150:], []byte(p2))
	content[150+len(p2)] = 0

	// Path 3: /opt/postgres_micro/share/locale at offset 250
	p3 := prefix + "/share/locale"
	content[249] = 0
	copy(content[250:], []byte(p3))
	content[250+len(p3)] = 0

	binPath := filepath.Join(tmpDir, "test-binary")
	if err := os.WriteFile(binPath, content, 0755); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Replace in correct order: longer suffixes first
	replacements := [][2]string{
		{prefix + "/share/locale", newPrefix + "/share/locale"},
		{prefix + "/share", newPrefix + "/share"},
		{prefix + "/lib", newPrefix + "/lib"},
	}

	err = patchBinaryPaths(binPath, replacements)
	if err != nil {
		t.Fatalf("patchBinaryPaths: %v", err)
	}

	patched, _ := os.ReadFile(binPath)

	// Verify each patched string reads correctly
	readNullTerminated := func(data []byte, offset int) string {
		end := offset
		for end < len(data) && data[end] != 0 {
			end++
		}
		return string(data[offset:end])
	}

	got1 := readNullTerminated(patched, 50)
	want1 := newPrefix + "/share"
	if got1 != want1 {
		t.Errorf("path1: got %q, want %q", got1, want1)
	}

	got2 := readNullTerminated(patched, 150)
	want2 := newPrefix + "/lib"
	if got2 != want2 {
		t.Errorf("path2: got %q, want %q", got2, want2)
	}

	got3 := readNullTerminated(patched, 250)
	want3 := newPrefix + "/share/locale"
	if got3 != want3 {
		t.Errorf("path3: got %q, want %q", got3, want3)
	}

	// Old prefix should be completely gone
	if bytes.Contains(patched, []byte(prefix)) {
		t.Error("patched binary still contains old prefix")
	}
}

func TestPatchBinaryPaths_NothingToPatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	content := []byte("no matching paths here at all")
	binPath := filepath.Join(tmpDir, "test-binary")
	if err := os.WriteFile(binPath, content, 0755); err != nil {
		t.Fatalf("write: %v", err)
	}

	err = patchBinaryPaths(binPath, [][2]string{
		{"/opt/postgres_micro/share", "/tmp/new/share"},
	})
	if err != nil {
		t.Fatalf("should succeed with nothing to patch: %v", err)
	}
}

func TestPatchBinaryPaths_PreservesFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix permissions not supported on Windows")
	}

	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldPath := "/opt/pg/share"
	content := make([]byte, 64)
	content[0] = 0
	copy(content[1:], []byte(oldPath))
	content[1+len(oldPath)] = 0

	binPath := filepath.Join(tmpDir, "test-binary")
	if err := os.WriteFile(binPath, content, 0755); err != nil {
		t.Fatalf("write: %v", err)
	}

	err = patchBinaryPaths(binPath, [][2]string{
		{oldPath, "/tmp/x/share"},
	})
	if err != nil {
		t.Fatalf("patchBinaryPaths: %v", err)
	}

	info, _ := os.Stat(binPath)
	if info.Mode().Perm() != 0755 {
		t.Errorf("permissions changed: got %o, want 0755", info.Mode().Perm())
	}
}

func TestRelocateBinaries_NilShareDir(t *testing.T) {
	m := &Manager{}
	err := relocateBinaries(m)
	if err != nil {
		t.Errorf("should be no-op with empty shareDir: %v", err)
	}
}

func TestRelocateBinaries_PathsAlreadyMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "vibe-relocate-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	shareDir := tmpDir + "/share"
	// Create a binary where the compiled-in path matches shareDir
	fakeBinary := make([]byte, 256)
	fakeBinary[49] = 0
	copy(fakeBinary[50:], []byte(shareDir))
	fakeBinary[50+len(shareDir)] = 0

	binPath := filepath.Join(tmpDir, "postgres")
	if err := os.WriteFile(binPath, fakeBinary, 0755); err != nil {
		t.Fatalf("write: %v", err)
	}

	m := &Manager{
		tmpDir:          tmpDir,
		shareDir:        shareDir,
		postgresBinPath: binPath,
		initdbBinPath:   binPath, // reuse for test
	}

	err = relocateBinaries(m)
	if err != nil {
		t.Errorf("should be no-op when paths match: %v", err)
	}
}
