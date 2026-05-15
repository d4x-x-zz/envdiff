package extractor_test

import (
	"testing"

	"github.com/user/envdiff/internal/extractor"
)

var src = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envdiff",
	"APP_VERSION": "1.0.0",
	"SECRET_KEY":  "s3cr3t",
}

func TestExtract_ExplicitKeys(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Keys = []string{"DB_HOST", "APP_NAME"}
	r := extractor.Extract(src, opts)
	if len(r.Env) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Env))
	}
	if r.Env["DB_HOST"] != "localhost" {
		t.Errorf("unexpected value for DB_HOST: %s", r.Env["DB_HOST"])
	}
	if len(r.Missed) != 0 {
		t.Errorf("expected no missed keys, got %v", r.Missed)
	}
}

func TestExtract_GlobPattern(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Patterns = []string{"DB_*"}
	r := extractor.Extract(src, opts)
	if len(r.Env) != 2 {
		t.Fatalf("expected 2 DB_ keys, got %d", len(r.Env))
	}
	if _, ok := r.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := r.Env["DB_PORT"]; !ok {
		t.Error("expected DB_PORT in result")
	}
}

func TestExtract_MultiplePatterns(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Patterns = []string{"APP_*", "SECRET_*"}
	r := extractor.Extract(src, opts)
	if len(r.Env) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Env))
	}
}

func TestExtract_MissingKey_Tracked(t *testing.T) {
	opts := extractor.Options{IgnoreMissing: false, Keys: []string{"DB_HOST", "MISSING_KEY"}}
	r := extractor.Extract(src, opts)
	if len(r.Missed) != 1 || r.Missed[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in Missed, got %v", r.Missed)
	}
	if _, ok := r.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to still be extracted")
	}
}

func TestExtract_NoPatterns_ReturnsEmpty(t *testing.T) {
	opts := extractor.DefaultOptions()
	r := extractor.Extract(src, opts)
	if len(r.Env) != 0 {
		t.Errorf("expected empty result with no keys/patterns, got %d", len(r.Env))
	}
}

func TestExtract_EmptySrc(t *testing.T) {
	opts := extractor.DefaultOptions()
	opts.Keys = []string{"FOO"}
	r := extractor.Extract(map[string]string{}, opts)
	if len(r.Env) != 0 {
		t.Errorf("expected empty env, got %d keys", len(r.Env))
	}
	if len(r.Missed) != 1 {
		t.Errorf("expected 1 missed key, got %v", r.Missed)
	}
}
