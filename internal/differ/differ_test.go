package differ_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/differ"
)

func TestDiff_NoDifferences(t *testing.T) {
	left := map[string]string{"HOST": "localhost", "PORT": "5432"}
	right := map[string]string{"HOST": "localhost", "PORT": "5432"}

	result := differ.Diff(left, right)

	if result.HasDiff() {
		t.Errorf("expected no diff, got %+v", result)
	}
}

func TestDiff_MissingInRight(t *testing.T) {
	left := map[string]string{"HOST": "localhost", "SECRET": "abc"}
	right := map[string]string{"HOST": "localhost"}

	result := differ.Diff(left, right)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "SECRET" {
		t.Errorf("expected SECRET missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected nothing missing in left, got %v", result.MissingInLeft)
	}
}

func TestDiff_MissingInLeft(t *testing.T) {
	left := map[string]string{"HOST": "localhost"}
	right := map[string]string{"HOST": "localhost", "NEW_KEY": "value"}

	result := differ.Diff(left, right)

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY missing in left, got %v", result.MissingInLeft)
	}
}

func TestDiff_MismatchedValues(t *testing.T) {
	left := map[string]string{"DB_URL": "postgres://old", "PORT": "5432"}
	right := map[string]string{"DB_URL": "postgres://new", "PORT": "5432"}

	result := differ.Diff(left, right)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "DB_URL" || m.LeftValue != "postgres://old" || m.RightValue != "postgres://new" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	result := differ.Diff(map[string]string{}, map[string]string{})
	if result.HasDiff() {
		t.Error("expected no diff for two empty maps")
	}
}

func TestDiff_HasDiff(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"B": "2"}

	result := differ.Diff(left, right)

	if !result.HasDiff() {
		t.Error("expected HasDiff to return true")
	}
}
