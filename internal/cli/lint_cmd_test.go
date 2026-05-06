package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeLintTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestRunLint_NoIssues(t *testing.T) {
	f := writeLintTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	out, err := captureOutput(func(o *os.File) error {
		return RunLint([]string{f}, o)
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "no lint issues found") {
		t.Errorf("expected clean message, got: %s", out)
	}
}

func TestRunLint_LowercaseKey(t *testing.T) {
	f := writeLintTempEnv(t, "app_host=localhost\n")
	out, err := captureOutput(func(o *os.File) error {
		return RunLint([]string{f}, o)
	})
	if err == nil {
		t.Fatal("expected error for lowercase key")
	}
	if !strings.Contains(out, "app_host") {
		t.Errorf("expected key in output, got: %s", out)
	}
}

func TestRunLint_AllowLowerFlag(t *testing.T) {
	f := writeLintTempEnv(t, "app_host=localhost\n")
	_, err := captureOutput(func(o *os.File) error {
		return RunLint([]string{"-allow-lower", f}, o)
	})
	if err != nil {
		t.Fatalf("expected no error with -allow-lower, got: %v", err)
	}
}

func TestRunLint_MissingFile(t *testing.T) {
	_, err := captureOutput(func(o *os.File) error {
		return RunLint([]string{"/nonexistent/.env"}, o)
	})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunLint_NoArgs(t *testing.T) {
	_, err := captureOutput(func(o *os.File) error {
		return RunLint([]string{}, o)
	})
	if err == nil {
		t.Fatal("expected error when no file provided")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage message in error, got: %v", err)
	}
}
