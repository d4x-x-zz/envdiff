package encoder

import (
	"strings"
	"testing"
)

func TestEncode_Shell(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := DefaultOptions()
	out, err := Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `export BAZ="qux"`) {
		t.Errorf("expected shell export for BAZ, got:\n%s", out)
	}
	if !strings.Contains(out, `export FOO="bar"`) {
		t.Errorf("expected shell export for FOO, got:\n%s", out)
	}
}

func TestEncode_Docker(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	opts := DefaultOptions()
	opts.Format = FormatDocker
	out, err := Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected docker line for HOST, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=5432") {
		t.Errorf("expected docker line for PORT, got:\n%s", out)
	}
}

func TestEncode_YAML(t *testing.T) {
	env := map[string]string{"APP_NAME": "myapp", "DB_URL": "postgres://user:pass@host/db"}
	opts := DefaultOptions()
	opts.Format = FormatYAML
	out, err := Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME: myapp") {
		t.Errorf("expected yaml line for APP_NAME, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_URL:") {
		t.Errorf("expected yaml line for DB_URL, got:\n%s", out)
	}
}

func TestEncode_OmitEmpty(t *testing.T) {
	env := map[string]string{"KEY": "value", "EMPTY": ""}
	opts := DefaultOptions()
	opts.OmitEmpty = true
	out, err := Encode(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY to be omitted, got:\n%s", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY to be present, got:\n%s", out)
	}
}

func TestEncode_UnknownFormat(t *testing.T) {
	env := map[string]string{"K": "v"}
	opts := DefaultOptions()
	opts.Format = Format("toml")
	_, err := Encode(env, opts)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.Format != FormatShell {
		t.Errorf("expected default format shell, got %s", opts.Format)
	}
	if !opts.SortKeys {
		t.Error("expected SortKeys to be true by default")
	}
}
