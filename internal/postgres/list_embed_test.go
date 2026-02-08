package postgres

import (
	"testing"
)

func TestListEmbedded(t *testing.T) {
	entries, err := embeddedPostgres.ReadDir("embed")
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}
	t.Logf("Found %d embedded files:", len(entries))
	for _, e := range entries {
		info, _ := e.Info()
		t.Logf("  %s - %d bytes", e.Name(), info.Size())
	}
}
