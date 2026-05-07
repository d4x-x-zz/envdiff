package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemplateTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestRunTemplate_BasicOutput(t *testing.T) {
	p := writeTemplateTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nSECRET_KEY=abc123\n")
	out, err := captureOutput(func() error {
		return RunTemplate([]string{p})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
}

func TestRunTemplate_NoTypedPlaceholders(t *testing.T) {
	p := writeTemplateTempEnv(t, "APP_NAME=myapp\nAPP_ENV=production\n")
	out, err := captureOutput(func() error {
		return RunTemplate([]string{"--no-types", p})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME in output, got: %s", out)
	}
}

func TestRunTemplate_CommentOriginal(t *testing.T) {
	p := writeTemplateTempEnv(t, "HOST=localhost\nPORT=8080\n")
	out, err := captureOutput(func() error {
		return RunTemplate([]string{"--comment", p})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "#") {
		t.Errorf("expected comments in output, got: %s", out)
	}
}

func TestRunTemplate_WriteToFile(t *testing.T) {
	p := writeTemplateTempEnv(t, "DB_URL=postgres://localhost/db\n")
	out := filepath.Join(t.TempDir(), ".env.template")
	_, err := captureOutput(func() error {
		return RunTemplate([]string{"--output", out, p})
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "DB_URL") {
		t.Errorf("expected DB_URL in output file, got: %s", string(data))
	}
}

func TestRunTemplate_MissingFile(t *testing.T) {
	err := RunTemplate([]string{"/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
