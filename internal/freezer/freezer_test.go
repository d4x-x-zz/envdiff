package freezer

import (
	"testing"
)

func TestFreeze_BasicKeys(t *testing.T) {
	env := map[string]string{"B": "2", "A": "1", "C": "3"}
	f := Freeze(env, DefaultOptions())
	if len(f.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(f.Keys))
	}
	if f.Keys[0] != "A" || f.Keys[1] != "B" || f.Keys[2] != "C" {
		t.Errorf("keys not sorted: %v", f.Keys)
	}
}

func TestFreeze_FingerprintDeterministic(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	f1 := Freeze(env, DefaultOptions())
	f2 := Freeze(env, DefaultOptions())
	if f1.Fingerprint != f2.Fingerprint {
		t.Error("fingerprint should be deterministic")
	}
}

func TestFreeze_FingerprintChangesOnValueChange(t *testing.T) {
	env1 := map[string]string{"KEY": "old"}
	env2 := map[string]string{"KEY": "new"}
	f1 := Freeze(env1, DefaultOptions())
	f2 := Freeze(env2, DefaultOptions())
	if f1.Fingerprint == f2.Fingerprint {
		t.Error("fingerprint should differ when values change")
	}
}

func TestFreeze_KeysOnlyMode_IgnoresValueChange(t *testing.T) {
	opts := Options{IncludeValues: false}
	env1 := map[string]string{"KEY": "old"}
	env2 := map[string]string{"KEY": "new"}
	f1 := Freeze(env1, opts)
	f2 := Freeze(env2, opts)
	if f1.Fingerprint != f2.Fingerprint {
		t.Error("keys-only fingerprint should not change when only values differ")
	}
}

func TestFreeze_KeyPrefix_FiltersKeys(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "APP_NAME": "test", "DB_PORT": "5432"}
	f := Freeze(env, Options{IncludeValues: true, KeyPrefix: "DB_"})
	if len(f.Keys) != 2 {
		t.Fatalf("expected 2 keys with DB_ prefix, got %d", len(f.Keys))
	}
}

func TestFrozen_Changed_DetectsAddedKey(t *testing.T) {
	env := map[string]string{"A": "1"}
	f := Freeze(env, DefaultOptions())
	env["B"] = "2"
	if !f.Changed(env, DefaultOptions()) {
		t.Error("expected Changed to return true after adding a key")
	}
}

func TestFrozen_Changed_ReturnsFalseWhenUnchanged(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	f := Freeze(env, DefaultOptions())
	copy := map[string]string{"A": "1", "B": "2"}
	if f.Changed(copy, DefaultOptions()) {
		t.Error("expected Changed to return false for identical env")
	}
}
