package splitter

import (
	"testing"
)

var sampleEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"REDIS_HOST":  "127.0.0.1",
	"REDIS_PORT":  "6379",
	"APP_NAME":    "envdiff",
	"SECRET_KEY":  "abc123",
}

func TestSplit_BasicPrefixes(t *testing.T) {
	opts := DefaultOptions()
	opts.Prefixes = []string{"DB_", "REDIS_"}

	groups := Split(sampleEnv, opts)

	if len(groups["DB_"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(groups["DB_"]))
	}
	if groups["DB_"]["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", groups["DB_"]["HOST"])
	}
	if len(groups["REDIS_"]) != 2 {
		t.Errorf("expected 2 REDIS keys, got %d", len(groups["REDIS_"]))
	}
}

func TestSplit_StripPrefix_False(t *testing.T) {
	opts := DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	opts.StripPrefix = false

	groups := Split(sampleEnv, opts)

	if _, ok := groups["DB_"]["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST to be present when StripPrefix=false")
	}
	if _, ok := groups["DB_"]["HOST"]; ok {
		t.Error("expected key HOST to be absent when StripPrefix=false")
	}
}

func TestSplit_IncludeOther_True(t *testing.T) {
	opts := DefaultOptions()
	opts.Prefixes = []string{"DB_", "REDIS_"}

	groups := Split(sampleEnv, opts)

	other := groups["_other"]
	if len(other) != 2 {
		t.Errorf("expected 2 other keys, got %d", len(other))
	}
	if other["APP_NAME"] != "envdiff" {
		t.Errorf("expected APP_NAME in _other, got %v", other)
	}
}

func TestSplit_IncludeOther_False(t *testing.T) {
	opts := DefaultOptions()
	opts.Prefixes = []string{"DB_"}
	opts.IncludeOther = false

	groups := Split(sampleEnv, opts)

	if _, ok := groups["_other"]; ok {
		t.Error("expected _other group to be absent")
	}
}

func TestSplit_EmptyMap(t *testing.T) {
	opts := DefaultOptions()
	opts.Prefixes = []string{"DB_"}

	groups := Split(map[string]string{}, opts)
	if len(groups) != 0 {
		t.Errorf("expected empty result, got %d groups", len(groups))
	}
}

func TestSortedGroupNames_OtherLast(t *testing.T) {
	groups := map[string]map[string]string{
		"_other": {"X": "1"},
		"REDIS_": {"HOST": "x"},
		"DB_":    {"PORT": "5432"},
	}
	names := SortedGroupNames(groups)
	if names[len(names)-1] != "_other" {
		t.Errorf("expected _other last, got %v", names)
	}
	if names[0] != "DB_" || names[1] != "REDIS_" {
		t.Errorf("expected sorted prefixes first, got %v", names)
	}
}
