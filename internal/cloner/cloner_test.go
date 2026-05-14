package cloner_test

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/envdiff/internal/cloner"
)

func base() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"SECRET":  "",
	}
}

func TestClone_NoOptions_IsDeepCopy(t *testing.T) {
	src := base()
	out := cloner.Clone(src, cloner.DefaultOptions())
	if len(out) != len(src) {
		t.Fatalf("expected %d keys, got %d", len(src), len(out))
	}
	src["DB_HOST"] = "changed"
	if out["DB_HOST"] == "changed" {
		t.Fatal("clone shares memory with source")
	}
}

func TestClone_KeyPrefix(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.KeyPrefix = "CLONE_"
	out := cloner.Clone(base(), opts)
	for k := range out {
		if !strings.HasPrefix(k, "CLONE_") {
			t.Errorf("key %q missing prefix", k)
		}
	}
}

func TestClone_KeySuffix(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.KeySuffix = "_COPY"
	out := cloner.Clone(base(), opts)
	for k := range out {
		if !strings.HasSuffix(k, "_COPY") {
			t.Errorf("key %q missing suffix", k)
		}
	}
}

func TestClone_UppercaseKeys(t *testing.T) {
	src := map[string]string{"db_host": "localhost", "db_port": "5432"}
	opts := cloner.DefaultOptions()
	opts.UppercaseKeys = true
	out := cloner.Clone(src, opts)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST after uppercase")
	}
	if _, ok := out["DB_PORT"]; !ok {
		t.Error("expected DB_PORT after uppercase")
	}
}

func TestClone_OmitEmpty(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.OmitEmpty = true
	out := cloner.Clone(base(), opts)
	if _, ok := out["SECRET"]; ok {
		t.Error("expected empty SECRET key to be omitted")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestClone_ValueTransform(t *testing.T) {
	opts := cloner.DefaultOptions()
	opts.ValueTransform = strings.ToUpper
	out := cloner.Clone(map[string]string{"KEY": "hello"}, opts)
	if out["KEY"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["KEY"])
	}
}

func TestClone_EmptyMap(t *testing.T) {
	out := cloner.Clone(map[string]string{}, cloner.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}
