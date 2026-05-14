package stripper_test

import (
	"testing"

	"github.com/user/envdiff/internal/stripper"
)

func base() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"APP_SECRET":  "s3cr3t",
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "pass",
		"DEBUG":       "true",
		"LOG_LEVEL":   "info",
	}
}

func TestStrip_ExplicitKeys(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Keys = []string{"DEBUG", "LOG_LEVEL"}
	out, removed := stripper.Strip(base(), opts)
	if len(removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(removed))
	}
	if _, ok := out["DEBUG"]; ok {
		t.Error("DEBUG should have been stripped")
	}
	if _, ok := out["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL should have been stripped")
	}
}

func TestStrip_ByPrefix(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	out, removed := stripper.Strip(base(), opts)
	if len(removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(removed))
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should have been stripped")
	}
}

func TestStrip_BySuffix(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Suffixes = []string{"_SECRET", "_PASSWORD"}
	out, removed := stripper.Strip(base(), opts)
	if len(removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(removed))
	}
	if _, ok := out["APP_SECRET"]; ok {
		t.Error("APP_SECRET should have been stripped")
	}
}

func TestStrip_DryRun(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Keys = []string{"APP_NAME"}
	opts.DryRun = true
	out, removed := stripper.Strip(base(), opts)
	if len(removed) != 1 {
		t.Fatalf("expected 1 in removed list, got %d", len(removed))
	}
	if _, ok := out["APP_NAME"]; !ok {
		t.Error("DryRun should not mutate the map")
	}
}

func TestStrip_NoRules(t *testing.T) {
	opts := stripper.DefaultOptions()
	out, removed := stripper.Strip(base(), opts)
	if len(removed) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(removed))
	}
	if len(out) != len(base()) {
		t.Errorf("map should be unchanged")
	}
}

func TestStrip_EmptyMap(t *testing.T) {
	opts := stripper.DefaultOptions()
	opts.Keys = []string{"ANYTHING"}
	out, removed := stripper.Strip(map[string]string{}, opts)
	if len(removed) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(removed))
	}
	if len(out) != 0 {
		t.Error("output map should be empty")
	}
}
