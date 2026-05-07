package patcher_test

import (
	"testing"

	"github.com/user/envdiff/internal/patcher"
)

func base() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"APP_PORT": "8080",
		"DB_PASS":  "secret",
	}
}

func TestApply_SetNewKey(t *testing.T) {
	out, err := patcher.Apply(base(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "NEW_KEY", Value: "hello"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
	if len(out) != 4 {
		t.Errorf("expected 4 keys, got %d", len(out))
	}
}

func TestApply_UpdateExistingKey(t *testing.T) {
	out, err := patcher.Apply(base(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_ENV", Value: "staging"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", out["APP_ENV"])
	}
}

func TestApply_DeleteKey(t *testing.T) {
	out, err := patcher.Apply(base(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "DB_PASS"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be deleted")
	}
}

func TestApply_DeleteMissingKey_NoError(t *testing.T) {
	_, err := patcher.Apply(base(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "GHOST"},
	}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestApply_DeleteMissingKey_ErrorMode(t *testing.T) {
	opts := patcher.DefaultOptions()
	opts.ErrorOnMissingDelete = true
	_, err := patcher.Apply(base(), []patcher.Patch{
		{Op: patcher.OpDelete, Key: "GHOST"},
	}, opts)
	if err == nil {
		t.Fatal("expected error for missing delete, got nil")
	}
}

func TestApply_ErrorOnNoChange(t *testing.T) {
	opts := patcher.DefaultOptions()
	opts.ErrorOnNoChange = true
	_, err := patcher.Apply(base(), []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_ENV", Value: "production"},
	}, opts)
	if err == nil {
		t.Fatal("expected error for no-change set, got nil")
	}
}

func TestApply_UnknownOp(t *testing.T) {
	_, err := patcher.Apply(base(), []patcher.Patch{
		{Op: "rename", Key: "APP_ENV"},
	}, patcher.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestApply_DoesNotMutateSrc(t *testing.T) {
	src := base()
	_, _ = patcher.Apply(src, []patcher.Patch{
		{Op: patcher.OpSet, Key: "APP_ENV", Value: "dev"},
	}, patcher.DefaultOptions())
	if src["APP_ENV"] != "production" {
		t.Error("Apply must not mutate the source map")
	}
}
