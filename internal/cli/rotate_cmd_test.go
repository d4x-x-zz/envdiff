package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeRotateTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeRotateTempEnv: %v", err)
	}
	return p
}

func TestRunRotate_BasicRename(t *testing.T) {
	f := writeRotateTempEnv(t, "DB_HOST=localhost\nDB_PASS=secret\n")
	out, err := captureOutput(func(stdout, stderr *os.File) error {
		return cli.RunRotate([]string{"--file", f, "--renames", "DB_HOST:DATABASE_HOST"}, stdout, stderr)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DATABASE_HOST=localhost") {
		t.Errorf("expected DATABASE_HOST in output, got:\n%s", out)
	}
	if strings.Contains(out, "DB_HOST=") {
		t.Errorf("expected DB_HOST removed, got:\n%s", out)
	}
}

func TestRunRotate_StrategyDeprecate(t *testing.T) {
	f := writeRotateTempEnv(t, "OLD_KEY=value\n")
	out, err := captureOutput(func(stdout, stderr *os.File) error {
		return cli.RunRotate([]string{"--file", f, "--renames", "OLD_KEY:NEW_KEY", "--strategy", "deprecate"}, stdout, stderr)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "NEW_KEY=value") {
		t.Errorf("expected NEW_KEY=value, got:\n%s", out)
	}
	if !strings.Contains(out, "OLD_KEY=") {
		t.Errorf("expected OLD_KEY= (deprecated), got:\n%s", out)
	}
}

func TestRunRotate_WriteToFile(t *testing.T) {
	f := writeRotateTempEnv(t, "FOO=bar\n")
	outFile := filepath.Join(t.TempDir(), "out.env")
	_, err := captureOutput(func(stdout, stderr *os.File) error {
		return cli.RunRotate([]string{"--file", f, "--renames", "FOO:BAR", "--output", outFile}, stdout, stderr)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "BAR=bar") {
		t.Errorf("expected BAR=bar in file, got:\n%s", data)
	}
}

func TestRunRotate_MissingFile(t *testing.T) {
	_, err := captureOutput(func(stdout, stderr *os.File) error {
		return cli.RunRotate([]string{"--file", "/nonexistent/.env", "--renames", "A:B"}, stdout, stderr)
	})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRunRotate_InvalidStrategy(t *testing.T) {
	f := writeRotateTempEnv(t, "A=1\n")
	_, err := captureOutput(func(stdout, stderr *os.File) error {
		return cli.RunRotate([]string{"--file", f, "--renames", "A:B", "--strategy", "bogus"}, stdout, stderr)
	})
	if err == nil {
		t.Error("expected error for invalid strategy")
	}
}
