package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/vibesql/vibe/internal/version"
)

func TestPrintVersion(t *testing.T) {
	output := captureOutput(func() {
		printVersion()
	})

	if !strings.Contains(output, "VibeSQL Version Information") {
		t.Errorf("Expected version header in output, got: %s", output)
	}

	info := version.Get()
	if !strings.Contains(output, info.Version) {
		t.Errorf("Expected version %s in output, got: %s", info.Version, output)
	}

	if !strings.Contains(output, "Git Commit:") {
		t.Errorf("Expected 'Git Commit:' in output, got: %s", output)
	}

	if !strings.Contains(output, "Build Date:") {
		t.Errorf("Expected 'Build Date:' in output, got: %s", output)
	}

	if !strings.Contains(output, "Go Version:") {
		t.Errorf("Expected 'Go Version:' in output, got: %s", output)
	}
}

func TestPrintUsage(t *testing.T) {
	output := captureOutput(func() {
		printUsage()
	})

	expectedStrings := []string{
		"VibeSQL",
		"Usage:",
		"Commands:",
		"serve",
		"version",
		"help",
		"Examples:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected '%s' in help output, got: %s", expected, output)
		}
	}
}

func TestPrintUsageDescribesServeCommand(t *testing.T) {
	output := captureOutput(func() {
		printUsage()
	})

	if !strings.Contains(output, "Start the HTTP server and embedded PostgreSQL") {
		t.Errorf("Expected serve command description in help output")
	}
}

func TestPrintUsageDescribesVersionCommand(t *testing.T) {
	output := captureOutput(func() {
		printUsage()
	})

	if !strings.Contains(output, "Print version information") {
		t.Errorf("Expected version command description in help output")
	}
}

func TestPrintUsageIncludesExamples(t *testing.T) {
	output := captureOutput(func() {
		printUsage()
	})

	examples := []string{
		"vibe serve",
		"vibe version",
		"vibe help",
	}

	for _, example := range examples {
		if !strings.Contains(output, example) {
			t.Errorf("Expected example '%s' in help output, got: %s", example, output)
		}
	}
}

func TestUsageTextFormat(t *testing.T) {
	if !strings.Contains(usageText, "VibeSQL") {
		t.Error("usageText should contain 'VibeSQL'")
	}

	if !strings.Contains(usageText, "serve") {
		t.Error("usageText should contain 'serve' command")
	}

	if !strings.Contains(usageText, "version") {
		t.Error("usageText should contain 'version' command")
	}

	if !strings.Contains(usageText, "help") {
		t.Error("usageText should contain 'help' command")
	}
}

func TestMainNoArgs(t *testing.T) {
	if os.Getenv("TEST_MAIN_NO_ARGS") == "1" {
		os.Args = []string{"vibe"}
		main()
		return
	}

	output := runMainWithArgsExpectError(t, []string{"vibe"}, "TEST_MAIN_NO_ARGS")
	
	if !strings.Contains(output, "Usage:") {
		t.Errorf("Expected help output when no args provided, got: %s", output)
	}
}

func TestMainHelpCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"help", []string{"vibe", "help"}},
		{"--help", []string{"vibe", "--help"}},
		{"-h", []string{"vibe", "-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("TEST_MAIN_HELP") == "1" {
				os.Args = tt.args
				main()
				return
			}

			output := runMainWithArgs(t, tt.args, "TEST_MAIN_HELP")
			
			if !strings.Contains(output, "Usage:") {
				t.Errorf("Expected help output for %v, got: %s", tt.args, output)
			}
		})
	}
}

func TestMainVersionCommand(t *testing.T) {
	if os.Getenv("TEST_MAIN_VERSION") == "1" {
		os.Args = []string{"vibe", "version"}
		main()
		return
	}

	output := runMainWithArgs(t, []string{"vibe", "version"}, "TEST_MAIN_VERSION")
	
	if !strings.Contains(output, "VibeSQL Version Information") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

func TestMainUnknownCommand(t *testing.T) {
	if os.Getenv("TEST_MAIN_UNKNOWN") == "1" {
		os.Args = []string{"vibe", "unknown"}
		main()
		return
	}

	output := runMainWithArgsExpectError(t, []string{"vibe", "unknown"}, "TEST_MAIN_UNKNOWN")
	
	if !strings.Contains(output, "Unknown command") {
		t.Errorf("Expected 'Unknown command' error, got: %s", output)
	}
	
	if !strings.Contains(output, "unknown") {
		t.Errorf("Expected command name in error, got: %s", output)
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func runMainWithArgs(t *testing.T, args []string, envVar string) string {
	t.Helper()
	
	cmd := os.Args[0]
	
	testCmd := exec.Command(cmd, "-test.run=^"+t.Name()+"$")
	testCmd.Env = append(os.Environ(), envVar+"=1")
	
	output, err := testCmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				return string(output)
			}
		}
		t.Fatalf("Failed to run command: %v", err)
	}
	
	return string(output)
}

func runMainWithArgsExpectError(t *testing.T, args []string, envVar string) string {
	t.Helper()
	
	cmd := os.Args[0]
	
	testCmd := exec.Command(cmd, "-test.run=^"+t.Name()+"$")
	testCmd.Env = append(os.Environ(), envVar+"=1")
	
	output, err := testCmd.CombinedOutput()
	if err == nil {
		t.Fatal("Expected command to fail but it succeeded")
	}
	
	return string(output)
}
