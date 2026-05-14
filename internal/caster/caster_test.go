package caster_test

import (
	"testing"

	"github.com/user/envdiff/internal/caster"
)

func resultMap(results []caster.Result) map[string]caster.Type {
	m := make(map[string]caster.Type, len(results))
	for _, r := range results {
		m[r.Key] = r.Type
	}
	return m
}

func TestCast_BoolValues(t *testing.T) {
	env := map[string]string{"ENABLED": "true", "DEBUG": "false", "FLAG": "1"}
	res := resultMap(caster.Cast(env, caster.DefaultOptions()))
	for k, want := range map[string]caster.Type{"ENABLED": caster.TypeBool, "DEBUG": caster.TypeBool, "FLAG": caster.TypeBool} {
		if res[k] != want {
			t.Errorf("key %s: got %s, want %s", k, res[k], want)
		}
	}
}

func TestCast_IntValues(t *testing.T) {
	env := map[string]string{"PORT": "8080", "TIMEOUT": "30"}
	res := resultMap(caster.Cast(env, caster.DefaultOptions()))
	for k := range env {
		if res[k] != caster.TypeInt {
			t.Errorf("key %s: got %s, want int", k, res[k])
		}
	}
}

func TestCast_FloatValues(t *testing.T) {
	env := map[string]string{"RATE": "3.14", "THRESHOLD": "0.5"}
	res := resultMap(caster.Cast(env, caster.DefaultOptions()))
	for k := range env {
		if res[k] != caster.TypeFloat {
			t.Errorf("key %s: got %s, want float", k, res[k])
		}
	}
}

func TestCast_URLValues(t *testing.T) {
	env := map[string]string{"API_URL": "https://api.example.com", "HOOK": "http://localhost/hook"}
	res := resultMap(caster.Cast(env, caster.DefaultOptions()))
	for k := range env {
		if res[k] != caster.TypeURL {
			t.Errorf("key %s: got %s, want url", k, res[k])
		}
	}
}

func TestCast_DSNValues(t *testing.T) {
	env := map[string]string{"DATABASE_URL": "postgres://user:pass@localhost/mydb"}
	res := resultMap(caster.Cast(env, caster.DefaultOptions()))
	if res["DATABASE_URL"] != caster.TypeDSN {
		t.Errorf("got %s, want dsn", res["DATABASE_URL"])
	}
}

func TestCast_StringFallback(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "REGION": "us-east-1"}
	res := resultMap(caster.Cast(env, caster.DefaultOptions()))
	for k := range env {
		if res[k] != caster.TypeString {
			t.Errorf("key %s: got %s, want string", k, res[k])
		}
	}
}

func TestCast_DisableURL(t *testing.T) {
	opts := caster.DefaultOptions()
	opts.DetectURL = false
	env := map[string]string{"ENDPOINT": "https://example.com"}
	res := resultMap(caster.Cast(env, opts))
	if res["ENDPOINT"] != caster.TypeString {
		t.Errorf("got %s, want string when URL detection disabled", res["ENDPOINT"])
	}
}

func TestCast_EmptyMap(t *testing.T) {
	res := caster.Cast(map[string]string{}, caster.DefaultOptions())
	if len(res) != 0 {
		t.Errorf("expected empty result, got %d items", len(res))
	}
}
