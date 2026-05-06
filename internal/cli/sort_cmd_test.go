package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSortTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestRunSort_Alphabetical(t *testing.T) {
	path := writeSortTempEnv(t, "ZEBRA=1\nAPPLE=2\nMIDDLE=3\n")
	args := []string{"sort", path}
	out, err := captureOutput(func() error {
		return RunSort(args)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APPLE") {
		t.Errorf("expected first line to start with APPLE, got %s", lines[0])
	}
}

func TestRunSort_GroupByPrefix(t *testing.T) {
	path := writeSortTempEnv(t, "DB_HOST=localhost\nAPP_NAME=myapp\nDB_PORT=5432\nAPP_ENV=prod\n")
	args := []string{"sort", "--group", path}
	out, err := captureOutput(func() error {
		return RunSort(args)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_") || !strings.Contains(out, "APP_") {
		t.Errorf("expected grouped output, got: %s", out)
	}
}

func TestRunSort_WriteToFile(t *testing.T) {
	path := writeSortTempEnv(t, "ZEBRA=1\nAPPLE=2\n")
	outPath := filepath.Join(t.TempDir(), "sorted.env")
	args := []string{"sort", "--output", outPath, path}
	if err := RunSort(args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !strings.Contains(string(data), "APPLE") {
		t.Errorf("expected output file to contain APPLE, got: %s", string(data))
	}
}

func TestRunSort_MissingFile(t *testing.T) {
	args := []string{"sort", "/nonexistent/.env"}
	if err := RunSort(args); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRunSort_NoArgs(t *testing.T) {
	args := []string{"sort"}
	if err := RunSort(args); err == nil {
		t.Error("expected error when no file provided, got nil")
	}
}
