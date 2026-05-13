package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/grouper"
)

func TestGroup_BasicPrefixes(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "myapp",
	}
	opts := grouper.DefaultOptions()
	groups := grouper.Group(env, opts)

	if len(groups["DB"]) != 2 {
		t.Fatalf("expected 2 DB keys, got %d", len(groups["DB"]))
	}
	if groups["APP"]["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME in APP group")
	}
}

func TestGroup_NoSeparator_GoesToOther(t *testing.T) {
	env := map[string]string{
		"HOSTNAME": "box1",
		"DB_HOST":  "localhost",
	}
	opts := grouper.DefaultOptions()
	groups := grouper.Group(env, opts)

	if _, ok := groups["OTHER"]["HOSTNAME"]; !ok {
		t.Errorf("expected HOSTNAME in OTHER group")
	}
}

func TestGroup_AllowList_FiltersGroups(t *testing.T) {
	env := map[string]string{
		"DB_HOST":    "localhost",
		"CACHE_HOST": "redis",
		"APP_NAME":   "myapp",
	}
	opts := grouper.DefaultOptions()
	opts.AllowList = []string{"DB"}
	groups := grouper.Group(env, opts)

	if _, ok := groups["DB"]; !ok {
		t.Errorf("expected DB group")
	}
	if _, ok := groups["OTHER"]["CACHE_HOST"]; !ok {
		t.Errorf("expected CACHE_HOST in OTHER (not in allowlist)")
	}
	if _, ok := groups["OTHER"]["APP_NAME"]; !ok {
		t.Errorf("expected APP_NAME in OTHER (not in allowlist)")
	}
}

func TestGroup_EmptyMap(t *testing.T) {
	groups := grouper.Group(map[string]string{}, grouper.DefaultOptions())
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestSortedGroupNames(t *testing.T) {
	groups := map[string]map[string]string{
		"Z": {"Z_KEY": "1"},
		"A": {"A_KEY": "2"},
		"M": {"M_KEY": "3"},
	}
	names := grouper.SortedGroupNames(groups)
	expected := []string{"A", "M", "Z"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("pos %d: want %s got %s", i, expected[i], n)
		}
	}
}

func TestGroup_CustomSeparator(t *testing.T) {
	env := map[string]string{
		"DB.HOST": "localhost",
		"DB.PORT": "5432",
	}
	opts := grouper.DefaultOptions()
	opts.Separator = "."
	groups := grouper.Group(env, opts)

	if len(groups["DB"]) != 2 {
		t.Fatalf("expected 2 DB keys with dot separator, got %d", len(groups["DB"]))
	}
}
