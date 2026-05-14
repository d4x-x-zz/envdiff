package digester_test

import (
	"testing"

	"github.com/user/envdiff/internal/digester"
)

func TestDigest_Deterministic(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := digester.DefaultOptions()

	r1 := digester.Digest(env, opts)
	r2 := digester.Digest(env, opts)

	if r1.Hex != r2.Hex {
		t.Errorf("expected same digest on repeated calls, got %q and %q", r1.Hex, r2.Hex)
	}
}

func TestDigest_KeyCount(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	r := digester.Digest(env, digester.DefaultOptions())
	if r.KeyCount != 3 {
		t.Errorf("expected KeyCount=3, got %d", r.KeyCount)
	}
}

func TestDigest_ChangedValueProducesDifferentHash(t *testing.T) {
	env1 := map[string]string{"SECRET": "old"}
	env2 := map[string]string{"SECRET": "new"}
	opts := digester.DefaultOptions()

	if digester.Digest(env1, opts).Hex == digester.Digest(env2, opts).Hex {
		t.Error("expected different digest for different values")
	}
}

func TestDigest_KeysOnlyMode_IgnoresValueChange(t *testing.T) {
	env1 := map[string]string{"KEY": "alpha"}
	env2 := map[string]string{"KEY": "beta"}
	opts := digester.Options{IncludeValues: false}

	if digester.Digest(env1, opts).Hex != digester.Digest(env2, opts).Hex {
		t.Error("expected same digest when IncludeValues=false and only values differ")
	}
}

func TestDigest_KeyPrefix_FiltersKeys(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"}
	opts := digester.Options{IncludeValues: true, KeyPrefix: "APP_"}

	r := digester.Digest(env, opts)
	if r.KeyCount != 2 {
		t.Errorf("expected 2 keys with prefix APP_, got %d", r.KeyCount)
	}

	// digest without DB_URL should differ from full digest
	full := digester.Digest(env, digester.DefaultOptions())
	if r.Hex == full.Hex {
		t.Error("expected prefix-filtered digest to differ from full digest")
	}
}

func TestDigest_EmptyMap(t *testing.T) {
	r := digester.Digest(map[string]string{}, digester.DefaultOptions())
	if r.KeyCount != 0 {
		t.Errorf("expected KeyCount=0 for empty map, got %d", r.KeyCount)
	}
	if r.Hex == "" {
		t.Error("expected non-empty hex digest even for empty map")
	}
}
