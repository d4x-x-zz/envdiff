package renamer_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/renamer"
)

func baseEnv() map[string]string {
	return map[string]string{
		"OLD_KEY":    "value1",
		"ANOTHER":    "value2",
		"KEEP_ME":    "value3",
	}
}

func TestRename_BasicRename(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.Mapping = map[string]string{"OLD_KEY": "NEW_KEY"}

	out, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", out["NEW_KEY"])
	}
}

func TestRename_IgnoreMissing_True(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.Mapping = map[string]string{"DOES_NOT_EXIST": "SOMETHING"}
	opts.IgnoreMissing = true

	_, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("expected no error with IgnoreMissing=true, got: %v", err)
	}
}

func TestRename_IgnoreMissing_False(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.Mapping = map[string]string{"DOES_NOT_EXIST": "SOMETHING"}
	opts.IgnoreMissing = false

	_, err := renamer.Rename(baseEnv(), opts)
	if err == nil {
		t.Fatal("expected error when source key is missing and IgnoreMissing=false")
	}
}

func TestRename_ConflictSkip(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.Mapping = map[string]string{"OLD_KEY": "ANOTHER"}
	opts.ConflictStrategy = renamer.SkipOnConflict

	out, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// ANOTHER should retain its original value
	if out["ANOTHER"] != "value2" {
		t.Errorf("expected ANOTHER=value2, got %q", out["ANOTHER"])
	}
}

func TestRename_ConflictOverwrite(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.Mapping = map[string]string{"OLD_KEY": "ANOTHER"}
	opts.ConflictStrategy = renamer.OverwriteOnConflict

	out, err := renamer.Rename(baseEnv(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ANOTHER"] != "value1" {
		t.Errorf("expected ANOTHER to be overwritten with value1, got %q", out["ANOTHER"])
	}
}

func TestRename_ConflictError(t *testing.T) {
	opts := renamer.DefaultOptions()
	opts.Mapping = map[string]string{"OLD_KEY": "ANOTHER"}
	opts.ConflictStrategy = renamer.ErrorOnConflict

	_, err := renamer.Rename(baseEnv(), opts)
	if err == nil {
		t.Fatal("expected error on conflict with ErrorOnConflict strategy")
	}
}

func TestRename_OriginalUnmodified(t *testing.T) {
	env := baseEnv()
	opts := renamer.DefaultOptions()
	opts.Mapping = map[string]string{"OLD_KEY": "NEW_KEY"}

	_, _ = renamer.Rename(env, opts)
	if _, ok := env["OLD_KEY"]; !ok {
		t.Error("original map should not be modified")
	}
}
