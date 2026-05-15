package tagger

import (
	"sort"
	"testing"
)

func sortedTags(tags []string) []string {
	out := make([]string, len(tags))
	copy(out, tags)
	sort.Strings(out)
	return out
}

func TestTag_SecretKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "abc123",
	}
	res := Tag(env, DefaultOptions())

	for _, key := range []string{"DB_PASSWORD", "API_TOKEN"} {
		tags := res.Tags[key]
		found := false
		for _, t2 := range tags {
			if t2 == "secret" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected key %q to have tag 'secret', got %v", key, tags)
		}
	}
}

func TestTag_DatabaseKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	res := Tag(env, DefaultOptions())

	for key, tags := range res.Tags {
		found := false
		for _, tg := range tags {
			if tg == "database" {
				found = true
			}
		}
		if !found {
			t.Errorf("key %q: expected 'database' tag, got %v", key, tags)
		}
	}
}

func TestTag_NoMatchReturnsEmptySlice(t *testing.T) {
	env := map[string]string{"APP_NAME": "envdiff"}
	res := Tag(env, DefaultOptions())

	if len(res.Tags["APP_NAME"]) != 0 {
		t.Errorf("expected no tags for APP_NAME, got %v", res.Tags["APP_NAME"])
	}
}

func TestTag_MultipleTagsOnSameKey(t *testing.T) {
	// DB_PASSWORD should match both 'secret' and 'database'
	env := map[string]string{"DB_PASSWORD": "s3cr3t"}
	res := Tag(env, DefaultOptions())

	tags := sortedTags(res.Tags["DB_PASSWORD"])
	if len(tags) < 2 {
		t.Errorf("expected at least 2 tags for DB_PASSWORD, got %v", tags)
	}
}

func TestTag_CaseInsensitiveFalse(t *testing.T) {
	opts := DefaultOptions()
	opts.CaseInsensitive = false
	// lowercase key won't match uppercase patterns
	env := map[string]string{"db_password": "x"}
	res := Tag(env, opts)

	if len(res.Tags["db_password"]) != 0 {
		t.Errorf("expected no tags when case-insensitive is off, got %v", res.Tags["db_password"])
	}
}

func TestTag_CustomRules(t *testing.T) {
	opts := Options{
		Rules:           map[string][]string{"infra": {"K8S_", "HELM_"}},
		CaseInsensitive: true,
	}
	env := map[string]string{"K8S_NAMESPACE": "prod", "APP_ENV": "production"}
	res := Tag(env, opts)

	if tags := res.Tags["K8S_NAMESPACE"]; len(tags) == 0 || tags[0] != "infra" {
		t.Errorf("expected 'infra' tag for K8S_NAMESPACE, got %v", tags)
	}
	if len(res.Tags["APP_ENV"]) != 0 {
		t.Errorf("expected no tags for APP_ENV, got %v", res.Tags["APP_ENV"])
	}
}

func TestTag_EmptyEnv(t *testing.T) {
	res := Tag(map[string]string{}, DefaultOptions())
	if len(res.Tags) != 0 {
		t.Errorf("expected empty tags map, got %v", res.Tags)
	}
}
