package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/redactor"
)

func TestRedact_DefaultOptions_MasksSensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "abc123",
		"APP_NAME":    "myapp",
	}
	out := redactor.Redact(env, redactor.DefaultOptions())
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "***" {
		t.Errorf("expected API_TOKEN to be redacted, got %q", out["API_TOKEN"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be kept, got %q", out["APP_NAME"])
	}
}

func TestRedact_CustomMask(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "topsecret"}
	opts := redactor.DefaultOptions()
	opts.Mask = "<REDACTED>"
	out := redactor.Redact(env, opts)
	if out["SECRET_KEY"] != "<REDACTED>" {
		t.Errorf("unexpected mask value: %q", out["SECRET_KEY"])
	}
}

func TestRedact_ExplicitKeys(t *testing.T) {
	env := map[string]string{
		"SOME_VAR":  "visible",
		"OTHER_VAR": "also-visible",
		"HIDE_ME":   "hidden",
	}
	opts := redactor.Options{
		ExplicitKeys: []string{"HIDE_ME"},
		Mask:         "***",
	}
	out := redactor.Redact(env, opts)
	if out["HIDE_ME"] != "***" {
		t.Errorf("expected HIDE_ME to be redacted")
	}
	if out["SOME_VAR"] != "visible" {
		t.Errorf("expected SOME_VAR to remain visible")
	}
}

func TestRedact_EmptyEnv(t *testing.T) {
	out := redactor.Redact(map[string]string{}, redactor.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestRedact_DefaultMaskWhenEmpty(t *testing.T) {
	env := map[string]string{"AUTH_TOKEN": "secret"}
	opts := redactor.Options{SensitivePatterns: []string{"auth"}, Mask: ""}
	out := redactor.Redact(env, opts)
	if out["AUTH_TOKEN"] != "***" {
		t.Errorf("expected default mask, got %q", out["AUTH_TOKEN"])
	}
}
