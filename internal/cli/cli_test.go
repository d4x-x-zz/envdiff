package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_NoDiff(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	if err := Run([]string{a, b}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRun_StrictWithDiff(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\n")
	b := writeTempEnv(t, "FOO=different\n")
	err := Run([]string{"--strict", a, b})
	if err == nil {
		t.Fatal("expected error in strict mode with differences")
	}
}

func TestRun_JSONFormat(t *testing.T) {
	a := writeTempEnv(t, "KEY=val\n")
	b := writeTempEnv(t, "KEY=val\n")
	if err := Run([]string{"--format", "json", a, b}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_MissingFile(t *testing.T) {
	err := Run([]string{"/nonexistent/.env", "/also/missing/.env"})
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}

func TestRun_TooFewArgs(t *testing.T) {
	err := Run([]string{"only-one-file"})
	if err == nil {
		t.Fatal("expected error for missing second file argument")
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	a := writeTempEnv(t, "K=v\n")
	b := writeTempEnv(t, "K=v\n")
	err := Run([]string{"--format", "xml", a, b})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
