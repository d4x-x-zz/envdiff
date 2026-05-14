package flattener_test

import (
	"testing"

	"github.com/user/envdiff/internal/flattener"
)

func TestFlatten_NoSeparator_PassThrough(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"PORT":     "8080",
	}
	out := flattener.Flatten(env, flattener.DefaultOptions())
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %s", out["APP_HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", out["PORT"])
	}
}

func TestFlatten_BasicNesting(t *testing.T) {
	env := map[string]string{
		"APP__DB__HOST": "db.local",
		"APP__DB__PORT": "5432",
	}
	out := flattener.Flatten(env, flattener.DefaultOptions())
	if _, ok := out["APP.DB.HOST"]; !ok {
		t.Error("expected key APP.DB.HOST")
	}
	if out["APP.DB.PORT"] != "5432" {
		t.Errorf("expected APP.DB.PORT=5432, got %s", out["APP.DB.PORT"])
	}
}

func TestFlatten_LowercaseKeys(t *testing.T) {
	env := map[string]string{
		"APP__HOST": "localhost",
	}
	opts := flattener.DefaultOptions()
	opts.LowercaseKeys = true
	out := flattener.Flatten(env, opts)
	if _, ok := out["app.host"]; !ok {
		t.Error("expected lowercase key app.host")
	}
}

func TestFlatten_DepthLimit(t *testing.T) {
	env := map[string]string{
		"A__B__C__D": "val",
	}
	opts := flattener.DefaultOptions()
	opts.Depth = 2
	out := flattener.Flatten(env, opts)
	// With depth=2: split into [A, B, C__D] -> A.B.C__D
	if _, ok := out["A.B.C__D"]; !ok {
		t.Errorf("expected key A.B.C__D, got keys: %v", keys(out))
	}
}

func TestFlatten_EmptyMap(t *testing.T) {
	out := flattener.Flatten(map[string]string{}, flattener.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	env := map[string]string{
		"APP::HOST": "localhost",
	}
	opts := flattener.DefaultOptions()
	opts.Separator = "::"
	out := flattener.Flatten(env, opts)
	if _, ok := out["APP.HOST"]; !ok {
		t.Errorf("expected APP.HOST, got keys: %v", keys(out))
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := flattener.DefaultOptions()
	if opts.Separator != "__" {
		t.Errorf("expected separator '__', got %q", opts.Separator)
	}
	if opts.Depth != 0 {
		t.Errorf("expected depth 0, got %d", opts.Depth)
	}
	if opts.LowercaseKeys {
		t.Error("expected LowercaseKeys=false by default")
	}
}

func keys(m map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
