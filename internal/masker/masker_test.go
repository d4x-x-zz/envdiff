package masker_test

import (
	"testing"

	"github.com/user/envdiff/internal/masker"
)

func TestMask_DefaultOptions(t *testing.T) {
	env := map[string]string{"API_KEY": "abcdef1234"}
	out := masker.Mask(env, masker.DefaultOptions())
	got := out["API_KEY"]
	// prefix=2 visible, rest masked
	if got[:2] != "ab" {
		t.Errorf("expected prefix 'ab', got %q", got[:2])
	}
	if len(got) != len("abcdef1234") {
		t.Errorf("length mismatch: got %d", len(got))
	}
}

func TestMask_VisibleSuffix(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.VisiblePrefix = 2
	opts.VisibleSuffix = 2
	env := map[string]string{"TOKEN": "abcdef"}
	out := masker.Mask(env, opts)
	got := out["TOKEN"]
	// expect "ab**ef"
	if got != "ab**ef" {
		t.Errorf("expected 'ab**ef', got %q", got)
	}
}

func TestMask_ShortValueFullyMasked(t *testing.T) {
	opts := masker.DefaultOptions() // MinLength=4
	env := map[string]string{"PIN": "abc"}
	out := masker.Mask(env, opts)
	if out["PIN"] != "***" {
		t.Errorf("expected '***', got %q", out["PIN"])
	}
}

func TestMask_EmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY": ""}
	out := masker.Mask(env, masker.DefaultOptions())
	if out["EMPTY"] != "" {
		t.Errorf("expected empty string, got %q", out["EMPTY"])
	}
}

func TestMask_CustomMaskChar(t *testing.T) {
	opts := masker.DefaultOptions()
	opts.MaskChar = "#"
	opts.VisiblePrefix = 1
	env := map[string]string{"KEY": "secret"}
	out := masker.Mask(env, opts)
	got := out["KEY"]
	if got[0] != 's' {
		t.Errorf("expected leading 's', got %q", string(got[0]))
	}
	for _, ch := range got[1:] {
		if ch != '#' {
			t.Errorf("expected '#' mask char, got %q", string(ch))
		}
	}
}

func TestMask_KeysUnchanged(t *testing.T) {
	env := map[string]string{"DB_PASS": "hunter2", "PORT": "5432"}
	out := masker.Mask(env, masker.DefaultOptions())
	for k := range env {
		if _, ok := out[k]; !ok {
			t.Errorf("key %q missing from output", k)
		}
	}
}
