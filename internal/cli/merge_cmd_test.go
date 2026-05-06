package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/cli"
)

func writeMergeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeMergeTempEnv: %v", err)
	}
	return p
}

func TestRunMerge_TwoFiles(t *testing.T) {
	f1 := writeMergeTempEnv(t, "A=1\nB=2\n")
	f2 := writeMergeTempEnv(t, "C=3\n")

	err := cli.RunMerge([]string{f1, f2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunMerge_InsufficientFiles(t *testing.T) {
	f1 := writeMergeTempEnv(t, "A=1\n")
	err := cli.RunMerge([]string{f1})
	if err == nil {
		t.Fatal("expected error for single file, got nil")
	}
}

func TestRunMerge_InvalidStrategy(t *testing.T) {
	f1 := writeMergeTempEnv(t, "A=1\n")
	f2 := writeMergeTempEnv(t, "B=2\n")
	err := cli.RunMerge([]string{"--strategy=bogus", f1, f2})
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestRunMerge_WriteToFile(t *testing.T) {
	f1 := writeMergeTempEnv(t, "A=1\n")
	f2 := writeMergeTempEnv(t, "B=2\n")
	out := filepath.Join(t.TempDir(), "merged.env")

	err := cli.RunMerge([]string{"--out=" + out, f1, f2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "A=1") || !strings.Contains(content, "B=2") {
		t.Errorf("unexpected output content: %s", content)
	}
}

func TestRunMerge_ConflictError(t *testing.T) {
	f1 := writeMergeTempEnv(t, "KEY=alpha\n")
	f2 := writeMergeTempEnv(t, "KEY=beta\n")

	err := cli.RunMerge([]string{"--strategy=error", f1, f2})
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}
