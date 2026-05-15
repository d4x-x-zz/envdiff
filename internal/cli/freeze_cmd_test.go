package cli

import (
	"os"
	"strings"
	"testing"
)

func writeFreezeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestRunFreeze_TextOutput(t *testing.T) {
	path := writeFreezeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\n")
	out, err := captureOutput(func(w *os.File) error {
		return RunFreeze([]string{path}, w)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "fingerprint:") {
		t.Error("expected fingerprint in output")
	}
	if !strings.Contains(out, "keys:        3") {
		t.Errorf("expected 3 keys, got:\n%s", out)
	}
}

func TestRunFreeze_JSONOutput(t *testing.T) {
	path := writeFreezeTempEnv(t, "FOO=1\nBAR=2\n")
	out, err := captureOutput(func(w *os.File) error {
		return RunFreeze([]string{"-format", "json", path}, w)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "fingerprint") {
		t.Error("expected JSON with fingerprint key")
	}
	if !strings.Contains(out, "key_count") {
		t.Error("expected JSON with key_count field")
	}
}

func TestRunFreeze_KeyPrefix(t *testing.T) {
	path := writeFreezeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\n")
	out, err := captureOutput(func(w *os.File) error {
		return RunFreeze([]string{"-prefix", "DB_", path}, w)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "keys:        2") {
		t.Errorf("expected 2 DB_ keys, got:\n%s", out)
	}
}

func TestRunFreeze_MissingFile(t *testing.T) {
	err := RunFreeze([]string{"/nonexistent/.env"}, os.Stdout)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunFreeze_NoArgs(t *testing.T) {
	err := RunFreeze([]string{}, os.Stdout)
	if err == nil {
		t.Error("expected error when no args provided")
	}
}
