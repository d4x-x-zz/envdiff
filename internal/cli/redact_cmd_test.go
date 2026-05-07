package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeRedactTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func TestRunRedact_DefaultMasksSensitive(t *testing.T) {
	p := writeRedactTempEnv(t, "SECRET_KEY=supersecret\nAPP_NAME=myapp\n")
	out, err := captureOutput(func() error {
		return redactMain([]string{p}, os.Stdout)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "supersecret") {
		t.Errorf("expected secret value to be masked, got: %s", out)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME in output, got: %s", out)
	}
}

func TestRunRedact_CustomMask(t *testing.T) {
	p := writeRedactTempEnv(t, "API_TOKEN=abc123\n")
	out, err := captureOutput(func() error {
		return redactMain([]string{"--mask", "***", p}, os.Stdout)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected custom mask in output, got: %s", out)
	}
}

func TestRunRedact_ExplicitKeys(t *testing.T) {
	p := writeRedactTempEnv(t, "MY_CUSTOM=value\nOTHER=hello\n")
	out, err := captureOutput(func() error {
		return redactMain([]string{"--keys", "MY_CUSTOM", p}, os.Stdout)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "value") {
		t.Errorf("expected MY_CUSTOM value to be masked, got: %s", out)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected OTHER value to remain, got: %s", out)
	}
}

func TestRunRedact_WriteToFile(t *testing.T) {
	p := writeRedactTempEnv(t, "PASSWORD=secret\n")
	out := filepath.Join(t.TempDir(), ".env.redacted")
	_, err := captureOutput(func() error {
		return redactMain([]string{"--output", out, p}, os.Stdout)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if strings.Contains(string(data), "secret") {
		t.Errorf("expected redacted output file, got: %s", string(data))
	}
}

func TestRunRedact_MissingFile(t *testing.T) {
	err := redactMain([]string{"/no/such/.env"}, os.Stdout)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
