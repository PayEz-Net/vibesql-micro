package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version == "" {
		t.Error("Version should not be empty")
	}

	if info.GitCommit == "" {
		t.Error("GitCommit should not be empty")
	}

	if info.BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}

	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}

	if info.OS == "" {
		t.Error("OS should not be empty")
	}

	if info.Arch == "" {
		t.Error("Arch should not be empty")
	}
}

func TestInfo_String(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2024-01-01",
		GoVersion: "go1.21.0",
		OS:        "linux",
		Arch:      "amd64",
	}

	result := info.String()
	expected := "VibeSQL 1.0.0 (commit: abc123, built: 2024-01-01, go: go1.21.0, linux/amd64)"

	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}

func TestInfo_Short(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2024-01-01",
		GoVersion: "go1.21.0",
		OS:        "linux",
		Arch:      "amd64",
	}

	result := info.Short()
	expected := "1.0.0"

	if result != expected {
		t.Errorf("Short() = %q, want %q", result, expected)
	}
}

func TestInfo_Full(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2024-01-01",
		GoVersion: "go1.21.0",
		OS:        "linux",
		Arch:      "amd64",
	}

	result := info.Full()

	// Check that all expected fields are present
	expectedFields := []string{
		"VibeSQL Version Information:",
		"Version:    1.0.0",
		"Git Commit: abc123",
		"Build Date: 2024-01-01",
		"Go Version: go1.21.0",
		"OS/Arch:    linux/amd64",
	}

	for _, field := range expectedFields {
		if !strings.Contains(result, field) {
			t.Errorf("Full() missing expected field: %q", field)
		}
	}
}

func TestVersionConstants(t *testing.T) {
	// Test that default constants are set
	if Version == "" {
		t.Error("Version constant should not be empty")
	}

	if GitCommit == "" {
		t.Error("GitCommit constant should not be empty")
	}

	if BuildDate == "" {
		t.Error("BuildDate constant should not be empty")
	}

	if GoVersion != runtime.Version() {
		t.Errorf("GoVersion = %q, want %q", GoVersion, runtime.Version())
	}
}

func TestInfo_AllFormatsMaintainConsistency(t *testing.T) {
	info := Get()

	// All format methods should use the same data
	short := info.Short()
	str := info.String()
	full := info.Full()

	// Short should just be the version
	if short != info.Version {
		t.Errorf("Short() = %q, but Version = %q", short, info.Version)
	}

	// String and Full should contain the version
	if !strings.Contains(str, info.Version) {
		t.Error("String() should contain Version")
	}

	if !strings.Contains(full, info.Version) {
		t.Error("Full() should contain Version")
	}

	// String and Full should contain the git commit
	if !strings.Contains(str, info.GitCommit) {
		t.Error("String() should contain GitCommit")
	}

	if !strings.Contains(full, info.GitCommit) {
		t.Error("Full() should contain GitCommit")
	}
}

func TestInfo_OSAndArchMatch(t *testing.T) {
	info := Get()

	if info.OS != runtime.GOOS {
		t.Errorf("Info.OS = %q, want %q", info.OS, runtime.GOOS)
	}

	if info.Arch != runtime.GOARCH {
		t.Errorf("Info.Arch = %q, want %q", info.Arch, runtime.GOARCH)
	}
}

func TestInfo_EmptyValues(t *testing.T) {
	// Test behavior with empty/default values
	info := Info{
		Version:   "",
		GitCommit: "",
		BuildDate: "",
		GoVersion: "",
		OS:        "",
		Arch:      "",
	}

	// Should not panic
	_ = info.String()
	_ = info.Short()
	_ = info.Full()
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Get()
	}
}

func BenchmarkInfo_String(b *testing.B) {
	info := Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = info.String()
	}
}

func BenchmarkInfo_Full(b *testing.B) {
	info := Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = info.Full()
	}
}
