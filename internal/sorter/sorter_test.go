package sorter_test

import (
	"testing"

	"envdiff/internal/sorter"
)

func TestSort_Alphabetical(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "1",
		"APPLE": "2",
		"MANGO": "3",
	}
	opts := sorter.Options{Alphabetical: true, GroupByPrefix: false}
	keys := sorter.Sort(env, opts)

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSort_GroupByPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST":    "localhost",
		"AWS_KEY":    "abc",
		"DB_PORT":    "5432",
		"AWS_SECRET": "xyz",
		"APP_NAME":   "envdiff",
	}
	opts := sorter.Options{Alphabetical: false, GroupByPrefix: true}
	keys := sorter.Sort(env, opts)

	if len(keys) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(keys))
	}
	// Groups should be: APP, AWS, DB — alphabetically ordered
	if keys[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME first, got %s", keys[0])
	}
	if keys[1] != "AWS_KEY" || keys[2] != "AWS_SECRET" {
		t.Errorf("unexpected AWS group order: %v", keys[1:3])
	}
	if keys[3] != "DB_HOST" || keys[4] != "DB_PORT" {
		t.Errorf("unexpected DB group order: %v", keys[3:5])
	}
}

func TestSort_EmptyMap(t *testing.T) {
	env := map[string]string{}
	keys := sorter.Sort(env, sorter.DefaultOptions())
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}

func TestSort_NoPrefixKey(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
		"HOST": "localhost",
	}
	opts := sorter.Options{GroupByPrefix: true}
	keys := sorter.Sort(env, opts)

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "HOST" || keys[1] != "PORT" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := sorter.DefaultOptions()
	if !opts.Alphabetical {
		t.Error("expected Alphabetical to be true by default")
	}
	if opts.GroupByPrefix {
		t.Error("expected GroupByPrefix to be false by default")
	}
}
