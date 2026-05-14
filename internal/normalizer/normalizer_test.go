package normalizer

import (
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "api_key": "abc"}
	opts := DefaultOptions()
	out := Normalize(env, opts)

	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be present")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY to be present")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	env := map[string]string{"HOST": "  localhost  ", "PORT": "\t8080\n"}
	opts := DefaultOptions()
	out := Normalize(env, opts)

	if out["HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", out["PORT"])
	}
}

func TestNormalize_StripQuotes_DoubleQuotes(t *testing.T) {
	env := map[string]string{"SECRET": `"my secret"`}
	opts := Options{StripQuotes: true, TrimValues: false, UppercaseKeys: false}
	out := Normalize(env, opts)

	if out["SECRET"] != "my secret" {
		t.Errorf("expected 'my secret', got %q", out["SECRET"])
	}
}

func TestNormalize_StripQuotes_SingleQuotes(t *testing.T) {
	env := map[string]string{"TOKEN": "'abc123'"}
	opts := Options{StripQuotes: true}
	out := Normalize(env, opts)

	if out["TOKEN"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", out["TOKEN"])
	}
}

func TestNormalize_StripQuotes_MismatchedQuotes(t *testing.T) {
	env := map[string]string{"VAL": `"oops'`}
	opts := Options{StripQuotes: true}
	out := Normalize(env, opts)

	if out["VAL"] != `"oops'` {
		t.Errorf("mismatched quotes should not be stripped, got %q", out["VAL"])
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"lower_key": "  val  "}
	opts := DefaultOptions()
	Normalize(env, opts)

	if _, ok := env["lower_key"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestNormalize_EmptyMap(t *testing.T) {
	out := Normalize(map[string]string{}, DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}
