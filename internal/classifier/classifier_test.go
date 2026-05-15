package classifier_test

import (
	"sort"
	"testing"

	"github.com/user/envdiff/internal/classifier"
)

func sortedKeys(keys []string) []string {
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	return sorted
}

func TestClassify_DatabaseKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "POSTGRES_USER": "admin"}
	res := classifier.Classify(env, classifier.DefaultOptions())
	keys := sortedKeys(res.Categories[classifier.CategoryDatabase])
	if len(keys) != 2 {
		t.Fatalf("expected 2 database keys, got %d", len(keys))
	}
}

func TestClassify_AuthKeys(t *testing.T) {
	env := map[string]string{"JWT_SECRET": "abc", "API_KEY": "xyz", "APP_NAME": "myapp"}
	res := classifier.Classify(env, classifier.DefaultOptions())
	keys := res.Categories[classifier.CategoryAuth]
	if len(keys) != 2 {
		t.Fatalf("expected 2 auth keys, got %d: %v", len(keys), keys)
	}
}

func TestClassify_NetworkKeys(t *testing.T) {
	env := map[string]string{"HTTP_PORT": "8080", "BASE_URL": "http://example.com"}
	res := classifier.Classify(env, classifier.DefaultOptions())
	keys := res.Categories[classifier.CategoryNetwork]
	if len(keys) != 2 {
		t.Fatalf("expected 2 network keys, got %d: %v", len(keys), keys)
	}
}

func TestClassify_OtherIncluded(t *testing.T) {
	env := map[string]string{"APP_NAME": "envdiff", "VERSION": "1.0"}
	opts := classifier.DefaultOptions()
	opts.IncludeOther = true
	res := classifier.Classify(env, opts)
	if len(res.Categories[classifier.CategoryOther]) != 2 {
		t.Fatalf("expected 2 other keys, got %d", len(res.Categories[classifier.CategoryOther]))
	}
}

func TestClassify_OtherExcluded(t *testing.T) {
	env := map[string]string{"APP_NAME": "envdiff", "VERSION": "1.0"}
	opts := classifier.DefaultOptions()
	opts.IncludeOther = false
	res := classifier.Classify(env, opts)
	if _, ok := res.Categories[classifier.CategoryOther]; ok {
		t.Fatal("expected other category to be absent")
	}
}

func TestClassify_EmptyMap(t *testing.T) {
	res := classifier.Classify(map[string]string{}, classifier.DefaultOptions())
	for cat, keys := range res.Categories {
		if len(keys) != 0 {
			t.Errorf("expected no keys for category %s, got %d", cat, len(keys))
		}
	}
}

func TestClassify_StorageKeys(t *testing.T) {
	env := map[string]string{"S3_BUCKET": "my-bucket", "UPLOAD_PATH": "/tmp"}
	res := classifier.Classify(env, classifier.DefaultOptions())
	keys := res.Categories[classifier.CategoryStorage]
	if len(keys) != 2 {
		t.Fatalf("expected 2 storage keys, got %d: %v", len(keys), keys)
	}
}
