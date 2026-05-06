package differ

import "testing"

func TestResult_Clean_WhenEmpty(t *testing.T) {
	r := Result{}
	if !r.Clean() {
		t.Error("empty result should be clean")
	}
}

func TestResult_Clean_WhenMissingRight(t *testing.T) {
	r := Result{MissingInRight: []string{"FOO"}}
	if r.Clean() {
		t.Error("result with MissingInRight should not be clean")
	}
}

func TestResult_Clean_WhenMissingLeft(t *testing.T) {
	r := Result{MissingInLeft: []string{"BAR"}}
	if r.Clean() {
		t.Error("result with MissingInLeft should not be clean")
	}
}

func TestResult_Clean_WhenMismatch(t *testing.T) {
	r := Result{Mismatched: []MismatchEntry{{Key: "X", LeftValue: "a", RightValue: "b"}}}
	if r.Clean() {
		t.Error("result with mismatches should not be clean")
	}
}

func TestResult_TotalIssues(t *testing.T) {
	r := Result{
		MissingInRight: []string{"A", "B"},
		MissingInLeft:  []string{"C"},
		Mismatched:     []MismatchEntry{{Key: "D"}},
	}
	if got := r.TotalIssues(); got != 4 {
		t.Errorf("expected 4 total issues, got %d", got)
	}
}

func TestResult_TotalIssues_Clean(t *testing.T) {
	r := Result{}
	if got := r.TotalIssues(); got != 0 {
		t.Errorf("expected 0 issues, got %d", got)
	}
}
