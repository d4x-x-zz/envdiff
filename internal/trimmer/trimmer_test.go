package trimmer_test

import (
	"testing"

	"github.com/user/envdiff/internal/trimmer"
)

func TestTrim_RemovesUnknownKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	required := []string{"A", "C"}

	res, err := trimmer.Trim(src, required, trimmer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Trimmed) != 2 {
		t.Fatalf("expected 2 trimmed keys, got %d", len(res.Trimmed))
	}
	if _, ok := res.Trimmed["B"]; ok {
		t.Error("expected B to be removed")
	}
	if len(res.Removed) != 1 || res.Removed[0] != "B" {
		t.Errorf("expected Removed=[B], got %v", res.Removed)
	}
}

func TestTrim_KeepUnknown(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	required := []string{"A"}
	opts := trimmer.DefaultOptions()
	opts.KeepUnknown = true

	res, err := trimmer.Trim(src, required, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Trimmed) != 2 {
		t.Fatalf("expected 2 trimmed keys, got %d", len(res.Trimmed))
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removed keys, got %v", res.Removed)
	}
}

func TestTrim_MissingRequiredKey(t *testing.T) {
	src := map[string]string{"A": "1"}
	required := []string{"A", "B"}

	res, err := trimmer.Trim(src, required, trimmer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "B" {
		t.Errorf("expected Missing=[B], got %v", res.Missing)
	}
}

func TestTrim_StrictMode_ReturnError(t *testing.T) {
	src := map[string]string{"A": "1"}
	required := []string{"A", "MISSING"}
	opts := trimmer.DefaultOptions()
	opts.Strict = true

	_, err := trimmer.Trim(src, required, opts)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestTrim_EmptySource(t *testing.T) {
	res, err := trimmer.Trim(map[string]string{}, []string{"A"}, trimmer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Trimmed) != 0 {
		t.Errorf("expected empty trimmed map")
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 missing key, got %d", len(res.Missing))
	}
}

func TestTrim_EmptyRequired(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	res, err := trimmer.Trim(src, []string{}, trimmer.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Trimmed) != 0 {
		t.Errorf("expected empty trimmed map when required is empty")
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed keys, got %d", len(res.Removed))
	}
}
