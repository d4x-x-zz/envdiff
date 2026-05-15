package profiler

import (
	"testing"
)

func TestAnalyze_EmptyMap(t *testing.T) {
	p := Analyze(map[string]string{}, DefaultOptions())
	if p.TotalKeys != 0 {
		t.Fatalf("expected 0 keys, got %d", p.TotalKeys)
	}
	if p.MinValueLen != 0 {
		t.Fatalf("expected MinValueLen 0, got %d", p.MinValueLen)
	}
}

func TestAnalyze_TotalKeys(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": ""}
	p := Analyze(env, DefaultOptions())
	if p.TotalKeys != 3 {
		t.Fatalf("expected 3, got %d", p.TotalKeys)
	}
}

func TestAnalyze_EmptyValues(t *testing.T) {
	env := map[string]string{"A": "", "B": "", "C": "hello"}
	p := Analyze(env, DefaultOptions())
	if p.EmptyValues != 2 {
		t.Fatalf("expected 2 empty, got %d", p.EmptyValues)
	}
}

func TestAnalyze_Density(t *testing.T) {
	env := map[string]string{"A": "val", "B": "", "C": "val2", "D": ""}
	p := Analyze(env, DefaultOptions())
	// 2 non-empty out of 4
	if p.Density != 0.5 {
		t.Fatalf("expected density 0.5, got %f", p.Density)
	}
}

func TestAnalyze_TypeBreakdown(t *testing.T) {
	env := map[string]string{
		"BOOL":  "true",
		"INT":   "42",
		"FLOAT": "3.14",
		"URL":   "https://example.com",
		"STR":   "hello",
	}
	p := Analyze(env, DefaultOptions())
	if p.TypeBreakdown["bool"] != 1 {
		t.Errorf("expected 1 bool, got %d", p.TypeBreakdown["bool"])
	}
	if p.TypeBreakdown["int"] != 1 {
		t.Errorf("expected 1 int, got %d", p.TypeBreakdown["int"])
	}
	if p.TypeBreakdown["float"] != 1 {
		t.Errorf("expected 1 float, got %d", p.TypeBreakdown["float"])
	}
	if p.TypeBreakdown["url"] != 1 {
		t.Errorf("expected 1 url, got %d", p.TypeBreakdown["url"])
	}
	if p.TypeBreakdown["string"] != 1 {
		t.Errorf("expected 1 string, got %d", p.TypeBreakdown["string"])
	}
}

func TestAnalyze_LengthStats(t *testing.T) {
	env := map[string]string{"A": "hi", "B": "hello", "C": "x"}
	p := Analyze(env, DefaultOptions())
	if p.MaxValueLen != 5 {
		t.Errorf("expected max 5, got %d", p.MaxValueLen)
	}
	if p.MinValueLen != 1 {
		t.Errorf("expected min 1, got %d", p.MinValueLen)
	}
	// avg = (2+5+1)/3 = 2.666...
	if p.AvgValueLen < 2.6 || p.AvgValueLen > 2.7 {
		t.Errorf("unexpected avg %f", p.AvgValueLen)
	}
}

func TestAnalyze_NoTypeBreakdown(t *testing.T) {
	env := map[string]string{"A": "true", "B": "123"}
	opts := DefaultOptions()
	opts.IncludeTypeBreakdown = false
	p := Analyze(env, opts)
	if len(p.TypeBreakdown) != 0 {
		t.Errorf("expected empty type breakdown")
	}
}

func TestSortedTypes(t *testing.T) {
	p := Profile{
		TypeBreakdown: map[string]int{"string": 2, "bool": 1, "int": 3},
	}
	types := SortedTypes(p)
	if types[0] != "bool" || types[1] != "int" || types[2] != "string" {
		t.Errorf("unexpected order: %v", types)
	}
}

// TestAnalyze_AllEmpty verifies that density is 0.0 and EmptyValues equals
// TotalKeys when every value in the map is an empty string.
func TestAnalyze_AllEmpty(t *testing.T) {
	env := map[string]string{"A": "", "B": "", "C": ""}
	p := Analyze(env, DefaultOptions())
	if p.EmptyValues != 3 {
		t.Errorf("expected 3 empty values, got %d", p.EmptyValues)
	}
	if p.Density != 0.0 {
		t.Errorf("expected density 0.0, got %f", p.Density)
	}
}
