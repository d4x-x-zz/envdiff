package aliaser

import (
	"testing"
)

var base = map[string]string{
	"DB_HOST": "localhost",
	"DB_PORT": "5432",
	"APP_ENV":  "production",
}

func TestAlias_BasicAlias(t *testing.T) {
	opts := DefaultOptions()
	opts.Aliases = map[string][]string{
		"DB_HOST": {"DATABASE_HOST"},
	}
	out, err := Alias(base, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
	// original key still present
	if out["DB_HOST"] != "localhost" {
		t.Errorf("original key DB_HOST should still be present")
	}
}

func TestAlias_MultipleDestinations(t *testing.T) {
	opts := DefaultOptions()
	opts.Aliases = map[string][]string{
		"DB_PORT": {"DATABASE_PORT", "POSTGRES_PORT"},
	}
	out, err := Alias(base, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, dest := range []string{"DATABASE_PORT", "POSTGRES_PORT"} {
		if out[dest] != "5432" {
			t.Errorf("expected %s=5432, got %q", dest, out[dest])
		}
	}
}

func TestAlias_IgnoreMissing_True(t *testing.T) {
	opts := DefaultOptions()
	opts.IgnoreMissing = true
	opts.Aliases = map[string][]string{
		"NONEXISTENT": {"SOME_DEST"},
	}
	_, err := Alias(base, opts)
	if err != nil {
		t.Fatalf("expected no error with IgnoreMissing=true, got: %v", err)
	}
}

func TestAlias_IgnoreMissing_False(t *testing.T) {
	opts := DefaultOptions()
	opts.IgnoreMissing = false
	opts.Aliases = map[string][]string{
		"NONEXISTENT": {"SOME_DEST"},
	}
	_, err := Alias(base, opts)
	if err == nil {
		t.Fatal("expected error with IgnoreMissing=false, got nil")
	}
}

func TestAlias_OverwriteExisting_False(t *testing.T) {
	env := map[string]string{
		"SRC":  "new_value",
		"DEST": "original_value",
	}
	opts := DefaultOptions()
	opts.OverwriteExisting = false
	opts.Aliases = map[string][]string{
		"SRC": {"DEST"},
	}
	out, err := Alias(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DEST"] != "original_value" {
		t.Errorf("expected DEST to remain original_value, got %q", out["DEST"])
	}
}

func TestAlias_OverwriteExisting_True(t *testing.T) {
	env := map[string]string{
		"SRC":  "new_value",
		"DEST": "original_value",
	}
	opts := DefaultOptions()
	opts.OverwriteExisting = true
	opts.Aliases = map[string][]string{
		"SRC": {"DEST"},
	}
	out, err := Alias(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DEST"] != "new_value" {
		t.Errorf("expected DEST=new_value, got %q", out["DEST"])
	}
}

func TestAlias_EmptyMap(t *testing.T) {
	opts := DefaultOptions()
	opts.Aliases = map[string][]string{
		"DB_HOST": {"DATABASE_HOST"},
	}
	out, err := Alias(map[string]string{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output map, got %d keys", len(out))
	}
}
