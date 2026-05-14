package deduper

import (
	"testing"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	r := Dedupe([]map[string]string{a, b}, DefaultOptions())

	if len(r.Env) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Env))
	}
	if len(r.Removed) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(r.Removed))
	}
}

func TestDedupe_SameKeySameValue(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar"}

	r := Dedupe([]map[string]string{a, b}, DefaultOptions())

	if len(r.Env) != 1 {
		t.Fatalf("expected 1 key, got %d", len(r.Env))
	}
	if len(r.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(r.Removed))
	}
	if r.Removed[0].SourceIndex != 1 {
		t.Errorf("expected source index 1, got %d", r.Removed[0].SourceIndex)
	}
}

func TestDedupe_SameKeyDifferentValue_KeepFirst(t *testing.T) {
	a := map[string]string{"FOO": "first"}
	b := map[string]string{"FOO": "second"}

	opts := DefaultOptions() // KeepFirst = true
	r := Dedupe([]map[string]string{a, b}, opts)

	if r.Env["FOO"] != "first" {
		t.Errorf("expected 'first', got %q", r.Env["FOO"])
	}
	if len(r.Removed) != 0 {
		t.Errorf("different value should not be counted as duplicate")
	}
}

func TestDedupe_SameKeyDifferentValue_KeepLast(t *testing.T) {
	a := map[string]string{"FOO": "first"}
	b := map[string]string{"FOO": "second"}

	opts := Options{KeepFirst: false}
	r := Dedupe([]map[string]string{a, b}, opts)

	if r.Env["FOO"] != "second" {
		t.Errorf("expected 'second', got %q", r.Env["FOO"])
	}
}

func TestDedupe_SkipValueCheck(t *testing.T) {
	a := map[string]string{"FOO": "alpha"}
	b := map[string]string{"FOO": "beta"}

	opts := Options{KeepFirst: true, SkipValueCheck: true}
	r := Dedupe([]map[string]string{a, b}, opts)

	if len(r.Removed) != 1 {
		t.Fatalf("expected 1 removed with SkipValueCheck, got %d", len(r.Removed))
	}
	if r.Env["FOO"] != "alpha" {
		t.Errorf("expected original value 'alpha', got %q", r.Env["FOO"])
	}
}

func TestDedupe_EmptyInput(t *testing.T) {
	r := Dedupe([]map[string]string{}, DefaultOptions())
	if len(r.Env) != 0 || len(r.Removed) != 0 {
		t.Error("expected empty result for empty input")
	}
}
