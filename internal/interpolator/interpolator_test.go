package interpolator

import (
	"os"
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := Interpolate(env, nil, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://${HOST}:8080",
	}
	out, err := Interpolate(env, nil, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://localhost:8080" {
		t.Errorf("got %q", out["URL"])
	}
}

func TestInterpolate_OverrideTakesPrecedence(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://${HOST}:8080",
	}
	overrides := map[string]string{"HOST": "example.com"}
	out, err := Interpolate(env, overrides, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://example.com:8080" {
		t.Errorf("got %q", out["URL"])
	}
}

func TestInterpolate_FallbackToEnv(t *testing.T) {
	os.Setenv("SYSTEM_HOST", "syshost")
	defer os.Unsetenv("SYSTEM_HOST")

	env := map[string]string{"URL": "http://${SYSTEM_HOST}:9000"}
	opts := DefaultOptions()
	opts.FallbackToEnv = true

	out, err := Interpolate(env, nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://syshost:9000" {
		t.Errorf("got %q", out["URL"])
	}
}

func TestInterpolate_IgnoreMissing_True(t *testing.T) {
	env := map[string]string{"URL": "http://${UNDEFINED}:8080"}
	opts := DefaultOptions()
	opts.IgnoreMissing = true

	out, err := Interpolate(env, nil, opts)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if out["URL"] != "http://${UNDEFINED}:8080" {
		t.Errorf("expected original placeholder preserved, got %q", out["URL"])
	}
}

func TestInterpolate_IgnoreMissing_False(t *testing.T) {
	env := map[string]string{"URL": "http://${UNDEFINED}:8080"}
	opts := DefaultOptions()
	opts.IgnoreMissing = false

	_, err := Interpolate(env, nil, opts)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestInterpolate_EmptyMap(t *testing.T) {
	out, err := Interpolate(map[string]string{}, nil, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty result, got %v", out)
	}
}
