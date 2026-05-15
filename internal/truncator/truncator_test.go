package truncator

import (
	"strings"
	"testing"
)

func TestTruncate_ShortValuesUnchanged(t *testing.T) {
	env := map[string]string{"KEY": "hello"}
	out := Truncate(env, DefaultOptions())
	if out["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["KEY"])
	}
}

func TestTruncate_LongValueTruncated(t *testing.T) {
	long := strings.Repeat("a", 100)
	env := map[string]string{"KEY": long}
	opts := DefaultOptions() // maxLen=64, suffix="..."
	out := Truncate(env, opts)
	got := out["KEY"]
	if len([]rune(got)) != 64 {
		t.Errorf("expected length 64, got %d", len([]rune(got)))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestTruncate_CustomSuffix(t *testing.T) {
	env := map[string]string{"KEY": strings.Repeat("b", 20)}
	opts := Options{MaxLen: 10, Suffix: "~~"}
	out := Truncate(env, opts)
	got := out["KEY"]
	if len([]rune(got)) != 10 {
		t.Errorf("expected length 10, got %d", len([]rune(got)))
	}
	if !strings.HasSuffix(got, "~~") {
		t.Errorf("expected suffix '~~'", )
	}
}

func TestTruncate_KeysOnly_OtherKeysUntouched(t *testing.T) {
	long := strings.Repeat("x", 80)
	env := map[string]string{
		"TARGET": long,
		"OTHER":  long,
	}
	opts := Options{MaxLen: 20, Suffix: "...", KeysOnly: []string{"TARGET"}}
	out := Truncate(env, opts)
	if len([]rune(out["TARGET"])) != 20 {
		t.Errorf("TARGET should be truncated to 20, got %d", len([]rune(out["TARGET"])))
	}
	if out["OTHER"] != long {
		t.Errorf("OTHER should be unchanged")
	}
}

func TestTruncate_ExactLengthUnchanged(t *testing.T) {
	v := strings.Repeat("z", 64)
	env := map[string]string{"KEY": v}
	out := Truncate(env, DefaultOptions())
	if out["KEY"] != v {
		t.Errorf("value at exact maxLen should not be truncated")
	}
}

func TestTruncate_EmptyMap(t *testing.T) {
	out := Truncate(map[string]string{}, DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.MaxLen != 64 {
		t.Errorf("expected MaxLen 64, got %d", opts.MaxLen)
	}
	if opts.Suffix != "..." {
		t.Errorf("expected Suffix '...', got %q", opts.Suffix)
	}
	if opts.KeysOnly != nil {
		t.Errorf("expected KeysOnly nil")
	}
}
