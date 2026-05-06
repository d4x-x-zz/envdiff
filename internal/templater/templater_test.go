package templater_test

import (
	"strings"
	"testing"

	"envdiff/internal/templater"
)

func TestGenerate_TypedPlaceholders(t *testing.T) {
	env := map[string]string{
		"APP_NAME":  "myapp",
		"APP_PORT":  "8080",
		"DEBUG":     "true",
		"RATIO":     "3.14",
		"EMPTY_KEY": "",
	}
	opts := templater.DefaultOptions()
	out := templater.Generate(env, opts)

	expect := map[string]string{
		"APP_NAME":  "<string>",
		"APP_PORT":  "<number>",
		"DEBUG":     "<bool>",
		"RATIO":     "<number>",
		"EMPTY_KEY": "<string>",
	}
	for key, placeholder := range expect {
		line := key + "=" + placeholder
		if !strings.Contains(out, line) {
			t.Errorf("expected line %q in output, got:\n%s", line, out)
		}
	}
}

func TestGenerate_NoTypedPlaceholders(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	opts := templater.Options{UseTypedPlaceholders: false, CommentOriginal: false}
	out := templater.Generate(env, opts)

	for _, key := range []string{"DB_HOST", "DB_PORT"} {
		line := key + "="
		if !strings.Contains(out, line) {
			t.Errorf("expected bare key line %q in output", line)
		}
	}
	if strings.Contains(out, "<string>") || strings.Contains(out, "<number>") {
		t.Error("expected no typed placeholders when disabled")
	}
}

func TestGenerate_CommentOriginal(t *testing.T) {
	env := map[string]string{
		"SECRET": "hunter2",
	}
	opts := templater.Options{UseTypedPlaceholders: false, CommentOriginal: true}
	out := templater.Generate(env, opts)

	if !strings.Contains(out, "# original: hunter2") {
		t.Errorf("expected original value comment, got:\n%s", out)
	}
}

func TestGenerate_EmptyMap(t *testing.T) {
	out := templater.Generate(map[string]string{}, templater.DefaultOptions())
	if out != "" {
		t.Errorf("expected empty output for empty map, got: %q", out)
	}
}

func TestGenerate_SortedOutput(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}
	out := templater.Generate(env, templater.DefaultOptions())
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to be A_KEY, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to be Z_KEY, got %q", lines[2])
	}
}
