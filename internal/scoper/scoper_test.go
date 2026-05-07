package scoper_test

import (
	"testing"

	"github.com/user/envdiff/internal/scoper"
)

func TestScope_StripPrefix(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
	}
	opts := scoper.Options{Prefix: "APP_", StripPrefix: true}
	res := scoper.Scope(env, opts)

	if len(res.Scoped) != 2 {
		t.Fatalf("expected 2 scoped keys, got %d", len(res.Scoped))
	}
	if res.Scoped["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", res.Scoped["HOST"])
	}
	if res.Scoped["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", res.Scoped["PORT"])
	}
	if len(res.Excluded) != 1 {
		t.Fatalf("expected 1 excluded key, got %d", len(res.Excluded))
	}
}

func TestScope_KeepPrefix(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost",
		"DB_HOST":  "db.local",
	}
	opts := scoper.Options{Prefix: "APP_", StripPrefix: false}
	res := scoper.Scope(env, opts)

	if _, ok := res.Scoped["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in scoped map")
	}
}

func TestScope_EmptyPrefix(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAR": "2"}
	opts := scoper.Options{Prefix: "", StripPrefix: true}
	res := scoper.Scope(env, opts)

	if len(res.Scoped) != 2 {
		t.Errorf("expected all keys scoped, got %d", len(res.Scoped))
	}
	if len(res.Excluded) != 0 {
		t.Errorf("expected no excluded keys, got %d", len(res.Excluded))
	}
}

func TestScope_EmptyMap(t *testing.T) {
	res := scoper.Scope(map[string]string{}, scoper.DefaultOptions())
	if len(res.Scoped) != 0 || len(res.Excluded) != 0 {
		t.Error("expected empty results for empty input")
	}
}

func TestScope_NoMatch(t *testing.T) {
	env := map[string]string{"DB_HOST": "db.local", "DB_PORT": "5432"}
	opts := scoper.Options{Prefix: "APP_", StripPrefix: true}
	res := scoper.Scope(env, opts)

	if len(res.Scoped) != 0 {
		t.Errorf("expected 0 scoped keys, got %d", len(res.Scoped))
	}
	if len(res.Excluded) != 2 {
		t.Errorf("expected 2 excluded keys, got %d", len(res.Excluded))
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys := scoper.SortedKeys(m)
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("pos %d: expected %s got %s", i, expected[i], k)
		}
	}
}
