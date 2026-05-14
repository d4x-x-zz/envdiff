package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeCompareTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeCompareTempEnv: %v", err)
	}
	return p
}

func TestRunCompare_NoDiff(t *testing.T) {
	a := writeCompareTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	b := writeCompareTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	out, err := captureOutput(func() error {
		return cli.RunCompare([]string{a, b}, os.Stdout)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = out // output captured via captureOutput
}

func TestRunCompare_MissingKey(t *testing.T) {
	a := writeCompareTempEnv(t, "DB_HOST=localhost\nSECRET=abc\n")
	b := writeCompareTempEnv(t, "DB_HOST=localhost\n")
	var sb strings.Builder
	if err := cli.RunCompare([]string{a, b}, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "MISSING") {
		t.Errorf("expected MISSING in output, got: %s", sb.String())
	}
}

func TestRunCompare_MismatchedValue(t *testing.T) {
	a := writeCompareTempEnv(t, "DB_HOST=localhost\n")
	b := writeCompareTempEnv(t, "DB_HOST=remotehost\n")
	var sb strings.Builder
	if err := cli.RunCompare([]string{a, b}, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "MISMATCH") {
		t.Errorf("expected MISMATCH in output, got: %s", sb.String())
	}
}

func TestRunCompare_JSONFormat(t *testing.T) {
	a := writeCompareTempEnv(t, "DB_HOST=localhost\n")
	b := writeCompareTempEnv(t, "DB_HOST=remotehost\n")
	var sb strings.Builder
	if err := cli.RunCompare([]string{"-format", "json", a, b}, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "mismatch") {
		t.Errorf("expected json with mismatch, got: %s", sb.String())
	}
}

func TestRunCompare_InsufficientFiles(t *testing.T) {
	a := writeCompareTempEnv(t, "KEY=val\n")
	var sb strings.Builder
	err := cli.RunCompare([]string{a}, &sb)
	if err == nil {
		t.Fatal("expected error for single file")
	}
}

func TestRunCompare_MissingFile(t *testing.T) {
	var sb strings.Builder
	err := cli.RunCompare([]string{"/nonexistent/a.env", "/nonexistent/b.env"}, &sb)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunCompare_ThreeFiles(t *testing.T) {
	a := writeCompareTempEnv(t, "KEY=val\nONLY_A=1\n")
	b := writeCompareTempEnv(t, "KEY=val\n")
	c := writeCompareTempEnv(t, "KEY=val\n")
	var sb strings.Builder
	if err := cli.RunCompare([]string{a, b, c}, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "MISSING") {
		t.Errorf("expected MISSING for ONLY_A, got: %s", sb.String())
	}
}
