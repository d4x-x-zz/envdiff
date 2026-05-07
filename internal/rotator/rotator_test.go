package rotator_test

import (
	"testing"

	"github.com/user/envdiff/internal/rotator"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
		"APP_ENV": "production",
	}
}

func TestRotate_StrategyRemove(t *testing.T) {
	opts := rotator.DefaultOptions() // StrategyRemove
	rotations := []rotator.Rotation{{OldKey: "DB_HOST", NewKey: "DATABASE_HOST"}}
	out, err := rotator.Rotate(baseEnv(), rotations, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", out["DATABASE_HOST"])
	}
}

func TestRotate_StrategyDeprecate(t *testing.T) {
	opts := rotator.Options{Strategy: rotator.StrategyDeprecate}
	rotations := []rotator.Rotation{{OldKey: "DB_PASS", NewKey: "DATABASE_PASSWORD"}}
	out, err := rotator.Rotate(baseEnv(), rotations, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "" {
		t.Errorf("expected DB_PASS to be empty, got %q", out["DB_PASS"])
	}
	if out["DATABASE_PASSWORD"] != "secret" {
		t.Errorf("expected DATABASE_PASSWORD=secret, got %q", out["DATABASE_PASSWORD"])
	}
}

func TestRotate_StrategyKeep(t *testing.T) {
	opts := rotator.Options{Strategy: rotator.StrategyKeep}
	rotations := []rotator.Rotation{{OldKey: "APP_ENV", NewKey: "APP_ENVIRONMENT"}}
	out, err := rotator.Rotate(baseEnv(), rotations, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV preserved, got %q", out["APP_ENV"])
	}
	if out["APP_ENVIRONMENT"] != "production" {
		t.Errorf("expected APP_ENVIRONMENT=production, got %q", out["APP_ENVIRONMENT"])
	}
}

func TestRotate_MissingKey_NoError(t *testing.T) {
	opts := rotator.DefaultOptions()
	rotations := []rotator.Rotation{{OldKey: "NONEXISTENT", NewKey: "NEW_KEY"}}
	_, err := rotator.Rotate(baseEnv(), rotations, opts)
	if err != nil {
		t.Errorf("expected no error for missing key, got %v", err)
	}
}

func TestRotate_MissingKey_FailOnMissing(t *testing.T) {
	opts := rotator.Options{Strategy: rotator.StrategyRemove, FailOnMissing: true}
	rotations := []rotator.Rotation{{OldKey: "NONEXISTENT", NewKey: "NEW_KEY"}}
	_, err := rotator.Rotate(baseEnv(), rotations, opts)
	if err == nil {
		t.Error("expected error for missing key with FailOnMissing=true")
	}
}

func TestRotate_MultipleRotations(t *testing.T) {
	opts := rotator.DefaultOptions()
	rotations := []rotator.Rotation{
		{OldKey: "DB_HOST", NewKey: "DATABASE_HOST"},
		{OldKey: "DB_PASS", NewKey: "DATABASE_PASSWORD"},
	}
	out, err := rotator.Rotate(baseEnv(), rotations, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 { // APP_ENV + 2 new keys
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}
