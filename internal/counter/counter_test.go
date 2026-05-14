package counter_test

import (
	"testing"

	"envdiff/internal/counter"
)

func TestCount_EmptyInput(t *testing.T) {
	r := counter.Count(nil, counter.DefaultOptions())
	if r.Total != 0 {
		t.Fatalf("expected 0 total, got %d", r.Total)
	}
}

func TestCount_SingleMap(t *testing.T) {
	m := map[string]string{"FOO": "1", "BAR": "2"}
	r := counter.Count([]map[string]string{m}, counter.DefaultOptions())
	if r.Total != 2 {
		t.Fatalf("expected 2 total, got %d", r.Total)
	}
}

func TestCount_KeyPresentInAllMaps(t *testing.T) {
	maps := []map[string]string{
		{"FOO": "a", "BAR": "b"},
		{"FOO": "c"},
		{"FOO": "d", "BAZ": "e"},
	}
	r := counter.Count(maps, counter.DefaultOptions())
	for _, e := range r.Entries {
		if e.Key == "FOO" && e.Count != 3 {
			t.Fatalf("FOO: expected count 3, got %d", e.Count)
		}
	}
}

func TestCount_KeyPrefix_FiltersKeys(t *testing.T) {
	maps := []map[string]string{
		{"DB_HOST": "localhost", "APP_NAME": "test", "DB_PORT": "5432"},
	}
	opts := counter.DefaultOptions()
	opts.KeyPrefix = "DB_"
	r := counter.Count(maps, opts)
	if r.Total != 2 {
		t.Fatalf("expected 2 DB_ keys, got %d", r.Total)
	}
	for _, e := range r.Entries {
		if e.Key == "APP_NAME" {
			t.Fatal("APP_NAME should have been filtered out")
		}
	}
}

func TestCount_CaseInsensitive_MergesKeys(t *testing.T) {
	maps := []map[string]string{
		{"foo": "1"},
		{"FOO": "2"},
	}
	opts := counter.DefaultOptions()
	opts.CaseSensitive = false
	r := counter.Count(maps, opts)
	if r.Total != 1 {
		t.Fatalf("expected 1 merged key, got %d", r.Total)
	}
	if r.Entries[0].Count != 2 {
		t.Fatalf("expected count 2, got %d", r.Entries[0].Count)
	}
}

func TestCount_EntriesSortedAlphabetically(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"}
	r := counter.Count([]map[string]string{m}, counter.DefaultOptions())
	keys := make([]string, len(r.Entries))
	for i, e := range r.Entries {
		keys[i] = e.Key
	}
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("position %d: expected %s, got %s", i, k, keys[i])
		}
	}
}
