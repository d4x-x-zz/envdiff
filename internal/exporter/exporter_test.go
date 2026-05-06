package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/exporter"
)

func makeResult(missing, mismatched []string) *differ.Result {
	r := &differ.Result{}
	for _, k := range missing {
		r.Missing = append(r.Missing, differ.Diff{Key: k, Side: "right"})
	}
	for _, k := range mismatched {
		r.Mismatched = append(r.Mismatched, differ.Diff{Key: k, LeftVal: "a", RightVal: "b"})
	}
	return r
}

func TestExport_DotEnv_Redacted(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := makeResult(nil, nil)
	tmp := filepath.Join(t.TempDir(), "out.env")

	if err := exporter.Export(r, env, tmp, exporter.Options{Format: exporter.FormatDotEnv}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(tmp)
	content := string(b)
	if !strings.Contains(content, "BAZ=\n") || !strings.Contains(content, "FOO=\n") {
		t.Errorf("expected redacted keys, got:\n%s", content)
	}
	if strings.Contains(content, "bar") || strings.Contains(content, "qux") {
		t.Errorf("values should be redacted, got:\n%s", content)
	}
}

func TestExport_DotEnv_WithValues(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	r := makeResult(nil, nil)
	tmp := filepath.Join(t.TempDir(), "out.env")

	if err := exporter.Export(r, env, tmp, exporter.Options{Format: exporter.FormatDotEnv, IncludeVals: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(tmp)
	if !strings.Contains(string(b), "KEY=val") {
		t.Errorf("expected KEY=val, got: %s", string(b))
	}
}

func TestExport_JSON_Redacted(t *testing.T) {
	env := map[string]string{"ALPHA": "secret"}
	r := makeResult(nil, nil)
	tmp := filepath.Join(t.TempDir(), "out.json")

	if err := exporter.Export(r, env, tmp, exporter.Options{Format: exporter.FormatJSON}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(tmp)
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if v, ok := m["ALPHA"]; !ok || v != "" {
		t.Errorf("expected ALPHA with empty value, got %q", v)
	}
}

func TestExport_Markdown(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	r := makeResult([]string{"DB_PASS"}, nil)
	tmp := filepath.Join(t.TempDir(), "out.md")

	if err := exporter.Export(r, env, tmp, exporter.Options{Format: exporter.FormatMarkdown}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(tmp)
	content := string(b)
	if !strings.Contains(content, "| Key | Value |") {
		t.Errorf("missing markdown header, got:\n%s", content)
	}
	if !strings.Contains(content, "DB_HOST") || !strings.Contains(content, "DB_PASS") {
		t.Errorf("expected both keys in markdown, got:\n%s", content)
	}
}

func TestExport_MissingKeysIncluded(t *testing.T) {
	env := map[string]string{"EXISTING": "val"}
	r := makeResult([]string{"MISSING_KEY"}, nil)
	tmp := filepath.Join(t.TempDir(), "out.env")

	if err := exporter.Export(r, env, tmp, exporter.Options{Format: exporter.FormatDotEnv}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b, _ := os.ReadFile(tmp)
	if !strings.Contains(string(b), "MISSING_KEY=") {
		t.Errorf("expected MISSING_KEY in output, got:\n%s", string(b))
	}
}
