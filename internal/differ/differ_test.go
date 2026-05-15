package differ

import (
	"sort"
	"testing"
)

func TestDiff_NoDifferences(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	res := Diff(left, right)
	if !res.Clean() {
		t.Errorf("expected clean result, got %+v", res)
	}
}

func TestDiff_MissingInRight(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}
	res := Diff(left, right)
	if len(res.MissingInRight) != 1 || res.MissingInRight[0] != "B" {
		t.Errorf("expected B missing in right, got %v", res.MissingInRight)
	}
}

func TestDiff_MissingInLeft(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "C": "3"}
	res := Diff(left, right)
	if len(res.MissingInLeft) != 1 || res.MissingInLeft[0] != "C" {
		t.Errorf("expected C missing in left, got %v", res.MissingInLeft)
	}
}

func TestDiff_MismatchedValues(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old"}
	right := map[string]string{"A": "1", "B": "new"}
	res := Diff(left, right)
	if len(res.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(res.Mismatched))
	}
	m := res.Mismatched[0]
	if m.Key != "B" || m.LeftValue != "old" || m.RightValue != "new" {
		t.Errorf("unexpected mismatch: %+v", m)
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	res := Diff(map[string]string{}, map[string]string{})
	if !res.Clean() {
		t.Error("expected clean result for empty maps")
	}
}

func TestDiff_AllThreeIssues(t *testing.T) {
	left := map[string]string{"ONLY_LEFT": "x", "SHARED": "old"}
	right := map[string]string{"ONLY_RIGHT": "y", "SHARED": "new"}
	res := Diff(left, right)

	if len(res.MissingInRight) != 1 || res.MissingInRight[0] != "ONLY_LEFT" {
		t.Errorf("MissingInRight: %v", res.MissingInRight)
	}
	if len(res.MissingInLeft) != 1 || res.MissingInLeft[0] != "ONLY_RIGHT" {
		t.Errorf("MissingInLeft: %v", res.MissingInLeft)
	}
	if len(res.Mismatched) != 1 {
		t.Errorf("Mismatched: %v", res.Mismatched)
	}
	if res.TotalIssues() != 3 {
		t.Errorf("TotalIssues: got %d, want 3", res.TotalIssues())
	}
	_ = sort.Search // keep import used
}
