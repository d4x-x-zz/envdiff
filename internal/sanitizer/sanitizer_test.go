package sanitizer

import (
	"testing"
)

func ptr(s string) *string { return &s }

func TestSanitize_TrimWhitespace(t *testing.T) {
	env := map[string]string{"KEY": "  hello world  "}
	out := Sanitize(env, Options{TrimWhitespace: true})
	if got := out["KEY"]; got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestSanitize_NoTrim(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	out := Sanitize(env, Options{TrimWhitespace: false})
	if got := out["KEY"]; got != "  value  " {
		t.Errorf("expected '  value  ', got %q", got)
	}
}

func TestSanitize_ReplaceNewlines(t *testing.T) {
	env := map[string]string{"KEY": "line1\nline2"}
	out := Sanitize(env, Options{ReplaceNewlines: ptr(`\n`)})
	if got := out["KEY"]; got != `line1\nline2` {
		t.Errorf("unexpected value: %q", got)
	}
}

func TestSanitize_RemoveNewlines(t *testing.T) {
	empty := ""
	env := map[string]string{"KEY": "line1\nline2"}
	out := Sanitize(env, Options{ReplaceNewlines: &empty})
	if got := out["KEY"]; got != "line1line2" {
		t.Errorf("expected 'line1line2', got %q", got)
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	env := map[string]string{"KEY": "val\x01ue\x1f"}
	out := Sanitize(env, Options{StripControlChars: true})
	if got := out["KEY"]; got != "value" {
		t.Errorf("expected 'value', got %q", got)
	}
}

func TestSanitize_TabPreservedWhenStrippingControls(t *testing.T) {
	env := map[string]string{"KEY": "col1\tcol2"}
	out := Sanitize(env, Options{StripControlChars: true})
	if got := out["KEY"]; got != "col1\tcol2" {
		t.Errorf("expected tab to be preserved, got %q", got)
	}
}

func TestSanitize_MaxLength(t *testing.T) {
	env := map[string]string{"KEY": "abcdefghij"}
	out := Sanitize(env, Options{MaxLength: 5})
	if got := out["KEY"]; got != "abcde" {
		t.Errorf("expected 'abcde', got %q", got)
	}
}

func TestSanitize_MaxLength_Zero_NoTruncate(t *testing.T) {
	env := map[string]string{"KEY": "abcdefghij"}
	out := Sanitize(env, Options{MaxLength: 0})
	if got := out["KEY"]; got != "abcdefghij" {
		t.Errorf("expected full value, got %q", got)
	}
}

func TestSanitize_DefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	env := map[string]string{"KEY": "  hello\nworld\x01  "}
	out := Sanitize(env, opts)
	if got := out["KEY"]; got != `hello\nworld` {
		t.Errorf("unexpected result with default options: %q", got)
	}
}

func TestSanitize_EmptyMap(t *testing.T) {
	out := Sanitize(map[string]string{}, DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
