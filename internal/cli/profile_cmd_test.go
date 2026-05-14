package cli

import (
	"os"
	"strings"
	"testing"
)

func writeProfileTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRunProfile_TextOutput(t *testing.T) {
	file := writeProfileTempEnv(t, "PORT=8080\nDEBUG=true\nSECRET=\nAPI_URL=https://api.example.com\n")
	out, err := captureOutput(func() error {
		return RunProfile([]string{file}, nil)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Total keys") {
		t.Errorf("expected Total keys in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Density") {
		t.Errorf("expected Density in output")
	}
	if !strings.Contains(out, "bool") {
		t.Errorf("expected type breakdown with bool")
	}
}

func TestRunProfile_JSONOutput(t *testing.T) {
	file := writeProfileTempEnv(t, "PORT=9000\nENABLED=false\n")
	out, err := captureOutput(func() error {
		return RunProfile([]string{"--format", "json", file}, nil)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "TotalKeys") {
		t.Errorf("expected JSON with TotalKeys, got:\n%s", out)
	}
	if !strings.Contains(out, "TypeBreakdown") {
		t.Errorf("expected TypeBreakdown in JSON")
	}
}

func TestRunProfile_MissingFile(t *testing.T) {
	err := RunProfile([]string{"/nonexistent/path.env"}, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunProfile_NoArgs(t *testing.T) {
	err := RunProfile([]string{}, nil)
	if err == nil {
		t.Fatal("expected error when no file provided")
	}
}

func TestRunProfile_EmptyFile(t *testing.T) {
	file := writeProfileTempEnv(t, "")
	out, err := captureOutput(func() error {
		return RunProfile([]string{file}, nil)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Total keys   : 0") {
		t.Errorf("expected 0 total keys, got:\n%s", out)
	}
}
