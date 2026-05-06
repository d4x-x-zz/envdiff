package merger_test

import (
	"testing"

	"envdiff/internal/merger"
)

func TestMerge_SingleMap(t *testing.T) {
	input := []map[string]string{{"A": "1", "B": "2"}}
	got, err := merger.Merge(input, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "1" || got["B"] != "2" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestMerge_FirstWins(t *testing.T) {
	maps := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	got, err := merger.Merge(maps, merger.Options{Strategy: merger.StrategyFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", got["KEY"])
	}
}

func TestMerge_LastWins(t *testing.T) {
	maps := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	got, err := merger.Merge(maps, merger.Options{Strategy: merger.StrategyLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "second" {
		t.Errorf("expected 'second', got %q", got["KEY"])
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	maps := []map[string]string{
		{"KEY": "alpha"},
		{"KEY": "beta"},
	}
	_, err := merger.Merge(maps, merger.Options{Strategy: merger.StrategyError})
	if err == nil {
		t.Fatal("expected error for conflicting keys, got nil")
	}
}

func TestMerge_NoConflict_ErrorStrategy(t *testing.T) {
	maps := []map[string]string{
		{"KEY": "same"},
		{"KEY": "same"},
	}
	got, err := merger.Merge(maps, merger.Options{Strategy: merger.StrategyError})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "same" {
		t.Errorf("expected 'same', got %q", got["KEY"])
	}
}

func TestMerge_EmptyInput(t *testing.T) {
	got, err := merger.Merge(nil, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMerge_DisjointKeys(t *testing.T) {
	maps := []map[string]string{
		{"A": "1"},
		{"B": "2"},
	}
	got, err := merger.Merge(maps, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "1" || got["B"] != "2" {
		t.Errorf("unexpected result: %v", got)
	}
}
